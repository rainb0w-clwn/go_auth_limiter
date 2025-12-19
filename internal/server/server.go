package server

import (
	"github.com/rainb0w-clwn/go_auth_limiter/internal/interfaces"
	serverGRPC "github.com/rainb0w-clwn/go_auth_limiter/internal/server/grpc"
	serverHTTP "github.com/rainb0w-clwn/go_auth_limiter/internal/server/http"
)

type Options struct {
	HTTP serverHTTP.Options
	GRPC serverGRPC.Options
}

type Server struct {
	GRPC    *serverGRPC.Server
	HTTP    *serverHTTP.Server
	logger  appinterfaces.Logger
	options Options
}

func New(options Options, logger appinterfaces.Logger, app appinterfaces.Application) *Server {
	grpcServer := serverGRPC.New(
		options.GRPC,
		logger,
		app,
	)
	httpServer := serverHTTP.New(
		options.HTTP,
		logger,
	)
	return &Server{
		GRPC:    &grpcServer,
		HTTP:    &httpServer,
		logger:  logger,
		options: options,
	}
}
