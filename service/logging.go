package service

import (
	"context"
	"fmt"
	"github.com/awakari/int-bluesky/util"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"log/slog"
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

func (sl serviceLogging) ConsumeInterestEvents(ctx context.Context, evts []*pb.CloudEvent) (err error) {
	err = sl.svc.ConsumeInterestEvents(ctx, evts)
	sl.log.Log(ctx, util.LogLevel(err), fmt.Sprintf("service.ConsumeInterestEvents(%d): err=%s", len(evts), err))
	return
}

func (sl serviceLogging) ConsumePostEvents(ctx context.Context, evts []*pb.CloudEvent) (err error) {
	err = sl.svc.ConsumePostEvents(ctx, evts)
	sl.log.Log(ctx, util.LogLevel(err), fmt.Sprintf("service.ConsumePostEvents(%d): err=%s", len(evts), err))
	return
}
