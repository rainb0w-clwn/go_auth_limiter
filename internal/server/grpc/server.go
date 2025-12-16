package server

import (
	"context"
	"net"
	"time"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/config"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/interfaces"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/server/grpc/log"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/server/grpc/requestid"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/server/grpc/validate"
	proto "github.com/rainb0w-clwn/go_auth_limiter/proto/limiter"
	"google.golang.org/grpc"
)

type Options struct {
	Host, Port     string
	ConnectTimeout time.Duration
}

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type server struct {
	*grpc.Server
	logger interfaces.Logger
}

func New(options Options, logger interfaces.Logger, app interfaces.Application) Server {
	serverGRPC := grpc.NewServer(
		grpc.ConnectionTimeout(options.ConnectTimeout),
		grpc.ChainUnaryInterceptor(
			requestid.New(),
			validate.New(),
			log.New(logger),
		),
	)
	proto.RegisterAuthLimiterServer(serverGRPC, NewService(app, logger))
	return &server{serverGRPC, logger}
}

func (s *server) Start(ctx context.Context) error {
	cfg := config.GetFromContext(ctx)
	if cfg == nil {
		return config.ErrNoConfigInContext
	}

	var lc net.ListenConfig

	listener, err := lc.Listen(
		ctx,
		"tcp",
		net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port),
	)
	if err != nil {
		return err
	}

	err = s.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}

func (s *server) Stop(_ context.Context) error {
	s.GracefulStop()
	return nil
}
