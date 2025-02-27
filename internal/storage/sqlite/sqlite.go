package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// https://www.sqlite.org/lang_create_table_agewgsd213ad
	// https://my.sh/sqlite

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
	    id INTEGER PRIMARY KEY,
	    alias TEXT NOT NULL UNIQUE,
	    url TEXT NOT NULL UNIQUE);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)
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

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(urlToSave, alias)

	if err != nil {
		// Парсим ошибку, проверяем является ли она sqlite3 ошибкой, и является ли ошибкой ErrConstraintUnique
		if sqliteErr, ok := err.(sqlite3.Error); ok && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

// GetURL getting a URL by Alias from DB
func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"
	var resURL string

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement:  %w", op, err)
	}
	row := stmt.QueryRow(alias)
	err = row.Scan(&resURL)

	if err != nil {
		// Проверяем, связана ли ошибка с отсутствием данных
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement:  %w", op, err)
	}
	res, errExec := stmt.Exec(alias)
	if errExec != nil {
		return fmt.Errorf("%s: execute statement: %w", op, errExec)
	}
	if rows, errRows := res.RowsAffected(); errRows == nil && rows == 0 {
		return storage.ErrAliasNotFound
	}

	return nil
}

func (s *Storage) AliasExists(alias string) (bool, error) {
	const op = "storage.sqlite.AliasExists"
	var resURL string

	stmt, err := s.db.Prepare("SELECT * FROM url WHERE alias = ?")
	if err != nil {
		return false, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	_, errExec := stmt.Exec(alias)
	if errExec != nil {
		return false, fmt.Errorf("%s: execute statement: %w", op, errExec)
	}
	row := stmt.QueryRow(alias)
	err = row.Scan(&resURL)

	if err != nil {
		// Проверяем, связана ли ошибка с отсутствием данных
		if errors.Is(err, sql.ErrNoRows) {
			return false, storage.ErrURLNotFound
		}
		return false, fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return true, nil
}
