package quotes

import (
	"fmt"
	"sync"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/starshine-sys/pkgo"
)

func (bot *Bot) reactionAdd(ev *gateway.MessageReactionAddEvent) {
	if !ev.GuildID.IsValid() {
		return
	}

	if !bot.quotesEnabled(ev.GuildID) {
		return
	}

	// make sure only one quote gets added for a message
	if bot.mu[ev.MessageID] == nil {
		bot.mu[ev.MessageID] = &sync.Mutex{}
	}
	bot.mu[ev.MessageID].Lock()
	defer bot.mu[ev.MessageID].Unlock()

	_, err := bot.quoteMessage(ev.MessageID)
	if err == nil {
		return
	}

	u, err := bot.State.User(ev.UserID)
	if err != nil || u.Bot {
		return
	}

	if ev.Emoji.Name != "💬" && ev.Emoji.Name != "🗨️" {
		return
	}

	msg, err := bot.State.Message(ev.ChannelID, ev.MessageID)
	if err != nil {
		bot.Sugar.Errorf("Couldn't fetch message: %v", err)
		return
	}

	if msg.Content == "" {
		return
	}

	q := Quote{
		ServerID:  ev.GuildID,
		ChannelID: ev.ChannelID,
		MessageID: ev.MessageID,

		UserID:  msg.Author.ID,
		AddedBy: ev.UserID,

		Content: msg.Content,
	}
	if msg.WebhookID.IsValid() {
		pkMsg, err := bot.PK.Message(pkgo.Snowflake(ev.MessageID))
		if err != nil {
			return
		}

		q.UserID = discord.UserID(pkMsg.Sender)
		q.Proxied = true
	}

	q, err = bot.insertQuote(q)
	if err != nil {
		bot.Sugar.Errorf("Error inserting quote: %v", err)
	}

	err = bot.State.React(ev.ChannelID, ev.MessageID, ev.Emoji.APIString())
	if err != nil {
		bot.Sugar.Errorf("Couldn't react to the quote message: %v", err)
	}

	if bot.suppressMessages(ev.GuildID) {
		return
	}

	perms, err := bot.State.Permissions(ev.ChannelID, ev.UserID)
	if err == nil && !perms.Has(discord.PermissionSendMessages) {
		return
	}

	_, err = bot.State.SendMessageComplex(ev.ChannelID, api.SendMessageData{
		Content: fmt.Sprintf("New quote added with ID `%v` by %v.", q.HID, u.Username),
		AllowedMentions: &api.AllowedMentions{
			Parse: []api.AllowedMentionType{},
		},
		Reference: &discord.MessageReference{
			MessageID: ev.MessageID,
		},
	})
	if err != nil {
		bot.Sugar.Errorf("Error sending message: %v", err)
	}
}
