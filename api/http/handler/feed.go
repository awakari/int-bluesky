package handler

import (
	"fmt"
	"github.com/awakari/int-bluesky/api/grpc/interests"
	"github.com/awakari/int-bluesky/api/http/bluesky"
	"github.com/awakari/int-bluesky/model"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type FeedHandler struct {
	DidWeb       string
	SvcBluesky   bluesky.Service
	DidPlc       string
	Token        string
	SvcInterests interests.Service
	UrlPrivacy   string
	UrlTos       string
}

const maxId = "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
const fmtFeedPrefix = "at://%s/app.bsky.feed.generator/"
const pageSizeDefault = 100
const loopLimit = 10

func (h FeedHandler) DescribeFeedGenerator(ctx *gin.Context) {
	q := &interests.Query{
		Sort:   interests.Sort_FOLLOWERS,
		Order:  interests.Order_DESC,
		Limit:  pageSizeDefault,
		Public: true,
	}
	cursor := &interests.Cursor{
		Id:        maxId,
		Followers: math.MaxInt64,
	}
	idsPage, err := h.SvcInterests.Search(ctx, model.GroupIdDefault, model.UserIdDefault, q, cursor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to list feeds")
		return
	}

	var feeds []*bsky.FeedDescribeFeedGenerator_Feed
	for _, id := range idsPage {
		feeds = append(feeds, &bsky.FeedDescribeFeedGenerator_Feed{
			Uri: h.feedUrl(id),
		})
	}
	ctx.JSON(http.StatusOK, bsky.FeedDescribeFeedGenerator_Output{
		Did:   h.DidWeb,
		Feeds: feeds,
		Links: &bsky.FeedDescribeFeedGenerator_Links{
			PrivacyPolicy:  &h.UrlPrivacy,
			TermsOfService: &h.UrlTos,
		},
	})
}

func (h FeedHandler) Skeleton(ctx *gin.Context) {

	feedUrl := ctx.Query("feed")
	feedPrefix := fmt.Sprintf(fmtFeedPrefix, h.DidPlc)
	if !strings.HasPrefix(feedUrl, feedPrefix) {
		ctx.JSON(http.StatusBadRequest, "invalid feed url")
	}
	interestId := strings.TrimPrefix(feedUrl, feedPrefix)
	cursor := ctx.Query("cursor")
	limitRaw := ctx.Query("limit")
	limit, err := strconv.Atoi(limitRaw)
	if err != nil || limit <= 0 {
		limit = pageSizeDefault
	}

	var urls []string
	var urlsNext []string
	for i := 0; i < loopLimit; i++ {
		urlsNext, cursor, err = h.SvcBluesky.Posts(ctx, h.DidPlc, h.Token, interestId, cursor, limit-len(urls))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "failed to list posts")
			return
		}
		urls = append(urls, urlsNext...)
		if cursor == "" {
			break
		}
		if len(urlsNext) > 0 {
			break
		}
	}

	feed := make([]*bsky.FeedDefs_SkeletonFeedPost, 0)
	for _, u := range urls {
		feed = append(feed, &bsky.FeedDefs_SkeletonFeedPost{
			FeedContext: &interestId,
			Post:        u,
		})
	}

	ctx.JSON(http.StatusOK, bsky.FeedGetFeedSkeleton_Output{
		Cursor: &cursor,
		Feed:   feed,
	})
}

func (h FeedHandler) feedUrl(interestId string) (feedUrl string) {
	feedUrl = fmt.Sprintf("at://%s/app.bsky.feed.generator/%s", h.DidPlc, interestId)
	return
}
