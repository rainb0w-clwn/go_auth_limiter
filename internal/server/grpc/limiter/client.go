package limiter

import (
	"net"

	proto "github.com/rainb0w-clwn/go_auth_limiter/proto/limiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	proto.AuthLimiterClient
	*grpc.ClientConn
}

func NewClient(host, port string) (*Client, error) {
	conn, err := grpc.NewClient(
		net.JoinHostPort(host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		proto.NewAuthLimiterClient(conn),
		conn,
	}, nil
}
