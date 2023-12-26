package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"os"
	"time"
	"trip_service/internal/config"
)

func InitDB(ctx context.Context, config *config.Config) (*sqlx.DB, error) {
	time.Sleep(2 * time.Second)
	db, err := sqlx.Open("postgres", config.DB.DSN)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	fs := os.DirFS(config.DB.MigrationsDir)
	goose.SetBaseFS(fs)
	if err = goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err = goose.UpContext(ctx, db.DB, "."); err != nil {
		panic(err)
	}
	return db, nil
}
