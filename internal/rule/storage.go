package rule

import (
	"github.com/rainb0w-clwn/go_auth_limiter/internal/storage/postgres"
)

type sqlEntity struct {
	ID       int    `db:"id"`
	IP       string `db:"ip"`
	RuleType string `db:"type"`
}

type Storage struct {
	*postgres.Storage
}

func NewStorage(storage *postgres.Storage) *Storage {
	return &Storage{storage}
}

func (s *Storage) Create(rule Rule) (int, error) {
	query := `
		INSERT INTO ip_net_rule(ip, type) VALUES (:ip, :type)
		RETURNING id
	`

	params := map[string]any{
		"ip":   rule.IP,
		"type": rule.RuleType,
	}

	var id int
	stmt, err := s.DB.PrepareNamedContext(s.Ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	if err = stmt.GetContext(s.Ctx, &id, params); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) Delete(id int) error {
	query := `
		DELETE FROM ip_net_rule
		WHERE id = :id
	`

	_, err := s.DB.NamedExecContext(
		s.Ctx,
		query,
		map[string]any{"id": id},
	)
	return err
}

func (s *Storage) GetForType(ruleType Type) (*Rules, error) {
	query := `
		SELECT *
		FROM ip_net_rule
		WHERE type = :type
	`

	stmt, err := s.DB.PrepareNamedContext(s.Ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var rows []sqlEntity

	err = stmt.SelectContext(
		s.Ctx,
		&rows,
		map[string]any{
			"type": ruleType,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make(Rules, 0, len(rows))
	for _, r := range rows {
		result = append(result, *s.sqlEntityToEntity(&r))
	}

	return &result, nil
}

func (s *Storage) Find(ip string, ruleType Type) (*Rules, error) {
	query := `
		SELECT *
		FROM ip_net_rule
		WHERE ip = :ip 
			AND type = :type
	`

	stmt, err := s.DB.PrepareNamedContext(s.Ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var rows []sqlEntity

	err = stmt.SelectContext(
		s.Ctx,
		&rows,
		map[string]any{
			"ip":   ip,
			"type": ruleType,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make(Rules, 0, len(rows))
	for _, r := range rows {
		result = append(result, *s.sqlEntityToEntity(&r))
	}

	return &result, nil
}

func (s *Storage) sqlEntityToEntity(se *sqlEntity) *Rule {
	e := &Rule{
		ID:       se.ID,
		IP:       se.IP,
		RuleType: Type(se.RuleType),
	}

	return e
}
