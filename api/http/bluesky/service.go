package bluesky

import (
	"bytes"
	"context"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
	"net/http/httptest"
)

type Service interface {
	Login(ctx context.Context, id, password string) (did, token string, err error)
	CreatePost(ctx context.Context, post *bsky.FeedPost, did, token string) (uri string, err error)
}

type service struct {
	clientHttp *http.Client
}

type loginReq struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type loginResp struct {
	AccessJwt string `json:"accessJwt"`
	Did       string `json:"did"`
}

type createPostReq struct {
	Repo       string         `json:"repo"`
	Collection string         `json:"collection"`
	Record     *bsky.FeedPost `json:"record"`
}

type createPostResp struct {
	Uri string `json:"uri"`
	Cid string `json:"cid"`
}

const valContentTypeJson = "application/json"
const limitBodyLen = 262_144
const coll = "app.bsky.feed.post"

func NewService(clientHttp *http.Client) Service {
	return service{
		clientHttp: clientHttp,
	}
}

func (s service) Login(ctx context.Context, id, password string) (did, token string, err error) {
	body, _ := sonic.Marshal(loginReq{
		Identifier: id,
		Password:   password,
	})
	req := httptest.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://bsky.social/xrpc/com.atproto.server.createSession",
		bytes.NewReader(body),
	)
	req.Header.Set("Accept", valContentTypeJson)
	req.Header.Set("Content-Type", valContentTypeJson)
	var resp *http.Response
	resp, err = s.clientHttp.Do(req)
	var respData []byte
	if err == nil {
		defer resp.Body.Close()
		respData, err = io.ReadAll(io.LimitReader(req.Body, limitBodyLen))
	}
	var lr loginResp
	if err == nil {
		err = sonic.Unmarshal(respData, &lr)
	}
	if err == nil {
		did = lr.Did
		token = lr.AccessJwt
	}
	return
}

func (s service) CreatePost(ctx context.Context, post *bsky.FeedPost, did, token string) (uri string, err error) {
	body, _ := sonic.Marshal(&createPostReq{
		Repo:       did,
		Collection: coll,
		Record:     post,
	})
	req := httptest.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://bsky.social/xrpc/com.atproto.repo.createRecord",
		bytes.NewReader(body),
	)
	req.Header.Set("Accept", valContentTypeJson)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", valContentTypeJson)
	var resp *http.Response
	resp, err = s.clientHttp.Do(req)
	var respData []byte
	if err == nil {
		defer resp.Body.Close()
		respData, err = io.ReadAll(io.LimitReader(req.Body, limitBodyLen))

	}
	var cr createPostResp
	if err == nil {
		err = sonic.Unmarshal(respData, &cr)
	}
	if err == nil {
		uri = cr.Uri
	}
	return
}
