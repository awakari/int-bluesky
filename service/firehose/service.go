package firehose

import (
	"bytes"
	"fmt"
	"github.com/awakari/int-bluesky/model"
	"github.com/fxamacker/cbor/v2"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipld/go-car"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/reiver/go-bsky/firehose"
	"io"
	"strings"
)

type commitMessage struct {
	Repo string   `cbor:"repo"`
	Ops  []repoOp `cbor:"ops"`
	Time string   `cbor:"time"`
}

type repoOp struct {
	Action string `cbor:"action"`
	Path   string `cbor:"path"`
}

const prefixDid = "did:plc:"
const typeKey = "$type"
const typeValPost = "app.bsky.feed.post"
const fmtUserId = "https://bsky.app/profile/%s"
const fmtUrlPost = "https://bsky.app/profile/%s/post%s"

func DecodePost(data []byte) (raw map[string]any, userId string, err error) {

	var h any
	var dataPayload []byte
	if dataPayload, err = cbor.UnmarshalFirst(data, &h); err != nil { // skip the header
		err = fmt.Errorf("firehose: failed to decode the CBOR header: %s", err)
	}

	var cm commitMessage
	if err == nil {
		if _, err = cbor.UnmarshalFirst(dataPayload, &cm); err != nil {
			err = fmt.Errorf("firehose: failed to decode the CBOR payload: %s", err)
		}
	}

	raw = make(map[string]any)
	if err == nil {
		var repo string
		switch strings.HasPrefix(cm.Repo, prefixDid) {
		case true:
			repo = cm.Repo
			userId = fmt.Sprintf(fmtUserId, repo)
			raw[model.CeKeySubject] = userId
		default:
			return
		}
		var rKey string
		for _, op := range cm.Ops {
			if op.Action == "create" && strings.HasPrefix(op.Path, typeValPost) {
				rKey = strings.TrimPrefix(op.Path, typeValPost)
				break
			}
		}
		switch rKey {
		case "":
			return
		default:
			raw[model.CeKeyObjectUrl] = fmt.Sprintf(fmtUrlPost, repo, rKey)
		}
	}

	msg := firehose.Message(data)
	var header firehose.MessageHeader
	var payload firehose.MessagePayload
	if err = msg.Decode(&header, &payload); err != nil {
		err = fmt.Errorf("firehose: failed to decode message: %s", err)
	}

	var br *car.CarReader
	if err == nil {
		if br, err = payload.Blocks(); err != nil {
			err = fmt.Errorf("firehose: failed get the message payload blocks: %s", err)
		}
	}

	if err == nil {
		for {
			var block blocks.Block
			block, err = br.Next()
			if err == io.EOF {
				err = nil
				break
			}
			if err != nil {
				err = fmt.Errorf("firehose: failed to iterate blocks: %s", err)
				break
			}
			if err = decodeBlock(block, raw); err != nil {
				break
			}
		}
	}

	if t, tOk := raw[typeKey]; !tOk || t != typeValPost {
		raw = nil // discard
	}
	if _, uOk := raw[model.CeKeyObjectUrl]; !uOk {
		raw = nil // discard
	}

	return
}

func decodeBlock(block blocks.Block, raw map[string]any) (err error) {

	nb := basicnode.Prototype.Any.NewBuilder()
	err = dagcbor.Decode(nb, bytes.NewReader(block.RawData()))
	if err != nil {
		err = fmt.Errorf("firehose: failed to decode the block: %s", err)
	}

	if err == nil {
		node := nb.Build()
		iter := node.MapIterator()
		for !iter.Done() {
			var k, v datamodel.Node
			k, v, err = iter.Next()
			if err != nil {
				err = fmt.Errorf("firehose: failed to iterate the block nodes: %s", err)
				break
			}
			ks, _ := k.AsString()
			switch v.Kind() {
			case datamodel.Kind_Bool:
				raw[ks], _ = v.AsBool()
			case datamodel.Kind_Float:
				raw[ks], _ = v.AsFloat()
			case datamodel.Kind_Int:
				raw[ks], _ = v.AsInt()
			case datamodel.Kind_String:
				raw[ks], _ = v.AsString()
			}
		}
	}

	return
}
