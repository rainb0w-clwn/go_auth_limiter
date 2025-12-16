package postgres

import (
	"context"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // driver import
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/config"
)

type Storage struct {
	DB  *sqlx.DB
	Ctx context.Context
}

var ErrConnectFailed = errors.New("error connecting to DB")

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	cfg := config.GetFromContext(ctx)
	if cfg == nil {
		return config.ErrNoConfigInContext
	}

	db, err := sqlx.Open(cfg.DB.Type, cfg.DB.DSN)
	if err != nil {
		return fmt.Errorf(ErrConnectFailed.Error()+":%w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf(ErrConnectFailed.Error()+":%w", err)
	}

	s.DB = db
	s.Ctx = ctx
	if cfg.DB.AutoMigrate {
		return s.migrate(cfg.DB.MigrationsDir)
	}
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	if s.DB == nil {
		return nil
	}

	err := s.DB.Close()
	s.DB = nil
	s.Ctx = nil

	return err
}

func (s *Storage) migrate(migrationDir string) error {
	if s.DB == nil {
		return fmt.Errorf("database connection is not established")
	}

	if err := goose.SetDialect(s.DB.DriverName()); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(s.DB.DB, migrationDir); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
