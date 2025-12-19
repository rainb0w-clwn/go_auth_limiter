package app

import (
	"context"
	"time"

	"github.com/rainb0w-clwn/go_auth_limiter/internal/bucket/refillrate"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/config"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/interfaces"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter/auth"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter/composite"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter/tokenbucket/gb"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/rule"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/storage/postgres"
)

type App struct {
	rule    rule.IService
	limiter limiter.IService

	logger appinterfaces.Logger
	config *config.Config
}

func New(ctx context.Context, config *config.Config, logger appinterfaces.Logger) (appinterfaces.Application, error) {
	postgresStorage := postgres.New()

	if err := postgresStorage.Connect(ctx); err != nil {
		return nil, err
	}
	ruleStorage := rule.NewStorage(postgresStorage)
	ruleService := rule.NewService(ruleStorage)

	limitStorage := limiter.NewStorage(postgresStorage)
	bucketLimiter := composite.New(
		limitStorage,
		refillrate.New(config.App.RefillRate.Count, config.App.RefillRate.Time),
	)
	limiterService := auth.New(ruleService, bucketLimiter)

	// Init Limiter Garbage Collector
	limiterGB := gb.New(bucketLimiter, config.App.GarbageCollector.TTL)

	if config.App.GarbageCollector.Enabled {
		go func() {
			for {
				select {
				case <-ctx.Done():
					logger.Info("GB finished.")

					return
				case <-time.After(config.App.GarbageCollector.Interval):
					logger.Info("GB sweeping..")

					err := limiterGB.Sweep()
					if err != nil {
						logger.Error("GB error", "error", err)
					}
				}
			}
		}()
	}

	return &App{
		rule:    ruleService,
		limiter: limiterService,

		logger: logger,
		config: config,
	}, nil
}

func (a *App) LimitCheck(ip, login, password string) (bool, error) {
	return a.limiter.SatisfyLimit(limiter.UserIdentityDto{
		limiter.IPLimit.String():       ip,
		limiter.LoginLimit.String():    login,
		limiter.PasswordLimit.String(): password,
	})
}

func (a *App) LimitReset(ip, login string) error {
	return a.limiter.ResetLimit(limiter.UserIdentityDto{
		limiter.IPLimit.String():    ip,
		limiter.LoginLimit.String(): login,
	})
}

func (a *App) WhiteListAdd(ip string) error {
	return a.rule.WhiteListAdd(ip)
}

func (a *App) WhiteListDelete(ip string) error {
	return a.rule.WhiteListDelete(ip)
}

func (a *App) BlackListAdd(ip string) error {
	return a.rule.BlackListAdd(ip)
}

func (a *App) BlackListDelete(ip string) error {
	return a.rule.BlackListDelete(ip)
}
