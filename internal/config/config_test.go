package config

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var yamlData = `
logger:
  level: "debug"
http:
  host: "0.0.0.0"
  port: 3000
  readTimeout: 5s
  writeTimeout: 10s
grpc:
  host: "0.0.0.0"
  port: 50051
  connectTimeout: 5s
db:
  type: pgx
  dsn: "postgres://user:pass@localhost:5432/db"
  migrationsDir: "./migrations/"
storage: "memory"
rmq:
  host: "0.0.0.0"
  port: "5672"
  login: "guest"
  password: "guest"
scheduler:
  period: 60s
  queue: "calendar_events"
`

func TestNewConfig(t *testing.T) {
	r := bytes.NewReader([]byte(yamlData))
	cfg, err := New(r)
	require.NoError(t, err)
	require.Equal(t, "debug", cfg.Log.Level)
	require.Equal(t, "0.0.0.0", cfg.HTTP.Host)
	require.Equal(t, "3000", cfg.HTTP.Port)
	require.Equal(t, 5*time.Second, cfg.HTTP.ReadTimeout)
	require.Equal(t, 10*time.Second, cfg.HTTP.WriteTimeout)
	require.Equal(t, 30*time.Second, cfg.HTTP.IdleTimeout)
	require.Equal(t, 2*time.Second, cfg.HTTP.ReadHeaderTimeout)
	require.Equal(t, "0.0.0.0", cfg.GRPC.Host)
	require.Equal(t, "50051", cfg.GRPC.Port)
	require.Equal(t, 5*time.Second, cfg.GRPC.ConnectTimeout)
	require.Equal(t, "pgx", cfg.DB.Type)
	require.Equal(t, "postgres://user:pass@localhost:5432/db", cfg.DB.DSN)
	require.Equal(t, "./migrations/", cfg.DB.MigrationsDir)
	require.Equal(t, true, cfg.DB.AutoMigrate)
	require.Equal(t, 3, cfg.App.RefillRate.Count)
	require.Equal(t, 30*time.Second, cfg.App.RefillRate.Time)
	require.Equal(t, true, cfg.App.GarbageCollector.Enabled)
	require.Equal(t, 600*time.Second, cfg.App.GarbageCollector.TTL)
	require.Equal(t, 60*time.Second, cfg.App.GarbageCollector.Interval)
}

func TestConfigContext(t *testing.T) {
	r := bytes.NewReader([]byte(yamlData))
	cfg, err := New(r)
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}
	ctx := context.Background()
	ctxWithCfg := cfg.WithContext(ctx)
	cfgFromCtx := GetFromContext(ctxWithCfg)
	require.Same(t, cfg, cfgFromCtx)
}

func TestNewConfig_InvalidYAML(t *testing.T) {
	r := bytes.NewReader([]byte("invalid"))
	_, err := New(r)
	require.Error(t, err)
}

func TestEnv(t *testing.T) {
	r := bytes.NewReader([]byte(yamlData))

	err := os.Setenv("HTTP_HOST", "1.1.1.1")
	require.NoError(t, err)
	err = os.Setenv("DB_MIGRATIONS_DIR", "test")
	require.NoError(t, err)
	err = os.Setenv("SCHEDULER_PERIOD", "0s")
	require.NoError(t, err)
	err = os.Setenv("SCHEDULER_RETENTION_PERIOD", "1h")
	require.NoError(t, err)

	cfg, err := New(r)
	require.NoError(t, err)
	require.Equal(t, "1.1.1.1", cfg.HTTP.Host)
	require.Equal(t, "test", cfg.DB.MigrationsDir)
	require.Equal(t, "3000", cfg.HTTP.Port)
}
