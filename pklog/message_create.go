package pklog

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/starshine-sys/pkgo"
)

// messageCreate is used as a backup for pkMessageCreate in case proxy logging isn't enabled.
func (bot *Bot) messageCreate(m *gateway.MessageCreateEvent) {
	var shouldLog bool
	bot.DB.Pool.QueryRow(context.Background(), "select (pk_log_channel != 0) from servers where id = $1", m.GuildID).Scan(&shouldLog)
	if !shouldLog {
		return
	}

	// only check webhook messages
	if !m.WebhookID.IsValid() {
		return
	}

	// wait 5 seconds
	time.Sleep(5 * time.Second)

	// check if the message exists in the database; if so, return
	_, err := bot.Get(m.ID)
	if err == nil {
		return
	}

	pkm, err := bot.PK.Message(pkgo.Snowflake(m.ID))
	if err != nil {
		// Message is either not proxied or we got an error from the PK API. Either way, return
		return
	}

	msg := Message{
		MsgID:     m.ID,
		UserID:    discord.UserID(pkm.Sender),
		ChannelID: m.ChannelID,
		ServerID:  m.GuildID,

		Username: m.Author.Username,
		Member:   pkm.Member.ID,
		System:   pkm.System.ID,

		Content: m.Content,
	}

	// insert the message, ignore errors as those shouldn't impact anything
	bot.Insert(msg)
}
