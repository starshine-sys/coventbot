package pklog

import (
	"context"
	"encoding/hex"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
)

// Message is a single proxied message
type Message struct {
	MsgID     discord.MessageID
	UserID    discord.UserID
	ChannelID discord.ChannelID
	ServerID  discord.GuildID

	Username string
	Member   string
	System   string

	Content string
}

// Insert inserts a message
func (bot *Bot) Insert(m Message) (err error) {
	if m.Content == "" {
		m.Content = "None"
	}
	out, err := Encrypt([]byte(m.Content), bot.AESKey)
	if err != nil {
		return err
	}
	m.Content = hex.EncodeToString(out)

	out, err = Encrypt([]byte(m.Username), bot.AESKey)
	if err != nil {
		return err
	}
	m.Username = hex.EncodeToString(out)

	_, err = bot.DB.Pool.Exec(context.Background(), `insert into pk_messages
(msg_id, user_id, channel_id, server_id, username, member, system, content) values
($1, $2, $3, $4, $5, $6, $7, $8)`, m.MsgID, m.UserID, m.ChannelID, m.ServerID, m.Username, m.Member, m.System, m.Content)
	return err
}

// Get gets a single message
func (bot *Bot) Get(id discord.MessageID) (m *Message, err error) {
	m = &Message{}

	var exists bool
	err = bot.DB.Pool.QueryRow(context.Background(), "select exists(select * from pk_messages where msg_id = $1)", id).Scan(&exists)
	if !exists {
		return nil, ErrNotExists
	}

	err = pgxscan.Get(context.Background(), bot.DB.Pool, m, "select * from pk_messages where msg_id = $1", id)

	b, err := hex.DecodeString(m.Content)
	if err != nil {
		return nil, err
	}

	out, err := Decrypt(b, bot.AESKey)
	if err != nil {
		return nil, err
	}

	m.Content = string(out)

	b, err = hex.DecodeString(m.Username)
	if err != nil {
		return nil, err
	}

	out, err = Decrypt(b, bot.AESKey)
	if err != nil {
		return nil, err
	}

	m.Username = string(out)
	return
}

// Delete deletes a message from the database
func (bot *Bot) Delete(id discord.MessageID) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "delete from pk_messages where msg_id = $1", id)
	return
}
