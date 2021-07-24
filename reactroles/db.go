package reactroles

import (
	"context"
	"database/sql"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
)

// Message ...
type Message struct {
	ServerID  discord.GuildID
	ChannelID discord.ChannelID
	MessageID discord.MessageID

	Title       sql.NullString
	Description sql.NullString
	Mention     sql.NullBool
}

// Entry ...
type Entry struct {
	MessageID discord.MessageID
	Emote     string
	RoleID    discord.RoleID
}

func (bot *Bot) getEntries(mID discord.MessageID) (entries []Entry, err error) {
	err = pgxscan.Select(context.Background(), bot.DB.Pool, &entries, "select * from react_role_entries where message_id = $1", mID)
	return
}

func (bot *Bot) message(id discord.MessageID) (m Message, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &m, "select * from react_roles where message_id = $1", id)
	return
}
