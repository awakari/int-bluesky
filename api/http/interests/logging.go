package interests

import (
	"context"
	"fmt"
	"github.com/awakari/int-bluesky/model"
	"github.com/awakari/int-bluesky/util"
	"log/slog"
)

type logging struct {
	svc Service
	log *slog.Logger
}

func NewLogging(svc Service, log *slog.Logger) Service {
	return logging{
		svc: svc,
		log: log,
	}
}

func (l logging) Read(ctx context.Context, groupId, userId, interestId string) (d model.InterestData, err error) {
	d, err = l.svc.Read(ctx, groupId, userId, interestId)
	l.log.Log(ctx, util.LogLevel(err), fmt.Sprintf("interests.Read(%s, %s, %s): %v, %s", groupId, userId, interestId, d, err))
	return
}
