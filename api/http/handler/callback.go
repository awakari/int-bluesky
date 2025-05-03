package handler

import (
	"fmt"
	"github.com/awakari/int-bluesky/api/http/bluesky"
	"github.com/awakari/int-bluesky/api/http/reader"
	"github.com/awakari/int-bluesky/config"
	"github.com/awakari/int-bluesky/service/converter"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/utf8"
	ceProto "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	ce "github.com/cloudevents/sdk-go/v2/event"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type CallbackHandler interface {
	Confirm(ctx *gin.Context)
	Deliver(ctx *gin.Context)
}

type callbackHandler struct {
	topicPrefixBase string
	host            string
	cfgEvtType      config.EventTypeConfig
	svcConv         converter.Service
	svcBluesky      bluesky.Service
	blueskyDid      string
	blueskyToken    string
}

const keyHubChallenge = "hub.challenge"
const keyHubTopic = "hub.topic"
const linkSelfSuffix = ">; rel=\"self\""
const keyAckCount = "X-Ack-Count"

func NewCallbackHandler(
	topicPrefixBase, host string,
	cfgEvtType config.EventTypeConfig,
	svcConv converter.Service,
	svcBluesky bluesky.Service,
	blueskyDid string,
	blueskyToken string,
) CallbackHandler {
	return callbackHandler{
		topicPrefixBase: topicPrefixBase,
		host:            host,
		cfgEvtType:      cfgEvtType,
		svcConv:         svcConv,
		svcBluesky:      svcBluesky,
		blueskyDid:      blueskyDid,
		blueskyToken:    blueskyToken,
	}
}

func (ch callbackHandler) Confirm(ctx *gin.Context) {
	topic := ctx.Query(keyHubTopic)
	challenge := ctx.Query(keyHubChallenge)
	if strings.HasPrefix(topic, ch.topicPrefixBase+"/sub/"+reader.FmtJson) {
		ctx.String(http.StatusOK, challenge)
	} else {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("invalid topic: %s", topic))
	}
	return
}

func (ch callbackHandler) Deliver(ctx *gin.Context) {

	var topic string
	for k, vals := range ctx.Request.Header {
		if strings.ToLower(k) == "link" {
			for _, l := range vals {
				if strings.HasSuffix(l, linkSelfSuffix) && len(l) > len(linkSelfSuffix) {
					topic = l[1 : len(l)-len(linkSelfSuffix)]
				}
			}
		}
	}
	if topic == "" {
		ctx.String(http.StatusBadRequest, "self link header missing in the request")
		return
	}

	var interestId string
	topicParts := strings.Split(topic, "/")
	topicPartsLen := len(topicParts)
	if topicPartsLen > 0 {
		interestId = topicParts[topicPartsLen-1]
	}
	if interestId == "" {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("invalid self link header value in the request: %s", topic))
		return
	}

	defer ctx.Request.Body.Close()
	var evts []*ce.Event
	if err := sonic.ConfigDefault.NewDecoder(ctx.Request.Body).Decode(&evts); err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("failed to deserialize the request payload: %s", err))
		return
	}

	var countDelivered uint64
	var err error
	for _, evt := range evts {
		var evtProto *pb.CloudEvent
		evtProto, err = ceProto.ToProto(evt)
		var dataTxt string
		if err == nil {
			err = evt.DataAs(&dataTxt)
		}
		if err == nil && utf8.ValidateString(dataTxt) {
			evtProto.Data = &pb.CloudEvent_TextData{
				TextData: dataTxt,
			}
		}
		if err == nil {
			switch evtProto.Type {
			case ch.cfgEvtType.InterestsUpdated:
				// TODO
				fallthrough
			default:
				var post *bsky.FeedPost
				var errNotify error
				post, errNotify = ch.svcConv.EventToPost(ctx, evtProto, interestId, nil)
				if errNotify == nil {
					_, err = ch.svcBluesky.CreatePost(ctx, post, ch.blueskyDid, ch.blueskyToken)
				}
			}
		}
		if err != nil {
			break
		}
		countDelivered++
	}

	ctx.Writer.Header().Add(keyAckCount, strconv.FormatUint(countDelivered, 10))
	switch {
	case countDelivered < 1 && err != nil:
		ctx.String(http.StatusInternalServerError, err.Error())
	default:
		ctx.Status(http.StatusOK)
	}

	return
}
