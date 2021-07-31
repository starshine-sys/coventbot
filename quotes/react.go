package quotes

import (
	"fmt"
	"sync"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
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
	bot.muMu.Lock()
	mu := bot.mu[ev.MessageID]
	if mu == nil {
		mu = &sync.Mutex{}
		bot.mu[ev.MessageID] = mu
	}
	bot.muMu.Unlock()
	mu.Lock()
	defer mu.Unlock()

	_, err := bot.quoteMessage(ev.MessageID)
	if err == nil {
		return
	}

	s, _ := bot.Router.StateFromGuildID(ev.GuildID)

	u, err := s.User(ev.UserID)
	if err != nil || u.Bot {
		return
	}

	if ev.Emoji.Name != "💬" && ev.Emoji.Name != "🗨️" {
		return
	}

	msg, err := s.Message(ev.ChannelID, ev.MessageID)
	if err != nil {
		bot.Sugar.Errorf("Couldn't fetch message: %v", err)
		return
	}

	content := msg.Content
	if content == "" {
		if len(msg.Embeds) > 0 {
			if msg.Embeds[0].Description != "" {
				content = msg.Embeds[0].Description
			}
		} else if len(msg.Attachments) > 0 {
			content = fmt.Sprintf("*[(click to see attachment)](https://discord.com/channels/%v/%v/%v)*", ev.GuildID, ev.ChannelID, ev.MessageID)
		} else {
			return
		}
	}

	q := Quote{
		ServerID:  ev.GuildID,
		ChannelID: ev.ChannelID,
		MessageID: ev.MessageID,

		UserID:  msg.Author.ID,
		AddedBy: ev.UserID,

		Content: content,
	}
	if msg.WebhookID.IsValid() {
		pkMsg, err := bot.PK.Message(pkgo.Snowflake(ev.MessageID))
		if err != nil {
			return
		}

		q.UserID = discord.UserID(pkMsg.Sender)
		q.Proxied = true
	}

	blocked := bot.isUserBlocked(q.UserID)
	if blocked {
		s.DeleteUserReaction(ev.ChannelID, ev.MessageID, ev.UserID, ev.Emoji.APIString())
		return
	}

	q, err = bot.insertQuote(q)
	if err != nil {
		bot.Sugar.Errorf("Error inserting quote: %v", err)
	}

	err = s.React(ev.ChannelID, ev.MessageID, ev.Emoji.APIString())
	if err != nil {
		bot.Sugar.Errorf("Couldn't react to the quote message: %v", err)
	}

	if bot.suppressMessages(ev.GuildID) {
		return
	}

	perms, err := s.Permissions(ev.ChannelID, ev.UserID)
	if err == nil && !perms.Has(discord.PermissionSendMessages) {
		return
	}

	_, err = s.SendMessageComplex(ev.ChannelID, api.SendMessageData{
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
