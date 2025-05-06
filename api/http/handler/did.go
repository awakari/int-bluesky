package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type DidHandler struct {
	Id              string
	ServiceEndpoint string
}

func (h DidHandler) Handle(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"id": h.Id,
		// may be better to have a separate generator per interest feed
		"service": []gin.H{
			{
				"id":              "#bsky_fg",
				"type":            "BskyFeedGenerator",
				"serviceEndpoint": h.ServiceEndpoint,
			},
		},
	})
}
