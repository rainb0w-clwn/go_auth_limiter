package validate

import (
	"context"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func New() grpc.UnaryServerInterceptor {
	v, _ := protovalidate.New()
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		r, ok := req.(proto.Message)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "Invalid message")
		}
		if err := v.Validate(r); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return handler(ctx, req)
	}
}
