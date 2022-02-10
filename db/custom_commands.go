package db

import (
	"context"
	"strings"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type CustomCommand struct {
	ID      int64
	GuildID discord.GuildID
	Name    string
	Source  string
}

const ErrCommandNotFound = errors.Sentinel("no cc with that name found")

func (db *DB) CustomCommand(guildID discord.GuildID, name string) (c CustomCommand, err error) {
	err = pgxscan.Get(context.Background(), db.Pool, &c, "select * from custom_commands where guild_id = $1 and name = $2", guildID, strings.ToLower(name))
	if errors.Cause(err) == pgx.ErrNoRows {
		return c, ErrCommandNotFound
	}
	return c, errors.Cause(err)
}

func (db *DB) CustomCommandID(guildID discord.GuildID, id int64) (c CustomCommand, err error) {
	err = pgxscan.Get(context.Background(), db.Pool, &c, "select * from custom_commands where guild_id = $1 and id = $2", guildID, id)
	if errors.Cause(err) == pgx.ErrNoRows {
		return c, ErrCommandNotFound
	}
	return c, errors.Cause(err)
}

func (db *DB) SetCustomCommand(guildID discord.GuildID, name, source string) (c CustomCommand, err error) {
	err = pgxscan.Get(context.Background(), db.Pool, &c, `insert into custom_commands
	(guild_id, name, source) values
	($1, $2, $3) on conflict (guild_id, lower(name))
	do update set source = $3
	returning *`, guildID, name, source)
	return c, errors.Cause(err)
}

func (db *DB) AllCustomCommands(guildID discord.GuildID) (c []CustomCommand, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &c, "select * from custom_commands where guild_id = $1 order by id", guildID)
	return c, errors.Cause(err)
}
