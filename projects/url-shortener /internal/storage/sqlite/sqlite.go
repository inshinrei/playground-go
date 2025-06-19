package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"url-shortener_/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	const op = "sqlite.New"
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(data string, alias string) (int64, error) {
	const op = "sqlite.saveurl"
	stmt, err := s.db.Prepare(`INSERT INTO url(url, alias) VALUES(?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(data, alias)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) LoadURL(alias string) (string, error) {
	const op = "sqlite.loadurl"
	stmt, err := s.db.Prepare(`SELECT url FROM url WHERE alias = ?`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var res string
	err = stmt.QueryRow(alias).Scan(&res)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "sqlite.deleteurl"
	stmt, err := s.db.Prepare(`DELETE FROM url WHERE alias = ?`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(alias)
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrURLNotFound
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
