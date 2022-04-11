package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/starshine-sys/tribble/common"
	"go.uber.org/zap"

	// pgx driver for migrations
	_ "github.com/jackc/pgx/v4/stdlib"
)

// Misc errors
var (
	ErrTimedOut = errors.New("db: timed out")
)

// DB ...
type DB struct {
	Pool  *pgxpool.Pool
	Sugar *zap.SugaredLogger

	Config *common.BotConfig
}

// New ...
func New(url string, sugar *zap.SugaredLogger, c *common.BotConfig) (db *DB, err error) {
	err = runMigrations(url, sugar)
	if err != nil {
		return nil, err
	}

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

//go:embed migrations
var fs embed.FS

func runMigrations(url string, sugar *zap.SugaredLogger) (err error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return err
	}
	// defer this just in case
	// we call db.Close just before returning at the end to capture the error
	// if something else went wrong it's fine to discard the error here
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: fs,
		Root:       "migrations",
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	if n != 0 {
		sugar.Infof("Performed %v migrations!", n)
	}

	// close the db here because we don't use the stdlib database package normally
	// pgx can get a pgx.Conn from a stdlib database, but not a pgxpool.Pool
	err = db.Close()
	return err
}
