package bluesky

import (
	"bytes"
	"context"
	"fmt"
	apiHttpInterests "github.com/awakari/int-bluesky/api/http/interests"
	"github.com/awakari/int-bluesky/model"
	"github.com/awakari/int-bluesky/util"
	"github.com/bluesky-social/indigo/api/agnostic"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bytedance/sonic"
	"io"
	"net/http"
	"time"
)

type Service interface {
	Login(ctx context.Context, id, password string) (did, token string, err error)
	CreatePost(ctx context.Context, post *bsky.FeedPost, did, token string) (uri string, err error)
	Posts(ctx context.Context, did, token, interestId, cursor string) (urls []string, next string, err error)
	CreateFeed(ctx context.Context, didWeb, didPlc, token, interestId string) (err error)
}

type service struct {
	clientHttp   *http.Client
	svcInterests apiHttpInterests.Service
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

type actorPostsResp struct {
	Feed   []feedItem `json:"feed"`
	Cursor string     `json:"cursor"`
}

type feedItem struct {
	Post struct {
		Uri    string `json:"uri"`
		Record struct {
			Labels *bsky.FeedPost_Labels `json:"labels,omitempty" cborgen:"labels,omitempty"`
		} `json:"record"`
	} `json:"post"`
}

const valContentTypeJson = "application/json"
const limitBodyLen = 1_048_576
const coll = "app.bsky.feed.post"

func NewService(clientHttp *http.Client, svcInterests apiHttpInterests.Service) Service {
	return service{
		clientHttp:   clientHttp,
		svcInterests: svcInterests,
	}
}

func (s service) Login(ctx context.Context, id, password string) (did, token string, err error) {
	body, _ := sonic.Marshal(loginReq{
		Identifier: id,
		Password:   password,
	})
	req, _ := http.NewRequestWithContext(
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
		respData, err = io.ReadAll(io.LimitReader(resp.Body, limitBodyLen))
	}
	var lr loginResp
	if err == nil {
		err = sonic.Unmarshal(respData, &lr)
	}
	switch err {
	case nil:
		did = lr.Did
		token = lr.AccessJwt
	default:
		err = fmt.Errorf("response: %d, %+v, %s", resp.StatusCode, resp.Header, string(respData))
	}
	return
}

func (s service) CreatePost(ctx context.Context, post *bsky.FeedPost, did, token string) (uri string, err error) {
	body, _ := sonic.Marshal(&createPostReq{
		Repo:       did,
		Collection: coll,
		Record:     post,
	})
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://bsky.social/xrpc/com.atproto.repo.createRecord",
		bytes.NewReader(body),
	)
	req.Header.Set("Accept", valContentTypeJson)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", valContentTypeJson)
	var resp *http.Response
	if resp, err = s.clientHttp.Do(req); err != nil {
		err = fmt.Errorf("request failure: %s", err)
	}
	var respData []byte
	if err == nil {
		defer resp.Body.Close()
		if respData, err = io.ReadAll(io.LimitReader(resp.Body, limitBodyLen)); err != nil {
			err = fmt.Errorf("response read failure: %s", err)
		}
	}
	var cr createPostResp
	if err == nil {
		if err = sonic.Unmarshal(respData, &cr); err != nil {
			err = fmt.Errorf("response unmarshal failure: %s, status: %d, headers: %+v, data: %s", err, resp.StatusCode, resp.Header, string(respData))
		}
	}
	if err == nil {
		uri = cr.Uri
	}
	return
}

func (s service) Posts(ctx context.Context, did, token, interestId, cursor string) (urls []string, next string, err error) {
	reqUrl := fmt.Sprintf("https://bsky.social/xrpc/app.bsky.feed.getAuthorFeed?actor=%s", did)
	if cursor != "" {
		reqUrl += "&cursor=" + cursor
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	req.Header.Set("Accept", valContentTypeJson)
	req.Header.Set("Authorization", "Bearer "+token)
	var resp *http.Response
	if resp, err = s.clientHttp.Do(req); err != nil {
		err = fmt.Errorf("request failure: %s", err)
	}
	var respData []byte
	if err == nil {
		defer resp.Body.Close()
		if respData, err = io.ReadAll(io.LimitReader(resp.Body, limitBodyLen)); err != nil {
			err = fmt.Errorf("response read failure: %s", err)
		}
	}
	if err == nil && resp.StatusCode > 299 {
		err = fmt.Errorf("response: %d, %s", resp.StatusCode, string(respData))
	}
	var all actorPostsResp
	if err == nil {
		if err = sonic.Unmarshal(respData, &all); err != nil {
			err = fmt.Errorf("response unmarshal failure: %s, status: %d, headers: %+v, data: %s", err, resp.StatusCode, resp.Header, string(respData))
		}
	}
	if err == nil {
		next = all.Cursor
		for _, p := range all.Feed {
			if p.Post.Record.Labels == nil {
				continue
			}
			if p.Post.Record.Labels.LabelDefs_SelfLabels == nil {
				continue
			}
			lblVals := p.Post.Record.Labels.LabelDefs_SelfLabels.Values
			if len(lblVals) < 1 {
				continue
			}
			if lblVals[0].Val != interestId {
				continue
			}
			urls = append(urls, p.Post.Uri)
		}
	}
	return
}

func (s service) CreateFeed(ctx context.Context, didWeb, didPlc, token, interestId string) (err error) {
	var interest model.InterestData
	if interest, err = s.svcInterests.Read(ctx, model.GroupIdDefault, model.UserIdDefault, interestId); err != nil {
		err = fmt.Errorf("failed to read the interest %s: %w", interestId, err)
	}
	var req *http.Request
	if err == nil {
		interestLnk := fmt.Sprintf("https://awakari.com/sub-details.html?id=%s", interestId)
		body, _ := sonic.Marshal(&agnostic.RepoPutRecord_Input{
			Collection: "app.bsky.feed.generator",
			Record: map[string]any{
				"$type":               "app.bsky.feed.generator",
				"acceptsInteractions": util.Ptr(true),
				"createdAt":           time.Now().UTC().Format(time.RFC3339),
				"did":                 didWeb, // ← feed generator server DID
				"displayName":         util.TruncateStringUtf8(interest.Description, 24),
				"description":         interestLnk,
				"descriptionFacets": []*bsky.RichtextFacet{{
					Features: []*bsky.RichtextFacet_Features_Elem{{
						RichtextFacet_Link: &bsky.RichtextFacet_Link{
							Uri: interestLnk,
						},
					}},
					Index: &bsky.RichtextFacet_ByteSlice{
						ByteEnd: int64(len(interestLnk)),
					},
				}},
			},
			Repo: didPlc,
			Rkey: interestId,
		})
		req, _ = http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			"https://bsky.social/xrpc/com.atproto.repo.putRecord",
			bytes.NewReader(body),
		)
		req.Header.Set("Accept", valContentTypeJson)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", valContentTypeJson)
	}
	var resp *http.Response
	if resp, err = s.clientHttp.Do(req); err != nil {
		err = fmt.Errorf("request failure: %s", err)
	}
	var respData []byte
	if err == nil {
		defer resp.Body.Close()
		if respData, err = io.ReadAll(io.LimitReader(resp.Body, limitBodyLen)); err != nil {
			err = fmt.Errorf("response read failure: %s", err)
		}
	}
	if resp.StatusCode > 299 {
		err = fmt.Errorf("response status: %d, %s", resp.StatusCode, string(respData))
	}
	return
}
