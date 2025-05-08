package converter

import (
	"context"
	"fmt"
	"github.com/awakari/int-bluesky/util"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"log/slog"
	"time"
)

type serviceLogging struct {
	svc Service
	log *slog.Logger
}

func NewServiceLogging(svc Service, log *slog.Logger) Service {
	return serviceLogging{
		svc: svc,
		log: log,
	}
}

func (s serviceLogging) EventToPost(ctx context.Context, evt *pb.CloudEvent, interestId string, t *time.Time) (post *bsky.FeedPost, err error) {
	post, err = s.svc.EventToPost(ctx, evt, interestId, t)
	s.log.Log(ctx, util.LogLevel(err), fmt.Sprintf("converter.EventToPost(%s, %s): %s", evt.Id, interestId, err))
	return
}

func (s serviceLogging) PostToEvent(ctx context.Context, src *pb.CloudEvent) (dst *pb.CloudEvent, userId string, err error) {
	dst, userId, err = s.svc.PostToEvent(ctx, src)
	s.log.Log(ctx, util.LogLevel(err), fmt.Sprintf("converter.PostToEvent(%s): %s, %s", src.Id, userId, err))
	return
}
