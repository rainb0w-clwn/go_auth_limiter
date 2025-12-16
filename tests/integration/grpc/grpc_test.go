package grpc_test

import (
	"context"
	"os"
	"testing"

	proto "github.com/rainb0w-clwn/go_auth_limiter/proto/limiter"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var grpcAddr = "app:" + os.Getenv("GRPC_PORT")

func grpcClient(t *testing.T) proto.AuthLimiterClient {
	t.Helper()

	conn, err := grpc.NewClient(
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = conn.Close()
	})

	return proto.NewAuthLimiterClient(conn)
}

func ctx() context.Context {
	return context.Background()
}

func TestWhiteListAdd_OK(t *testing.T) {
	client := grpcClient(t)

	_, err := client.WhiteListAdd(ctx(), &proto.WhiteListAddRequest{
		IpNet: "10.10.10.0/24",
	})

	require.NoError(t, err)
}

func TestWhiteListAdd_InvalidArgument(t *testing.T) {
	client := grpcClient(t)

	_, err := client.WhiteListAdd(ctx(), &proto.WhiteListAddRequest{
		IpNet: "",
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestWhiteListDelete_NotFound(t *testing.T) {
	client := grpcClient(t)

	_, err := client.WhiteListDelete(ctx(), &proto.WhiteListDeleteRequest{
		IpNet: "192.168.100.0/24",
	})

	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestWhiteListDelete_InvalidArgument(t *testing.T) {
	client := grpcClient(t)

	_, err := client.WhiteListDelete(ctx(), &proto.WhiteListDeleteRequest{
		IpNet: "",
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestBlackListAdd_OK(t *testing.T) {
	client := grpcClient(t)

	_, err := client.BlackListAdd(ctx(), &proto.BlackListAddRequest{
		IpNet: "172.16.0.0/16",
	})

	require.NoError(t, err)
}

func TestBlackListAdd_InvalidArgument(t *testing.T) {
	client := grpcClient(t)

	_, err := client.BlackListAdd(ctx(), &proto.BlackListAddRequest{
		IpNet: "invalid-ip",
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestBlackListDelete_NotFound(t *testing.T) {
	client := grpcClient(t)

	_, err := client.BlackListDelete(ctx(), &proto.BlackListDeleteRequest{
		IpNet: "8.8.8.0/24",
	})

	require.Error(t, err)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func TestBlackListDelete_InvalidArgument(t *testing.T) {
	client := grpcClient(t)

	_, err := client.BlackListDelete(ctx(), &proto.BlackListDeleteRequest{
		IpNet: "",
	})

	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestBucketReset_OK(t *testing.T) {
	client := grpcClient(t)

	resp, err := client.LimitCheck(ctx(), &proto.LimitCheckRequest{
		Login:    "user1",
		Password: "secret",
		Ip:       "127.0.0.1",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	_, err = client.BucketReset(ctx(), &proto.BucketResetRequest{
		Login: "user1",
		Ip:    "127.0.0.1",
	})

	require.NoError(t, err)
}

func TestBucketReset_InvalidArgument(t *testing.T) {
	client := grpcClient(t)

	tests := []struct {
		name string
		req  *proto.BucketResetRequest
	}{
		{
			name: "empty login",
			req: &proto.BucketResetRequest{
				Login: "",
				Ip:    "127.0.0.1",
			},
		},
		{
			name: "empty ip",
			req: &proto.BucketResetRequest{
				Login: "user",
				Ip:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.BucketReset(ctx(), tt.req)

			require.Error(t, err)
			require.Equal(t, codes.InvalidArgument, status.Code(err))
		})
	}
}

func TestLimitCheck_OK(t *testing.T) {
	client := grpcClient(t)

	resp, err := client.LimitCheck(ctx(), &proto.LimitCheckRequest{
		Login:    "user1",
		Password: "secret",
		Ip:       "127.0.0.1",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestLimitCheck_InvalidArgument(t *testing.T) {
	client := grpcClient(t)

	tests := []struct {
		name string
		req  *proto.LimitCheckRequest
	}{
		{
			name: "empty login",
			req: &proto.LimitCheckRequest{
				Login:    "",
				Password: "pass",
				Ip:       "127.0.0.1",
			},
		},
		{
			name: "empty password",
			req: &proto.LimitCheckRequest{
				Login:    "user",
				Password: "",
				Ip:       "127.0.0.1",
			},
		},
		{
			name: "empty ip",
			req: &proto.LimitCheckRequest{
				Login:    "user",
				Password: "pass",
				Ip:       "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.LimitCheck(ctx(), tt.req)

			require.Error(t, err)
			require.Equal(t, codes.InvalidArgument, status.Code(err))
		})
	}
}
