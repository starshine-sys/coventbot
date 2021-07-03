package db

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
)

// Prefixes returns the server's prefixes
func (db *DB) Prefixes(id discord.GuildID) (prefixes []string, err error) {
	err = db.Pool.QueryRow(context.Background(), "select prefixes from servers where id = $1", id).
		Scan(&prefixes)
	return
}

// SetPrefixes sets the prefixes for a server
func (db *DB) SetPrefixes(id discord.GuildID, prefixes []string) (err error) {
	_, err = db.Pool.Exec(context.Background(), "update servers set prefixes = $1 where id = $2", prefixes, id)
	return
}
