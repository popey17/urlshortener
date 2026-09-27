package store

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct{ db *sql.DB }

type URL struct {
	ID         int64
	Code       string
	Original   string
	ClickCount int
	CreatedAt  string
}

func OpenDB(databaseURL string) (*Store, error) {
	dbStore := &Store{}

	db, err := sql.Open("pgx", databaseURL)
	if err := db.Ping(); err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = migrateFile(db)
	if err != nil {
		return nil, err
	}

	dbStore.db = db

	return dbStore, nil

}

func migrateFile(db *sql.DB) error {
	// migrationFile, err := os.ReadFile("migrations/001_init.sql")

	// if err != nil {
	// 	return err
	// }

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS urls (
    id          BIGSERIAL PRIMARY KEY,
    code        TEXT NOT NULL,
    original    TEXT NOT NULL,
    click_count INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (code)
	);`); err != nil {
		return err
	}

	return nil
}

func (s *Store) Create(code, original string) (*URL, error) {
	query := `INSERT INTO urls (code, original)
	VALUES ($1, $2)
	RETURNING id, code, original, click_count, created_at`

	u := &URL{}
	err := s.db.QueryRow(query, code, original).Scan(
		&u.ID, &u.Code, &u.Original, &u.ClickCount, &u.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return u, nil

}

func (s *Store) GetByCode(code string) (*URL, error) {
	u := &URL{}
	query := "SELECT id, code, original, click_count, created_at FROM urls WHERE code = $1"
	err := s.db.QueryRow(query, code).Scan(&u.ID, &u.Code, &u.Original, &u.ClickCount, &u.CreatedAt)

	if err != nil {
		return nil, sql.ErrNoRows
	}
	return u, nil
}

func (s *Store) IncrementClicks(code string) error {
	increaseCountQuery := "UPDATE urls SET click_count = click_count + 1 WHERE code = $1"
	result, err := s.db.Exec(increaseCountQuery, code)
	if err != nil {
		return fmt.Errorf("Error increasing the click count")
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	} else if n == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
