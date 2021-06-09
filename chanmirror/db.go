package chanmirror

import (
	"context"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
)

// Mirror ...
type Mirror struct {
	ServerID discord.GuildID

	FromChannel discord.ChannelID
	ToChannel   discord.ChannelID

	WebhookID discord.WebhookID
	Token     string
}

// Message ...
type Message struct {
	ServerID  discord.GuildID
	ChannelID discord.ChannelID
	MessageID discord.MessageID
	Original  discord.MessageID
	UserID    discord.UserID
}

func (bot *Bot) mirrors(id discord.GuildID) (m []Mirror, err error) {
	err = pgxscan.Select(context.Background(), bot.DB.Pool, &m, "select * from channel_mirror where server_id = $1", id)
	return
}

func (bot *Bot) mirrorFor(ch discord.ChannelID) (m Mirror, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &m, "select * from channel_mirror where from_channel = $1", ch)
	return
}

func (bot *Bot) mirrorTo(ch discord.ChannelID) (m Mirror, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &m, "select * from channel_mirror where to_channel = $1", ch)
	return
}

func (bot *Bot) setMirror(m Mirror) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), `insert into channel_mirror
	(server_id, from_channel, to_channel, webhook_id, token)
	values ($1, $2, $3, $4, $5) on conflict (from_channel) do update
	set to_channel = $3, webhook_id = $4, token = $5`, m.ServerID, m.FromChannel, m.ToChannel, m.WebhookID, m.Token)
	return
}

func (bot *Bot) insertMessage(m Message) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), `insert into channel_mirror_messages
	(server_id, channel_id, message_id, original, user_id)
	values ($1, $2, $3, $4, $5)`, m.ServerID, m.ChannelID, m.MessageID, m.Original, m.UserID)
	return
}

func (bot *Bot) message(id discord.MessageID) (m Message, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &m, "select * from channel_mirror_messages where message_id = $1 or original = $1", id)
	return
}
