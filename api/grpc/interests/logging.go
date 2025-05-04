package interests

import (
	"context"
	"fmt"
	"github.com/awakari/int-bluesky/util"
	"log/slog"
)

type loggingMiddleware struct {
	svc Service
	log *slog.Logger
}

func NewLoggingMiddleware(svc Service, log *slog.Logger) Service {
	return loggingMiddleware{
		svc: svc,
		log: log,
	}
}

func (lm loggingMiddleware) Search(ctx context.Context, groupId, userId string, q *Query, cursor *Cursor) (ids []string, err error) {
	ids, err = lm.svc.Search(ctx, groupId, userId, q, cursor)
	lm.log.Log(ctx, util.LogLevel(err), fmt.Sprintf("interests.Search(groupId=%s, userId=%s, q=%+v, cursor=%+v): %d, err=%s", groupId, userId, q, cursor, len(ids), err))
	return
}
