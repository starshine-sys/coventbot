package reactroles

import (
	"context"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
)

// Message ...
type Message struct {
	ServerID  discord.GuildID
	ChannelID discord.ChannelID
	MessageID discord.MessageID
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
