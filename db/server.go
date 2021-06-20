package db

import (
	"context"

	"github.com/diamondburned/arikawa/v2/discord"
)

// CreateServerIfNotExists returns true if the server exists
func (db *DB) CreateServerIfNotExists(guildID discord.GuildID) (exists bool, err error) {
	err = db.Pool.QueryRow(context.Background(), "select exists (select from servers where id = $1)", guildID).Scan(&exists)
	if err != nil {
		return exists, err
	}
	if !exists {
		_, err = db.Pool.Exec(context.Background(), "insert into servers (id, prefixes) values ($1, $2)", guildID, db.Config.Prefixes)
		return exists, err
	}
	return exists, nil
}
