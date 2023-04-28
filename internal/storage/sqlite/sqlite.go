package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zhayt/clean-arch-tmp-forum/config"
	"time"
)

func Dial(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("couldn't get db pool connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("couldn't connect to database: %w", err)
	}

	return db, nil
}
