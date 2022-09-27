package storage

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	AddLastTimeout       = 2 * time.Second
	GetLastTimeout       = 2 * time.Second
	CheckOriginalTimeout = 2 * time.Second
	AddRecordTimeout     = 2 * time.Second
	ClearHistoryTimeout  = 2 * time.Second
)

var errUpdate = errors.New("couldn't update last request ")

func (db Postgres) AddLast(chatID int64, title, exp string) error {
	ctx, cancel := context.WithTimeout(context.Background(), AddLastTimeout)
	defer cancel()

	res, err := db.pool.Exec(ctx,
		"INSERT INTO last(id, title, exp) VALUES($1, $2, $3)"+
			"ON CONFLICT (id) DO UPDATE SET title = $2, exp = $3", chatID, title, exp)

	if err != nil {
		return fmt.Errorf("query exec: %w", err)
	}

	if res.RowsAffected() == 0 {
		return errUpdate
	}

	return nil
}

func (db Postgres) GetLast(chatID int64) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GetLastTimeout)
	defer cancel()

	row := db.pool.QueryRow(ctx,
		"SELECT title, exp FROM last WHERE id = $1", chatID)

	var title, exp string

	if err := row.Scan(&title, &exp); err != nil {
		return "", "", fmt.Errorf("row.Scan: %w", err)
	}

	return title, exp, nil
}

func (db Postgres) CheckOriginal(chatID int64, url string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CheckOriginalTimeout)
	defer cancel()

	res, err := db.pool.Exec(ctx,
		"SELECT url FROM vacancies WHERE id = $1 AND url = $2", chatID, url)

	if err != nil {
		return false, fmt.Errorf("query exec: %w", err)
	} else if res.RowsAffected() == 0 {
		return true, nil
	}

	return false, nil
}

func (db Postgres) AddRecord(chatID int64, url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), AddRecordTimeout)
	defer cancel()

	if _, err := db.pool.Exec(ctx,
		"INSERT INTO vacancies(id, url) VALUES($1, $2)", chatID, url); err != nil {
		return fmt.Errorf("query exec: %w", err)
	}

	return nil
}

func (db Postgres) ClearHistory(chatID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), ClearHistoryTimeout)
	defer cancel()

	if _, err := db.pool.Exec(ctx,
		"DELETE FROM vacancies WHERE id = $1", chatID); err != nil {
		return fmt.Errorf("query exec: %w", err)
	}

	return nil
}
