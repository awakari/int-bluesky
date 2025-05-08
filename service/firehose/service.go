package firehose

import (
	"bytes"
	"fmt"
	"github.com/awakari/int-bluesky/model"
	"github.com/fxamacker/cbor/v2"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-car"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	basicnode "github.com/ipld/go-ipld-prime/node/basic"
	"github.com/multiformats/go-multihash"
	"github.com/reiver/go-bsky/firehose"
	"io"
	"reflect"
	"strings"
)

const prefixDid = "did:plc:"
const typeKey = "$type"
const typeValPost = "app.bsky.feed.post"
const fmtUserId = "https://bsky.app/profile/%s"
const fmtUrlPost = "https://bsky.app/profile/%s/post%s"

func DecodePost(data []byte) (raw map[string]any, userId string, err error) {

	raw = make(map[string]any)

	msg := firehose.Message(data)
	var header firehose.MessageHeader
	var payload firehose.MessagePayload
	if err = msg.Decode(&header, &payload); err != nil {
		err = fmt.Errorf("firehose: failed to decode message: %s", err)
	}

	if err == nil {

		repo, repoOk := payload["repo"]
		switch repoOk {
		case true:
			var repoStr string
			repoStr, repoOk = repo.(string)
			switch repoOk {
			case true:
				switch strings.HasPrefix(repoStr, prefixDid) {
				case true:
					userId = fmt.Sprintf(fmtUserId, repo)
					raw[model.CeKeySubject] = userId
				default:
					return
				}
			default:
				return
			}
		default:
			return
		}

		var postCid cid.Cid
		ops, opsOk := payload["ops"]
		switch opsOk {
		case true:
			var rKey string
			for _, opRaw := range ops.([]any) {

				op, ok := opRaw.(map[any]any)
				if !ok {
					err = fmt.Errorf("failed to cast op to map: %s", reflect.TypeOf(opRaw))
					break
				}
				if op["action"] != "create" {
					continue
				}

				postCidRaw, postCidOk := op["cid"]
				if !postCidOk {
					err = fmt.Errorf("no cid key in op")
				}
				var postCidTag cbor.Tag
				postCidTag, postCidOk = postCidRaw.(cbor.Tag)
				if !postCidOk {
					err = fmt.Errorf("failed to cast cid to tag: %s", reflect.TypeOf(postCidRaw))
					break
				}
				cidContent, cidContentOk := postCidTag.Content.([]byte)
				if !cidContentOk {
					err = fmt.Errorf("failed to cast cid content to bytes: %s", reflect.TypeOf(postCidTag.Content))
					break
				}
				var h []byte
				h, err = multihash.Encode(cidContent, multihash.SHA2_256)
				if err != nil {
					err = fmt.Errorf("failed to encode the cid multihash: %w", err)
					break
				}
				postCid = cid.NewCidV1(cid.DagCBOR, h)
				raw[model.CeKeyCid] = postCid.String()

				postPath := op["path"]
				if pathStr, pathOk := postPath.(string); pathOk && strings.HasPrefix(pathStr, typeValPost) {
					rKey = strings.TrimPrefix(pathStr, typeValPost)
					break
				}
			}
			switch rKey {
			case "":
				return
			default:
				raw[model.CeKeyBlueskyRKey] = rKey
				raw[model.CeKeyObjectUrl] = fmt.Sprintf(fmtUrlPost, repo, rKey)
			}
		default:
			return
		}
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
