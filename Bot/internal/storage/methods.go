package storage

import (
	"context"
	"errors"
)

func (db Postgres) AddLast(chatID int64, title, exp string) error {
	res, err := db.pool.Exec(context.Background(), "INSERT INTO last(id, title, exp) VALUES($1, $2, $3) ON CONFLICT (id) DO UPDATE SET title = $2, exp = $3", chatID, title, exp)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("Couldn't update last request ")
	}
	return nil
}

func (db Postgres) GetLast(chatID int64) (string, string, error) {
	row := db.pool.QueryRow(context.Background(), "SELECT title, exp FROM last WHERE id = $1", chatID)
	var title, exp string
	if err := row.Scan(&title, &exp); err != nil {
		return "", "", err
	}
	return title, exp, nil
}

func (db Postgres) CheckOriginal(chatID int64, url string) (bool, error) {
	res, err := db.pool.Exec(context.Background(), "SELECT url FROM vacancies WHERE id = $1 AND url = $2", chatID, url)
	if err != nil {
		return false, err
	}
	if res.RowsAffected() == 0 {
		return true, nil
	}
	return false, nil
}

func (db Postgres) AddRecord(chatID int64, url string) error {
	if _, err := db.pool.Exec(context.Background(), "INSERT INTO vacancies(id, url) VALUES($1, $2)", chatID, url); err != nil {
		return err
	}
	return nil
}

func (db Postgres) ClearHistory(chatID int64) error {
	if _, err := db.pool.Exec(context.Background(), "DELETE FROM vacancies WHERE id = $1", chatID); err != nil {
		return err
	}
	return nil
}
