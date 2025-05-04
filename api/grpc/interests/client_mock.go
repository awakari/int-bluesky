package interests

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type clientMock struct {
}

func NewClientMock() ServiceClient {
	return clientMock{}
}

func (cm clientMock) Search(ctx context.Context, req *SearchRequest, opts ...grpc.CallOption) (resp *SearchResponse, err error) {
	resp = &SearchResponse{}
	switch req.Cursor.Id {
	case "fail":
		err = status.Error(codes.Internal, "internal failure")
	case "fail_auth":
		err = status.Error(codes.Unauthenticated, "auth failure")
	default:
		switch req.Order {
		case Order_DESC:
			resp.Ids = []string{
				"sub1",
				"sub0",
			}
		default:
			resp.Ids = []string{
				"sub0",
				"sub1",
			}
		}
	}
	return
}
