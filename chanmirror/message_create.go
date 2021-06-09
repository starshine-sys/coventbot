package chanmirror

import (
	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/api/webhook"
	"github.com/diamondburned/arikawa/v2/gateway"
)

func (bot *Bot) messageCreate(m *gateway.MessageCreateEvent) {
	if !m.GuildID.IsValid() || m.Author.Bot {
		return
	}

	mirror, err := bot.mirrorFor(m.ChannelID)
	if err != nil {
		return
	}

	client := webhook.New(mirror.WebhookID, mirror.Token)

	name := m.Author.Username
	if m.Member != nil {
		if m.Member.Nick != "" {
			name = m.Member.Nick
		}
	}

	msg, err := client.ExecuteAndWait(webhook.ExecuteData{
		Content:   m.Content,
		Username:  name,
		AvatarURL: m.Author.AvatarURL(),
		AllowedMentions: &api.AllowedMentions{
			Parse: []api.AllowedMentionType{api.AllowUserMention},
		},
	})
	if err != nil {
		bot.Sugar.Errorf("Error sending mirror message: %v", err)
		return
	}

	err = bot.insertMessage(Message{
		ServerID:  m.GuildID,
		ChannelID: msg.ChannelID,
		MessageID: msg.ID,
		Original:  m.ID,
		UserID:    m.Author.ID,
	})
	if err != nil {
		bot.Sugar.Errorf("Error inserting message: %v", err)
		return
	}
}
