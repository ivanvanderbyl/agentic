package dstar

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type BuilderFunc[R any] func(client grpc.ClientConnInterface) *R

func NewQdrantClient[R any](addr string, clientBuilder BuilderFunc[R]) (*R, error) {
	config := new(tls.Config)

	interceptor := newAuthInterceptor(os.Getenv("QDRANT_API_KEY"))

	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithUnaryInterceptor(interceptor),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to qdrant")
	}

	collections_client := clientBuilder(conn)
	return collections_client, nil
}

func newAuthInterceptor(apiKey string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		newCtx := metadata.AppendToOutgoingContext(ctx, "api-key", apiKey)
		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}
