package interests

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Service interface {
	// Search returns all subscription ids matching the query.
	Search(ctx context.Context, groupId, userId string, q *Query, cursor *Cursor) (ids []string, err error)
}

type service struct {
	client ServiceClient
}

const keyGroupId = "x-awakari-group-id"
const keyUserId = "x-awakari-user-id"

// ErrNotFound indicates the interest is missing in the storage and can not be read/updated/deleted.
var ErrNotFound = errors.New("interest not found")

// ErrInternal indicates some unexpected internal failure.
var ErrInternal = errors.New("internal failure")

// ErrInvalid indicates the invalid request.
var ErrInvalid = errors.New("invalid request")

var ErrUnavailable = errors.New("unavailable")

var ErrAuth = errors.New("authentication failure")

func NewService(client ServiceClient) Service {
	return service{
		client: client,
	}
}

func (svc service) Search(ctx context.Context, groupId, userId string, q *Query, cursor *Cursor) (ids []string, err error) {
	ctx = metadata.AppendToOutgoingContext(ctx, keyGroupId, groupId, keyUserId, userId)
	req := SearchRequest{
		Cursor: cursor,
	}
	if q != nil {
		req.Limit = q.Limit
		req.Order = q.Order
		req.Pattern = q.Pattern
		req.Sort = q.Sort
	}
	var resp *SearchResponse
	resp, err = svc.client.Search(ctx, &req)
	if resp != nil {
		ids = resp.Ids
	}
	err = decodeError(err)
	return
}

func decodeError(src error) (dst error) {
	switch {
	case src == nil:
	default:
		s, isGrpcErr := status.FromError(src)
		switch {
		case !isGrpcErr:
			dst = src
		case s.Code() == codes.OK:
		case s.Code() == codes.NotFound:
			dst = fmt.Errorf("%w: %s", ErrNotFound, src)
		case s.Code() == codes.InvalidArgument:
			dst = fmt.Errorf("%w: %s", ErrInvalid, src)
		case s.Code() == codes.Unauthenticated:
			dst = fmt.Errorf("%w: %s", ErrAuth, src)
		case s.Code() == codes.Unavailable:
			dst = fmt.Errorf("%w: %s", ErrUnavailable, src)
		default:
			dst = fmt.Errorf("%w: %s", ErrInternal, src)
		}
	}
	return
}
