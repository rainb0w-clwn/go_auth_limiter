package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/meshapi/grpc-api-gateway/gateway"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/config"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/interfaces"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/server/http/health"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/server/http/log"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/server/http/requestid"
	proto "github.com/rainb0w-clwn/go_auth_limiter/proto/limiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Options struct {
	Host, Port                                                string
	ReadTimeout, WriteTimeout, IdleTimeout, ReadHeaderTimeout time.Duration
}

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type server struct {
	*http.Server
	logger appinterfaces.Logger
}

func New(options Options, logger appinterfaces.Logger) Server {
	serverHTTP := &http.Server{
		Addr:              net.JoinHostPort(options.Host, options.Port),
		ReadTimeout:       options.ReadTimeout,
		WriteTimeout:      options.WriteTimeout,
		IdleTimeout:       options.IdleTimeout,
		ReadHeaderTimeout: options.ReadHeaderTimeout,
	}
	return &server{
		serverHTTP,
		logger,
	}
}

func (s *server) Start(ctx context.Context) error {
	cfg := config.GetFromContext(ctx)
	if cfg == nil {
		return config.ErrNoConfigInContext
	}

	conn, err := grpc.NewClient(
		net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	mux := gateway.NewServeMux()
	for _, f := range []func(context.Context, *gateway.ServeMux, *grpc.ClientConn){
		proto.RegisterAuthLimiterHandler,
	} {
		f(ctx, mux, conn)
	}
	mux.Handle("GET", "/health", health.New())
	s.Handler = requestid.New(log.New(s.logger, mux))

	err = s.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *server) Stop(ctx context.Context) error {
	err := s.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
