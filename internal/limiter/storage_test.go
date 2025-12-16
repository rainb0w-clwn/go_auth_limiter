package limiter_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/limiter"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/storage/postgres"
	"github.com/stretchr/testify/require"
)

func newTestStorage(t *testing.T) (*limiter.Storage, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	pg := &postgres.Storage{
		DB:  sqlxDB,
		Ctx: context.Background(),
	}

	storage := &limiter.Storage{pg}
	return storage, mock
}

func TestStorage_GetLimits(t *testing.T) {
	storage, mock := newTestStorage(t)

	rows := sqlmock.NewRows([]string{"type", "value", "description"}).
		AddRow("login", 100, "login limit").
		AddRow("api", 200, nil)

	mock.ExpectPrepare("SELECT \\* FROM rate_limit").
		ExpectQuery().
		WillReturnRows(rows)

	result, err := storage.GetLimits()
	require.NoError(t, err)
	require.Len(t, *result, 2)

	require.Equal(t, limiter.Type("login"), (*result)[0].LimitType)
	require.Equal(t, 100, (*result)[0].Value)
	require.Equal(t, "login limit", (*result)[0].Description)

	require.Equal(t, limiter.Type("api"), (*result)[1].LimitType)
	require.Equal(t, 200, (*result)[1].Value)
	require.Empty(t, (*result)[1].Description)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_GetLimits_PrepareError(t *testing.T) {
	storage, mock := newTestStorage(t)

	mock.ExpectPrepare("SELECT \\* FROM rate_limit").
		WillReturnError(sql.ErrConnDone)

	_, err := storage.GetLimits()
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_GetLimits_SelectError(t *testing.T) {
	storage, mock := newTestStorage(t)

	mock.ExpectPrepare("SELECT \\* FROM rate_limit").
		ExpectQuery().
		WillReturnError(sql.ErrConnDone)

	_, err := storage.GetLimits()
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_GetLimitsByTypes(t *testing.T) {
	storage, mock := newTestStorage(t)

	rows := sqlmock.NewRows([]string{"type", "value", "description"}).
		AddRow("login", 100, "login limit").
		AddRow("api", 200, nil)

	types := []string{"login", "api"}
	query := "SELECT \\* FROM rate_limit WHERE type IN\\(\\$1, \\$2\\)"

	driverArgs := make([]driver.Value, len(types))
	for i, v := range types {
		driverArgs[i] = v
	}

	// Подготавливаем mock под итоговый SQL после Rebind
	mock.ExpectPrepare(query).
		ExpectQuery().
		WithArgs(driverArgs...).
		WillReturnRows(rows)

	result, err := storage.GetLimitsByTypes(types)
	require.NoError(t, err)
	require.Len(t, *result, 2)

	require.Equal(t, limiter.Type("login"), (*result)[0].LimitType)
	require.Equal(t, 100, (*result)[0].Value)
	require.Equal(t, "login limit", (*result)[0].Description)

	require.Equal(t, limiter.Type("api"), (*result)[1].LimitType)
	require.Equal(t, 200, (*result)[1].Value)
	require.Empty(t, (*result)[1].Description)

	require.NoError(t, mock.ExpectationsWereMet())
}
