package rule_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/rule"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/storage/postgres"
	"github.com/stretchr/testify/require"
)

func newTestStorage(t *testing.T) (*rule.Storage, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	pg := &postgres.Storage{
		DB:  sqlxDB,
		Ctx: context.Background(),
	}

	storage := rule.NewStorage(pg)

	return storage, mock
}

func TestStorage_Create(t *testing.T) {
	storage, mock := newTestStorage(t)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(42)

	mock.ExpectPrepare("INSERT INTO ip_net_rule").
		ExpectQuery().
		WithArgs(
			sqlmock.AnyArg(), // ip
			sqlmock.AnyArg(), // type
		).
		WillReturnRows(rows)

	id, err := storage.Create(rule.Rule{
		IP:       "127.0.0.1",
		RuleType: rule.WhiteList,
	})

	require.NoError(t, err)
	require.Equal(t, 42, id)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_Create_PrepareError(t *testing.T) {
	storage, mock := newTestStorage(t)

	mock.ExpectPrepare("INSERT INTO ip_net_rule").
		WillReturnError(sql.ErrConnDone)

	_, err := storage.Create(rule.Rule{
		IP:       "127.0.0.1",
		RuleType: rule.WhiteList,
	})

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_Delete(t *testing.T) {
	storage, mock := newTestStorage(t)

	mock.ExpectExec("DELETE FROM ip_net_rule").
		WithArgs(sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := storage.Delete(10)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_GetForType(t *testing.T) {
	storage, mock := newTestStorage(t)

	rows := sqlmock.NewRows([]string{"id", "ip", "type"}).
		AddRow(1, "127.0.0.1", "white").
		AddRow(2, "10.0.0.0/24", "white")

	mock.ExpectPrepare("SELECT \\*").
		ExpectQuery().
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	result, err := storage.GetForType(rule.WhiteList)

	require.NoError(t, err)
	require.Len(t, *result, 2)

	require.Equal(t, 1, (*result)[0].ID)
	require.Equal(t, "127.0.0.1", (*result)[0].IP)
	require.Equal(t, rule.WhiteList, (*result)[0].RuleType)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_Find(t *testing.T) {
	storage, mock := newTestStorage(t)

	rows := sqlmock.NewRows([]string{"id", "ip", "type"}).
		AddRow(5, "127.0.0.1", "black")

	mock.ExpectPrepare("SELECT \\*").
		ExpectQuery().
		WithArgs(
			sqlmock.AnyArg(), // ip
			sqlmock.AnyArg(), // type
		).
		WillReturnRows(rows)

	result, err := storage.Find("127.0.0.1", rule.BlackList)

	require.NoError(t, err)
	require.Len(t, *result, 1)

	require.Equal(t, 5, (*result)[0].ID)
	require.Equal(t, rule.BlackList, (*result)[0].RuleType)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_Find_Empty(t *testing.T) {
	storage, mock := newTestStorage(t)

	rows := sqlmock.NewRows([]string{"id", "ip", "type"})

	mock.ExpectPrepare("SELECT \\*").
		ExpectQuery().
		WillReturnRows(rows)

	result, err := storage.Find("127.0.0.1", rule.WhiteList)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, *result, 0)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStorage_GetForType_SelectError(t *testing.T) {
	storage, mock := newTestStorage(t)

	mock.ExpectPrepare("SELECT \\*").
		ExpectQuery().
		WillReturnError(sql.ErrConnDone)

	_, err := storage.GetForType(rule.WhiteList)

	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
