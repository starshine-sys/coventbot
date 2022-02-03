package db

import (
	"context"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jackc/pgx/v4"
)

func (db *DB) GuildStringGet(guildID discord.GuildID, key string) (val string, err error) {
	err = db.Pool.QueryRow(context.Background(), "select coalesce(config->$1, '') from guild_config where guild_id = $2", key, guildID).Scan(&val)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	return val, err
}

func (db *DB) GuildBoolGet(guildID discord.GuildID, key string) (val bool, err error) {
	err = db.Pool.QueryRow(context.Background(), "select coalesce(config->$1, 'false')::boolean from guild_config where guild_id = $2", key, guildID).Scan(&val)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	return val, err
}

func (db *DB) GuildIntGet(guildID discord.GuildID, key string) (val int64, err error) {
	err = db.Pool.QueryRow(context.Background(), "select coalesce(config->$1, '0')::bigint from guild_config where guild_id = $2", key, guildID).Scan(&val)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	return val, err
}

func (db *DB) GuildStringSet(guildID discord.GuildID, key, val string) error {
	sql := `insert into guild_config (guild_id, config)
	values ($1, hstore($2, $3))
	on conflict (guild_id) do update
	set config = guild_config.config || hstore($2, $3)`

	_, err := db.Pool.Exec(context.Background(), sql, guildID, key, val)
	return err
}

func (db *DB) GuildBoolSet(guildID discord.GuildID, key string, val bool) error {
	sql := `insert into guild_config (guild_id, config)
	values ($1, hstore($2, $3))
	on conflict (guild_id) do update
	set config = guild_config.config || hstore($2, $3)`

	_, err := db.Pool.Exec(context.Background(), sql, guildID, key, strconv.FormatBool(val))
	return err
}

func (db *DB) GuildIntSet(guildID discord.GuildID, key string, val int64) error {
	sql := `insert into guild_config (guild_id, config)
	values ($1, hstore($2, $3))
	on conflict (guild_id) do update
	set config = guild_config.config || hstore($2, $3)`

	_, err := db.Pool.Exec(context.Background(), sql, guildID, key, strconv.FormatInt(val, 10))
	return err
}

func (db *DB) GuildKeyDelete(guildID discord.GuildID, key string) error {
	sql := `insert into guild_config (guild_id, config)
	values ($1, ''::hstore)
	on conflict (guild_id) do update
	set config = delete(guild_config.config, $2)`

	_, err := db.Pool.Exec(context.Background(), sql, guildID, key)
	return err
}
