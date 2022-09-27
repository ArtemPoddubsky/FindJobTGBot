package storage

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	addLastTimeout       = 2 * time.Second
	getLastTimeout       = 2 * time.Second
	checkOriginalTimeout = 2 * time.Second
	addRecordTimeout     = 2 * time.Second
	clearHistoryTimeout  = 2 * time.Second
)

var errUpdate = errors.New("couldn't update last request ")

// AddLast updates or adds last formed request in database.
func (db Postgres) AddLast(chatID int64, title, exp string) error {
	ctx, cancel := context.WithTimeout(context.Background(), addLastTimeout)
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

// GetLast gets last formed request for specific user.
func (db Postgres) GetLast(chatID int64) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), getLastTimeout)
	defer cancel()

	row := db.pool.QueryRow(ctx,
		"SELECT title, exp FROM last WHERE id = $1", chatID)

	var title, exp string

	if err := row.Scan(&title, &exp); err != nil {
		return "", "", fmt.Errorf("row.Scan: %w", err)
	}

	return title, exp, nil
}

// CheckOriginal checks if vacancy's url appeared in request history for specific user.
func (db Postgres) CheckOriginal(chatID int64, url string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), checkOriginalTimeout)
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

// AddRecord adds vacancy's url to user requests history.
func (db Postgres) AddRecord(chatID int64, url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), addRecordTimeout)
	defer cancel()

	if _, err := db.pool.Exec(ctx,
		"INSERT INTO vacancies(id, url) VALUES($1, $2)", chatID, url); err != nil {
		return fmt.Errorf("query exec: %w", err)
	}

	return nil
}

// ClearHistory clears user requests history.
func (db Postgres) ClearHistory(chatID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), clearHistoryTimeout)
	defer cancel()

	if _, err := db.pool.Exec(ctx,
		"DELETE FROM vacancies WHERE id = $1", chatID); err != nil {
		return fmt.Errorf("query exec: %w", err)
	}

	return nil
}
