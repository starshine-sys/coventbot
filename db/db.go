package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/starshine-sys/coventbot/types"
	"go.uber.org/zap"
)

// Misc errors
var (
	ErrTimedOut = errors.New("db: timed out")
)

// DB ...
type DB struct {
	Pool  *pgxpool.Pool
	Sugar *zap.SugaredLogger

	Config *types.BotConfig
}

// New ...
func New(url string, sugar *zap.SugaredLogger, c *types.BotConfig) (db *DB, err error) {
	pool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	db = &DB{
		Pool:   pool,
		Sugar:  sugar,
		Config: c,
	}

	return
}
