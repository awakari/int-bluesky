package handler

import (
	"fmt"
	"github.com/awakari/int-bluesky/api/grpc/interests"
	"github.com/awakari/int-bluesky/api/http/bluesky"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
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
const groupId = "default"
const userId = "public"

func (h FeedHandler) DescribeFeedGenerator(ctx *gin.Context) {
	q := &interests.Query{
		Sort:   interests.Sort_FOLLOWERS,
		Order:  interests.Order_DESC,
		Limit:  100,
		Public: true,
	}
	cursor := &interests.Cursor{
		Id:        maxId,
		Followers: math.MaxInt64,
	}
	idsPage, err := h.SvcInterests.Search(ctx, groupId, userId, q, cursor)
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

	interestId := ctx.Query("feed")
	cursor := ctx.Query("cursor")
	urls, next, err := h.SvcBluesky.Posts(ctx, h.DidPlc, h.Token, interestId, cursor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to list posts")
		return
	}

	var feed []*bsky.FeedDefs_SkeletonFeedPost
	for _, u := range urls {
		feed = append(feed, &bsky.FeedDefs_SkeletonFeedPost{
			Post: u,
		})
	}

	ctx.JSON(http.StatusOK, bsky.FeedGetFeedSkeleton_Output{
		Cursor: &next,
		Feed:   feed,
	})
}

func (h FeedHandler) feedUrl(interestId string) (feedUrl string) {
	feedUrl = fmt.Sprintf("at://%s/app.bsky.feed.generator/%s", h.DidWeb, interestId)
	return
}
