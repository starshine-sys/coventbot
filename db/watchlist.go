package db

import (
	"context"

	"github.com/diamondburned/arikawa/v2/discord"
)

// Watchlist gets the given server's watchlist
func (db *DB) Watchlist(guildID discord.GuildID) (list []discord.UserID, err error) {
	err = db.Pool.QueryRow(context.Background(), "select watch_list from servers where id = $1", guildID).Scan(&list)
	return
}
