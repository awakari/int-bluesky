package interests

import (
	"context"
	"errors"
	"fmt"
	apiGrpc "github.com/awakari/int-bluesky/api/grpc/interests"
	"github.com/awakari/int-bluesky/model"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
)

type Service interface {
	Read(ctx context.Context, groupId, userId, subId string) (d model.InterestData, err error)
}

type service struct {
	clientHttp *http.Client
	url        string
	token      string
}

var protoJsonUnmarshalOpts = protojson.UnmarshalOptions{
	DiscardUnknown: true,
	AllowPartial:   true,
}

var ErrNoAuth = errors.New("unauthenticated request")
var ErrNotFound = errors.New("interest not found")

func NewService(clientHttp *http.Client, url, token string) Service {
	return service{
		clientHttp: clientHttp,
		url:        url,
		token:      token,
	}
}

func (svc service) Read(ctx context.Context, groupId, userId, interestId string) (d model.InterestData, err error) {

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, svc.url+"/"+interestId, nil)

	var resp *http.Response
	if err == nil {
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "Bearer "+svc.token)
		req.Header.Add(model.KeyGroupId, groupId)
		req.Header.Add(model.KeyUserId, userId)
		resp, err = svc.clientHttp.Do(req)
	}

	if err == nil {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			err = ErrNoAuth
		case http.StatusNotFound:
			err = fmt.Errorf("%w: %s", ErrNotFound, interestId)
		}
	}

	var respData []byte
	if err == nil {
		defer resp.Body.Close()
		respData, err = io.ReadAll(resp.Body)
	}

	var respProto apiGrpc.ReadResponse
	if err == nil {
		err = protoJsonUnmarshalOpts.Unmarshal(respData, &respProto)
	}

	if err == nil {
		d.Description = respProto.Description
		d.Enabled = respProto.Enabled
		d.Public = respProto.Public
		d.Followers = respProto.Followers
		if respProto.Expires != nil && respProto.Expires.IsValid() {
			d.Expires = respProto.Expires.AsTime()
		}
		if respProto.Created != nil && respProto.Created.IsValid() {
			d.Created = respProto.Created.AsTime()
		}
		if respProto.Updated != nil && respProto.Updated.IsValid() {
			d.Updated = respProto.Updated.AsTime()
		}
	}

	return
}
