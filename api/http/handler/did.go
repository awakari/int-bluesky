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
		"@context": []string{
			"https://www.w3.org/ns/did/v1",
		},
		"id": h.Id,
		// may be better to have a separate generator per interest feed
		"service": []gin.H{
			{
				"id":              "#awakari-feed-generator",
				"type":            "BskyFeedGenerator",
				"serviceEndpoint": h.ServiceEndpoint,
			},
		},
	})
}
