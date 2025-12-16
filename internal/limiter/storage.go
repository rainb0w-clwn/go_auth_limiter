package limiter

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rainb0w-clwn/go_auth_limiter/internal/storage/postgres"
)

type sqlEntity struct {
	LimitType   string         `db:"type"`
	Value       int            `db:"value"`
	Description sql.NullString `db:"description"`
}

type Storage struct {
	*postgres.Storage
}

func NewStorage(storage *postgres.Storage) IStorage {
	return &Storage{storage}
}

func (s *Storage) GetLimits() (*Limits, error) {
	query := `
		SELECT *
		FROM rate_limit
	`

	stmt, err := s.DB.PreparexContext(s.Ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var rows []sqlEntity

	err = stmt.SelectContext(
		s.Ctx,
		&rows,
	)
	if err != nil {
		return nil, err
	}

	result := make(Limits, 0, len(rows))
	for _, r := range rows {
		result = append(result, *s.sqlEntityToEntity(&r))
	}

	return &result, nil
}

func (s *Storage) GetLimitsByTypes(types []string) (*Limits, error) {
	arg := map[string]any{
		"types": types,
	}
	query := `SELECT * FROM rate_limit WHERE type IN(:types)`
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return nil, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query) // query теперь с $1,$2
	stmt, err := s.DB.PreparexContext(s.Ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var rows []sqlEntity

	err = stmt.SelectContext(
		s.Ctx,
		&rows,
		args...,
	)
	if err != nil {
		return nil, err
	}

	result := make(Limits, 0, len(rows))
	for _, r := range rows {
		result = append(result, *s.sqlEntityToEntity(&r))
	}

	return &result, nil
}

func (s *Storage) sqlEntityToEntity(se *sqlEntity) *Limit {
	e := &Limit{
		LimitType: Type(se.LimitType),
		Value:     se.Value,
	}

	if se.Description.Valid {
		e.Description = se.Description.String
	}

	return e
}
