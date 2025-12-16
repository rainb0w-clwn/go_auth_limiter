package requestid

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const XRequestIDKey = "x-request-id"

func New() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		requestID := getRequestID(ctx)
		ctx = metadata.AppendToOutgoingContext(ctx, XRequestIDKey, requestID)
		return handler(ctx, req)
	}
}

func getRequestID(ctx context.Context) string {
	requestID := getStringFromContext(ctx, XRequestIDKey)
	if requestID == "" {
		return newRequestID()
	}
	return requestID
}

func newRequestID() string {
	return uuid.NewString()
}

func getStringFromContext(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	header, ok := md[key]
	if !ok || len(header) == 0 {
		return ""
	}
	return header[0]
}
