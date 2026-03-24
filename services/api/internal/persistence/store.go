package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	db *sql.DB
}

func New(dsn string) (*Store, error) {
	if dsn == "" {
		return nil, errors.New("missing postgres dsn")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	store := &Store{db: db}
	if err := store.Init(ctx); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *Store) Init(ctx context.Context) error {
	const query = `
CREATE TABLE IF NOT EXISTS app_records (
  kind TEXT NOT NULL,
  id TEXT NOT NULL,
  payload JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (kind, id)
)`
	_, err := s.db.ExecContext(ctx, query)
	return err
}

func (s *Store) Save(ctx context.Context, kind string, id string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	const query = `
INSERT INTO app_records (kind, id, payload)
VALUES ($1, $2, $3)
ON CONFLICT (kind, id)
DO UPDATE SET payload = EXCLUDED.payload, updated_at = NOW()`
	_, err = s.db.ExecContext(ctx, query, kind, id, body)
	return err
}

func (s *Store) Load(ctx context.Context, kind string, id string, dst any) (bool, error) {
	const query = `SELECT payload FROM app_records WHERE kind = $1 AND id = $2`
	var body []byte
	err := s.db.QueryRowContext(ctx, query, kind, id).Scan(&body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	if err := json.Unmarshal(body, dst); err != nil {
		return false, fmt.Errorf("unmarshal %s/%s: %w", kind, id, err)
	}
	return true, nil
}

func (s *Store) ListByKind(ctx context.Context, kind string, limit int, dst any) error {
	const query = `
SELECT payload
FROM app_records
WHERE kind = $1
ORDER BY updated_at DESC
LIMIT $2`

	rows, err := s.db.QueryContext(ctx, query, kind, limit)
	if err != nil {
		return err
	}
	defer rows.Close()

	bodies := make([]json.RawMessage, 0)
	for rows.Next() {
		var body []byte
		if err := rows.Scan(&body); err != nil {
			return err
		}
		bodies = append(bodies, append(json.RawMessage(nil), body...))
	}
	if err := rows.Err(); err != nil {
		return err
	}
	encoded, err := json.Marshal(bodies)
	if err != nil {
		return err
	}
	return json.Unmarshal(encoded, dst)
}
