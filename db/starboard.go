package db

import (
	"context"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
)

// StarboardMessage is a single starboard message
type StarboardMessage struct {
	MessageID          discord.MessageID
	ChannelID          discord.ChannelID
	ServerID           discord.GuildID
	StarboardMessageID discord.MessageID
	WebhookID          *discord.WebhookID
}

// StarboardMessage gets a starboard message by ID
func (db *DB) StarboardMessage(id discord.MessageID) (m *StarboardMessage, err error) {
	m = &StarboardMessage{}

	err = pgxscan.Get(context.Background(), db.Pool, m, "select message_id, channel_id, server_id, starboard_message_id, webhook_id from starboard_messages where message_id = $1 or starboard_message_id = $1", id)
	return
}

// SaveStarboardMessage ...
func (db *DB) SaveStarboardMessage(s StarboardMessage) (err error) {
	_, err = db.Pool.Exec(context.Background(), "insert into starboard_messages (message_id, channel_id, server_id, starboard_message_id, webhook_id) values ($1, $2, $3, $4, $5)", s.MessageID, s.ChannelID, s.ServerID, s.StarboardMessageID, s.WebhookID)
	return err
}

type Webhook struct {
	ID        discord.WebhookID
	ChannelID discord.ChannelID
	Token     string
}

func (db *DB) StarboardWebhook(id discord.WebhookID) (wh Webhook, err error) {
	err = pgxscan.Get(context.Background(), db.Pool, &wh, "select * from starboard_webhooks where id = $1", id)
	return wh, errors.Cause(err)
}

func (db *DB) StarboardChannelWebhook(id discord.ChannelID) (wh Webhook, err error) {
	err = pgxscan.Get(context.Background(), db.Pool, &wh, "select * from starboard_webhooks where channel_id = $1", id)
	return wh, errors.Cause(err)
}
