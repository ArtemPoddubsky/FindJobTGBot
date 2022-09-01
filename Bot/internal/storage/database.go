package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"main/internal/config"
	"time"
)

const Attempts = 8
const Delay = 5 * time.Second

type Postgres struct {
	pool *pgxpool.Pool
}

func NewDB(cfg config.Config) *Postgres {
	psqlconn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.Db.Username, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.Database)
	return &Postgres{Connect(psqlconn)}
}

func Connect(psqlconn string) *pgxpool.Pool {
	var pool *pgxpool.Pool
	var err error
	err = TryToConnect(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), Delay)
		defer cancel()
		pool, err = pgxpool.Connect(ctx, psqlconn)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	return pool
}

func TryToConnect(fn func() error) (err error) {
	for count, j := 1, Attempts; j > 0; {
		log.Warnf("Trying to connect to database %d time", count)
		if err = fn(); err != nil {
			time.Sleep(Delay)
			count++
			j--
			continue
		}
		return nil
	}
	return
}
