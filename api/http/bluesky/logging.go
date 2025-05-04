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
	l.log.Log(ctx, util.LogLevel(err), fmt.Sprintf("service.CreatePost(%s): %s", did, err))
	return
}

func (l logging) Posts(ctx context.Context, did, token, interestId, cursor string) (urls []string, next string, err error) {
	urls, next, err = l.svc.Posts(ctx, did, token, interestId, cursor)
	l.log.Log(ctx, util.LogLevel(err), fmt.Sprintf("service.Posts(%s, %s, %s): %d, %s, %s", did, interestId, cursor, len(urls), next, err))
	return
}

func (l logging) CreateFeed(ctx context.Context, didWeb, didPlc, token, interestId string) (err error) {
	err = l.svc.CreateFeed(ctx, didWeb, didPlc, token, interestId)
	l.log.Log(ctx, util.LogLevel(err), fmt.Sprintf("service.CreateFeed(%s, %s, %s): %s", didWeb, didPlc, interestId, err))
	return
}
