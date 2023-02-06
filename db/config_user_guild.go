// SPDX-License-Identifier: AGPL-3.0-only
package db

import (
	"context"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jackc/pgx/v4"
)

func (db *DB) GuildUserStringGet(guildID discord.GuildID, userID discord.UserID, key string) (val string, err error) {
	err = db.Pool.QueryRow(context.Background(), "select coalesce(config->$1, '') from user_guild_config where user_id = $2 and guild_id = $3", key, userID, guildID).Scan(&val)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	return val, err
}

func (db *DB) GuildUserBoolGet(guildID discord.GuildID, userID discord.UserID, key string) (val bool, err error) {
	err = db.Pool.QueryRow(context.Background(), "select coalesce(config->$1, 'false')::boolean from user_guild_config where user_id = $2 and guild_id = $3", key, userID, guildID).Scan(&val)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	return val, err
}

func (db *DB) GuildUserIntGet(guildID discord.GuildID, userID discord.UserID, key string) (val int64, err error) {
	err = db.Pool.QueryRow(context.Background(), "select coalesce(config->$1, '0')::bigint from user_guild_config where user_id = $2 and guild_id = $3", key, userID, guildID).Scan(&val)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	return val, err
}

func (db *DB) GuildUserStringSet(guildID discord.GuildID, userID discord.UserID, key, val string) error {
	sql := `insert into user_guild_config (user_id, guild_id, config)
	values ($1, $2, hstore($3, $4))
	on conflict (user_id, guild_id) do update
	set config = user_guild_config.config || hstore($3, $4)`

	_, err := db.Pool.Exec(context.Background(), sql, userID, guildID, key, val)
	return err
}

func (db *DB) GuildUserBoolSet(guildID discord.GuildID, userID discord.UserID, key string, val bool) error {
	sql := `insert into user_guild_config (user_id, guild_id, config)
	values ($1, $2, hstore($3, $4))
	on conflict (user_id, guild_id) do update
	set config = user_guild_config.config || hstore($3, $4)`

	_, err := db.Pool.Exec(context.Background(), sql, userID, guildID, key, strconv.FormatBool(val))
	return err
}

func (db *DB) GuildUserIntSet(guildID discord.GuildID, userID discord.UserID, key string, val int64) error {
	sql := `insert into user_guild_config (user_id, guild_id, config)
	values ($1, $2, hstore($3, $4))
	on conflict (user_id, guild_id) do update
	set config = user_guild_config.config || hstore($3, $4)`

	_, err := db.Pool.Exec(context.Background(), sql, userID, guildID, key, strconv.FormatInt(val, 10))
	return err
}

func (db *DB) GuildUserKeyDelete(guildID discord.GuildID, userID discord.UserID, key string) error {
	sql := `insert into user_guild_config (user_id, guild_id, config)
	values ($1, $2, ''::hstore)
	on conflict (user_id, guild_id) do update
	set config = delete(user_guild_config.config, $3)`

	_, err := db.Pool.Exec(context.Background(), sql, userID, guildID, key)
	return err
}
