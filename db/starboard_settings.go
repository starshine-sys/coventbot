// SPDX-License-Identifier: AGPL-3.0-only
package db

import (
	"context"
	"errors"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
)

// StarboardSettings is the starboard settings for a server
type StarboardSettings struct {
	StarboardChannel   discord.ChannelID
	StarboardEmoji     string
	StarboardLimit     int
	StarboardUsername  string
	StarboardAvatarURL string
}

// Errors for setting the blacklist
var (
	ErrorAlreadyBlacklisted = errors.New("channel is already blacklisted")
	ErrorNotBlacklisted     = errors.New("channel is not blacklisted")
)

// Starboard gets the starboard settings for a server
func (db *DB) Starboard(id discord.GuildID) (s StarboardSettings, err error) {
	err = pgxscan.Get(context.Background(), db.Pool, &s, "select starboard_channel, starboard_emoji, starboard_limit, starboard_username, starboard_avatar_url from servers where id = $1", id)
	return
}

// SetStarboard sets the starboard settings for a server
func (db *DB) SetStarboard(id discord.GuildID, s StarboardSettings) (err error) {
	_, err = db.Pool.Exec(context.Background(), "update servers set starboard_channel = $1, starboard_emoji = $2, starboard_limit = $3, starboard_username = $4, starboard_avatar_url = $5 where id = $6", s.StarboardChannel, s.StarboardEmoji, s.StarboardLimit, s.StarboardUsername, s.StarboardAvatarURL, id)
	return err
}

// IsBlacklisted returns true if a channel is blacklisted
func (db *DB) IsBlacklisted(guildID discord.GuildID, channelID discord.ChannelID) (b bool) {
	_ = db.Pool.QueryRow(context.Background(), "select $1 = any(starboard_blacklist) from (select * from servers where id = $2) as server", channelID, guildID).Scan(&b)
	return b
}

// StarboardBlacklist gets the current starboard blacklist
func (db *DB) StarboardBlacklist(id discord.GuildID) (bl []uint64, err error) {
	err = db.Pool.QueryRow(context.Background(), "select starboard_blacklist from servers where id = $1", id).Scan(&bl)
	return
}

// AddToBlacklist adds the given channelID to the blacklist for guildID
func (db *DB) AddToBlacklist(guildID discord.GuildID, channelID discord.ChannelID) (err error) {
	if db.IsBlacklisted(guildID, channelID) {
		return ErrorAlreadyBlacklisted
	}
	_, err = db.Pool.Exec(context.Background(), "update servers set starboard_blacklist = array_append(starboard_blacklist, $1) where id = $2", channelID, guildID)
	return err
}

// RemoveFromBlacklist removes the given channelID from the blacklist for guildID
func (db *DB) RemoveFromBlacklist(guildID discord.GuildID, channelID discord.ChannelID) (err error) {
	if !db.IsBlacklisted(guildID, channelID) {
		return ErrorNotBlacklisted
	}
	_, err = db.Pool.Exec(context.Background(), "update servers set starboard_blacklist = array_remove(starboard_blacklist, $1) where id = $2", channelID, guildID)
	if err != nil {
		return err
	}
	return err
}
