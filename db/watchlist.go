package db

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
)

// Watchlist gets the given server's watchlist
func (db *DB) Watchlist(guildID discord.GuildID) (list []uint64, err error) {
	err = db.Pool.QueryRow(context.Background(), "select watch_list from servers where id = $1", guildID).Scan(&list)
	return
}

// IsWatchlisted returns true if a user is on the watchlist
func (db *DB) IsWatchlisted(guildID discord.GuildID, userID discord.UserID) (b bool) {
	db.Pool.QueryRow(context.Background(), "select $1 = any(watch_list) from (select * from servers where id = $2) as server", userID, guildID).Scan(&b)
	return b
}

// WatchlistChannel ...
func (db *DB) WatchlistChannel(guildID discord.GuildID) (c discord.ChannelID, err error) {
	err = db.Pool.QueryRow(context.Background(), "select watch_list_channel from servers where id = $1", guildID).Scan(&c)
	return
}

// AddToWatchlist adds the given userID to the watchlist for guildID
func (db *DB) AddToWatchlist(guildID discord.GuildID, userID discord.UserID) (err error) {
	if db.IsWatchlisted(guildID, userID) {
		return ErrorAlreadyBlacklisted
	}
	_, err = db.Pool.Exec(context.Background(), "update servers set watch_list = array_append(watch_list, $1) where id = $2", userID, guildID)
	return err
}

// RemoveFromWatchlist removes the given userID from the watchlist for guildID
func (db *DB) RemoveFromWatchlist(guildID discord.GuildID, userID discord.UserID) (err error) {
	if !db.IsWatchlisted(guildID, userID) {
		return ErrorNotBlacklisted
	}
	_, err = db.Pool.Exec(context.Background(), "update servers set watch_list = array_remove(watch_list, $1) where id = $2", userID, guildID)
	if err != nil {
		return err
	}
	return err
}
