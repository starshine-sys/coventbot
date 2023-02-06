// SPDX-License-Identifier: AGPL-3.0-only
package db

import (
	"context"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jackc/pgx/v4"
)

func (db *DB) UserStringGet(userID discord.UserID, key string) (val string, err error) {
	err = db.Pool.QueryRow(context.Background(), "select coalesce(config->$1, '') from user_config where user_id = $2", key, userID).Scan(&val)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	return val, err
}

func (db *DB) UserBoolGet(userID discord.UserID, key string) (val bool, err error) {
	err = db.Pool.QueryRow(context.Background(), "select coalesce(config->$1, 'false')::boolean from user_config where user_id = $2", key, userID).Scan(&val)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	return val, err
}

func (db *DB) UserIntGet(userID discord.UserID, key string) (val int64, err error) {
	err = db.Pool.QueryRow(context.Background(), "select coalesce(config->$1, '0')::bigint from user_config where user_id = $2", key, userID).Scan(&val)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	return val, err
}

func (db *DB) UserStringSet(userID discord.UserID, key, val string) error {
	sql := `insert into user_config (user_id, config)
	values ($1, hstore($2, $3))
	on conflict (user_id) do update
	set config = user_config.config || hstore($2, $3)`

	_, err := db.Pool.Exec(context.Background(), sql, userID, key, val)
	return err
}

func (db *DB) UserBoolSet(userID discord.UserID, key string, val bool) error {
	sql := `insert into user_config (user_id, config)
	values ($1, hstore($2, $3))
	on conflict (user_id) do update
	set config = user_config.config || hstore($2, $3)`

	_, err := db.Pool.Exec(context.Background(), sql, userID, key, strconv.FormatBool(val))
	return err
}

func (db *DB) UserIntSet(userID discord.UserID, key string, val int64) error {
	sql := `insert into user_config (user_id, config)
	values ($1, hstore($2, $3))
	on conflict (user_id) do update
	set config = user_config.config || hstore($2, $3)`

	_, err := db.Pool.Exec(context.Background(), sql, userID, key, strconv.FormatInt(val, 10))
	return err
}

func (db *DB) UserKeyDelete(userID discord.UserID, key string) error {
	sql := `insert into user_config (user_id, config)
	values ($1, ''::hstore)
	on conflict (user_id) do update
	set config = delete(user_config.config, $2)`

	_, err := db.Pool.Exec(context.Background(), sql, userID, key)
	return err
}
