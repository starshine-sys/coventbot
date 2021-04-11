package pklog

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/api/webhook"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/starshine-sys/tribble/etc"
)

func (bot *Bot) messageDelete(m *gateway.MessageDeleteEvent) {
	var logChannel discord.ChannelID
	bot.DB.Pool.QueryRow(context.Background(), "select pk_log_channel from servers where id = $1", m.GuildID).Scan(&logChannel)
	if !logChannel.IsValid() {
		return
	}

	// try getting the message
	msg, err := bot.Get(m.ID)
	if err != nil {
		return
	}

	// try getting the cached webhook
	var wh *discord.Webhook

	w, err := bot.GetWebhooks(m.GuildID)
	if err != nil {
		wh, err = bot.getWebhook(logChannel, bot.Router.Bot.Username+" Logging")
		if err != nil {
			bot.Sugar.Errorf("Error getting webhook: %v", err)
			return
		}

		bot.SetWebhooks(m.GuildID, &Webhooks{
			MessageWebhookID:    wh.ID,
			MessageWebhookToken: wh.Token,
		})
	} else {
		wh = &discord.Webhook{
			ID:    w.MessageWebhookID,
			Token: w.MessageWebhookToken,
		}
	}

	mention := msg.UserID.Mention()
	var author *discord.EmbedAuthor
	u, err := bot.State.User(msg.UserID)
	if err == nil {
		mention = fmt.Sprintf("%v\n%v#%v\nID: %v", u.Mention(), u.Username, u.Discriminator, u.ID)
		author = &discord.EmbedAuthor{
			Icon: u.AvatarURL(),
			Name: u.Username + u.Discriminator,
		}
	}

	e := discord.Embed{
		Author:      author,
		Title:       fmt.Sprintf("Message by \"%v\" deleted", msg.Username),
		Description: msg.Content,
		Color:       etc.ColourRed,
		Fields: []discord.EmbedField{
			{
				Name:  "​",
				Value: "​",
			},
			{
				Name:   "Channel",
				Value:  fmt.Sprintf("%v\nID: %v", msg.ChannelID.Mention(), msg.ChannelID),
				Inline: true,
			},
			{
				Name:   "Linked Discord account",
				Value:  mention,
				Inline: true,
			},
			{
				Name:  "​",
				Value: "​",
			},
			{
				Name:   "System ID",
				Value:  msg.System,
				Inline: true,
			},
			{
				Name:   "Member ID",
				Value:  msg.Member,
				Inline: true,
			},
		},
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v", msg.MsgID),
		},
		Timestamp: discord.NewTimestamp(msg.MsgID.Time()),
	}

	_, err = webhook.New(wh.ID, wh.Token).ExecuteAndWait(webhook.ExecuteData{
		AvatarURL: bot.Router.Bot.AvatarURL(),
		Embeds:    []discord.Embed{e},
	})
	if err == nil {
		bot.Delete(msg.MsgID)
	}
}

func (bot *Bot) getWebhook(id discord.ChannelID, name string) (*discord.Webhook, error) {
	ws, err := bot.State.ChannelWebhooks(id)
	if err == nil {
		for _, w := range ws {
			if w.Name == name {
				return &w, nil
			}
		}
	} else {
		return nil, err
	}

	w, err := bot.State.CreateWebhook(id, api.CreateWebhookData{
		Name: name,
	})
	return w, err
}
