package bluesky

import (
	"context"
	"fmt"
	"github.com/awakari/int-bluesky/util"
	"github.com/bluesky-social/indigo/api/bsky"
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

func (l logging) Login(ctx context.Context, id, password string) (did, token string, err error) {
	did, token, err = l.svc.Login(ctx, id, password)
	l.log.Log(ctx, util.LogLevel(err), fmt.Sprintf("service.Login(id=%s): did=%s, token=%s, %s", id, did, token, err))
	return
}

func (l logging) CreatePost(ctx context.Context, post *bsky.FeedPost, did, token string) (uri string, err error) {
	uri, err = l.svc.CreatePost(ctx, post, did, token)
	l.log.Log(ctx, util.LogLevel(err), fmt.Sprintf("service.CreatePost(%s, %s): %s", did, post.Text, err))
	return
}
