package db

import (
	"context"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
)

// StarboardMessage is a single starboard message
type StarboardMessage struct {
	MessageID          discord.MessageID
	ChannelID          discord.ChannelID
	ServerID           discord.GuildID
	StarboardMessageID discord.MessageID
}

// StarboardMessage gets a starboard message by ID
func (db *DB) StarboardMessage(id discord.MessageID) (m *StarboardMessage, err error) {
	m = &StarboardMessage{}

	err = pgxscan.Get(context.Background(), db.Pool, m, "select message_id, channel_id, server_id, starboard_message_id from starboard_messages where message_id = $1 or starboard_message_id = $1", id)
	return
}

// SaveStarboardMessage ...
func (db *DB) SaveStarboardMessage(s StarboardMessage) (err error) {
	_, err = db.Pool.Exec(context.Background(), "insert into starboard_messages (message_id, channel_id, server_id, starboard_message_id) values ($1, $2, $3, $4)", s.MessageID, s.ChannelID, s.ServerID, s.StarboardMessageID)
	return err
}
