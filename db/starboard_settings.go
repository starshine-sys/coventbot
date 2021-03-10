package db

import (
	"context"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
)

// StarboardSettings is the starboard settings for a server
type StarboardSettings struct {
	StarboardChannel discord.ChannelID
	StarboardEmoji   string
	StarboardLimit   int
}

// Starboard gets the starboard settings for a server
func (db *DB) Starboard(id discord.GuildID) (s StarboardSettings, err error) {
	err = pgxscan.Get(context.Background(), db.Pool, &s, "select starboard_channel, starboard_emoji, starboard_limit from servers where id = $1", id)
	return
}

// SetStarboard sets the starboard settings for a server
func (db *DB) SetStarboard(id discord.GuildID, s StarboardSettings) (err error) {
	_, err = db.Pool.Exec(context.Background(), "update servers set starboard_channel = $1, starboard_emoji = $2, starboard_limit = $3 where id = $4", s.StarboardChannel, s.StarboardEmoji, s.StarboardLimit, id)
	return err
}
