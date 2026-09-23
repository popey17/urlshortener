package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Store struct{ db *sql.DB }

type URL struct {
	ID         int64
	Code       string
	Original   string
	ClickCount int
	CreatedAt  string
}

func OpenDB(path string) (*Store, error) {
	dbStore := &Store{}

	err := os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
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
	migrationFile, err := os.ReadFile("migrations/001_init.sql")

	if err != nil {
		return err
	}

	if _, err := db.Exec(string(migrationFile)); err != nil {
		return err
	}

	return nil
}

func (s *Store) Create(code, original string) (*URL, error) {
	query := "INSERT INTO urls (code, original) VALUES (?, ?)"
	result, err := s.db.Exec(query, code, original)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	u := &URL{}
	getQuery := "SELECT id, code, original, click_count, created_at FROM urls WHERE id = ?"
	err = s.db.QueryRow(getQuery, id).Scan(&u.ID, &u.Code, &u.Original, &u.ClickCount, &u.CreatedAt)

	if err != nil {
		return nil, err
	}
	return u, nil

}

func (s *Store) GetByCode(code string) (*URL, error) {
	u := &URL{}
	query := "SELECT id, code, original, click_count, created_at FROM urls WHERE code = ?"
	err := s.db.QueryRow(query, code).Scan(&u.ID, &u.Code, &u.Original, &u.ClickCount, &u.CreatedAt)

	if err != nil {
		return nil, sql.ErrNoRows
	}
	return u, nil
}

func (s *Store) IncrementClicks(code string) error {
	increaseCountQuery := "UPDATE urls SET click_count = click_count + 1 WHERE code = ?"
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
