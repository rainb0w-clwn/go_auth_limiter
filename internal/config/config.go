package config

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/creasty/defaults"
	gocfg "github.com/dsbasko/go-cfg"
	"gopkg.in/yaml.v3"
)

type key int

const (
	ctxKey key = iota
)

var ErrNoConfigInContext = errors.New("no config found in context")

type Config struct {
	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL" default:"debug"`
	} `yaml:"log"`
	HTTP struct {
		Host              string        `yaml:"host" env:"HTTP_HOST"`
		Port              string        `yaml:"port" env:"HTTP_PORT"`
		ReadTimeout       time.Duration `default:"3s" yaml:"readTimeout" env:"HTTP_READ_TIMEOUT"`
		WriteTimeout      time.Duration `default:"3s" yaml:"writeTimeout" env:"HTTP_WRITE_TIMEOUT"`
		IdleTimeout       time.Duration `default:"30s" yaml:"idleTimeout" env:"HTTP_IDLE_TIMEOUT"`
		ReadHeaderTimeout time.Duration `default:"2s" yaml:"readHeaderTimeout" env:"HTTP_READ_HEADER_TIMEOUT"`
	} `yaml:"http"`
	GRPC struct {
		Host           string        `yaml:"host" env:"GRPC_HOST"`
		Port           string        `yaml:"port" env:"GRPC_PORT"`
		ConnectTimeout time.Duration `default:"5s" yaml:"connectTimeout" env:"GRPC_CONNECT_TIMEOUT"`
	} `yaml:"grpc"`
	DB struct {
		Type          string `yaml:"type" env:"DB_TYPE"`
		DSN           string `yaml:"dsn" env:"DB_DSN"`
		MigrationsDir string `default:"./migrations/" yaml:"migrationsDir" env:"DB_MIGRATIONS_DIR"`
		AutoMigrate   bool   `default:"true" yaml:"autoMigrate" env:"DB_AUTO_MIGRATE"`
	} `yaml:"db"`
	App struct {
		RefillRate struct {
			Count int           `default:"3" yaml:"count" env:"APP_REFILL_RATE_COUNT"`
			Time  time.Duration `default:"30s" yaml:"time" env:"APP_REFILL_RATE_TIME"`
		} `yaml:"refillRate"`
		GarbageCollector struct {
			Enabled  bool          `default:"true" yaml:"app" env:"APP_ENABLE_GARBAGE_COLLECTOR"`
			TTL      time.Duration `default:"600s" yaml:"ttl" env:"APP_TTL"`
			Interval time.Duration `default:"60s" yaml:"interval" env:"APP_INTERVAL"`
		} `yaml:"garbageCollector"`
	} `yaml:"app"`
}

func New(r io.Reader) (*Config, error) {
	config := &Config{}
	if err := yaml.NewDecoder(r).Decode(config); err != nil {
		return nil, err
	}
	if err := gocfg.ReadEnv(config); err != nil {
		return nil, err
	}
	if err := defaults.Set(config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey, c)
}

func GetFromContext(ctx context.Context) *Config {
	return ctx.Value(ctxKey).(*Config)
}
