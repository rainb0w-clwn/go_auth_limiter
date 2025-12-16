package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/app"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/config"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/logger"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/server"
	serverGRPC "github.com/rainb0w-clwn/go_auth_limiter/internal/server/grpc"
	serverHTTP "github.com/rainb0w-clwn/go_auth_limiter/internal/server/http"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/service/limiter"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/configs/config.yml", "Path to configuration file")
}

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	configR, err := os.Open(configFile)
	if err != nil {
		log.Printf("%s", "Error opening config file: "+err.Error())
		return 1
	}

	cfg, err := config.New(configR)
	if err != nil {
		log.Printf("%s", "Error parsing config file: "+err.Error())
		return 1
	}
	ctx = cfg.WithContext(ctx)

	logg := logger.New(logger.GetLevelOrPanic(cfg.Log.Level), os.Stdout)

	application, err := app.New(ctx, cfg, logg)
	if err != nil {
		logg.Error("Error init application", err)
		return 1
	}

	srv := server.New(
		server.Options{
			GRPC: serverGRPC.Options{
				Host:           cfg.GRPC.Host,
				Port:           cfg.GRPC.Port,
				ConnectTimeout: cfg.GRPC.ConnectTimeout,
			},
			HTTP: serverHTTP.Options{
				Host:         cfg.HTTP.Host,
				Port:         cfg.HTTP.Port,
				ReadTimeout:  cfg.HTTP.ReadTimeout,
				WriteTimeout: cfg.HTTP.WriteTimeout,
			},
		},
		logg,
		application,
	)

	service := limiter.New(srv, logg)
	err = service.Run(ctx)
	if err != nil {
		logg.Error("Error starting calendar: %v", err)
		return 1
	}

	return 0
}
