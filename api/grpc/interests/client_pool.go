package interests

import (
	"context"
	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
)

type clientPool struct {
	connPool *grpcpool.Pool
}

func NewClientPool(connPool *grpcpool.Pool) ServiceClient {
	return clientPool{
		connPool: connPool,
	}
}

func (cp clientPool) Search(ctx context.Context, req *SearchRequest, opts ...grpc.CallOption) (resp *SearchResponse, err error) {
	var conn *grpcpool.ClientConn
	conn, err = cp.connPool.Get(ctx)
	if err == nil {
		defer conn.Close()
	}
	var client ServiceClient
	if err == nil {
		client = NewServiceClient(conn)
		resp, err = client.Search(ctx, req, opts...)
	}
	return
}
