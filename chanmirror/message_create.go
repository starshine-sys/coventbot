package chanmirror

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/starshine-sys/pkgo"
)

func (bot *Bot) messageCreate(m *gateway.MessageCreateEvent) {
	if !m.GuildID.IsValid() || m.Author.Bot {
		if m.WebhookID.IsValid() {
			bot.pkMessage(m)
		}

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

	embeds := []discord.Embed{}

	if m.Reference != nil {
		s, _ := bot.Router.StateFromGuildID(m.GuildID)
		ref, err := s.Message(m.Reference.ChannelID, m.Reference.MessageID)
		if err == nil {
			content := fmt.Sprintf("**[Reply to:](https://discord.com/channels/%v/", m.GuildID)

			dbm, err := bot.message(m.Reference.MessageID)
			if err == nil {
				content += fmt.Sprintf("%v/%v)**", dbm.ChannelID, dbm.MessageID)
			} else {
				content += fmt.Sprintf("%v/%v)**", m.Reference.ChannelID, m.Reference.MessageID)
			}

			msgContent := ref.Content
			if len(msgContent) > 100 {
				msgContent = msgContent[:100]

				if strings.Count(msgContent, "||")%2 == 1 && strings.Count(ref.Content, "||")%2 == 0 {
					msgContent += "||"
				}

				msgContent += "..."
			}

			if ref.Content == "" {
				msgContent = "*(attachment)*"
			}

			embeds = append(embeds, discord.Embed{
				Author: &discord.EmbedAuthor{
					Icon: ref.Author.AvatarURLWithType(discord.PNGImage),
					Name: ref.Author.Username + " \u2004\u21a9\ufe0f",
				},
				Description: fmt.Sprintf("%v %v", content, msgContent),
			})
		}
	}

	if len(embeds) > 10 {
		embeds = embeds[:9]
	}

	for i, a := range m.Attachments {
		if hasAnySuffix(a.Filename, ".png", ".jpeg", ".jpg", ".gif", ".webp") {
			embeds = append(embeds, discord.Embed{
				Image: &discord.EmbedImage{
					URL: a.URL,
				},
			})
		} else {
			if len(m.Content)+len("[Attachment]()\n")+len(a.URL) < 1990 {
				m.Content += fmt.Sprintf("[Attachment %v](%v)\n", i+1, a.URL)
			}
		}
	}

	msg, err := client.ExecuteAndWait(webhook.ExecuteData{
		Content:   m.Content,
		Username:  name,
		AvatarURL: m.Author.AvatarURL(),
		AllowedMentions: &api.AllowedMentions{
			Parse: []api.AllowedMentionType{api.AllowUserMention},
		},
		Embeds: embeds,
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

func (bot *Bot) pkMessage(m *gateway.MessageCreateEvent) {
	mirror, err := bot.mirrorFor(m.ChannelID)
	if err != nil {
		return
	}

	pkMsg, err := bot.PK.Message(pkgo.Snowflake(m.ID))
	if err != nil {
		return
	}

	orig, err := bot.message(discord.MessageID(pkMsg.Original))
	if err == nil {
		s, _ := bot.Router.StateFromGuildID(m.GuildID)

		s.DeleteMessage(orig.ChannelID, orig.MessageID)
	}

	client := webhook.New(mirror.WebhookID, mirror.Token)

	embeds := m.Embeds

	if len(embeds) > 10 {
		embeds = embeds[:9]
	}

	for i, a := range m.Attachments {
		if hasAnySuffix(a.Filename, ".png", ".jpeg", ".jpg", ".gif", ".webp") {
			embeds = append(embeds, discord.Embed{
				Image: &discord.EmbedImage{
					URL: a.URL,
				},
			})
		} else {
			if len(m.Content)+len("[Attachment]()\n")+len(a.URL) < 1990 {
				m.Content += fmt.Sprintf("[Attachment %v](%v)\n", i+1, a.URL)
			}
		}
	}

	msg, err := client.ExecuteAndWait(webhook.ExecuteData{
		Content:   m.Content,
		Username:  m.Author.Username,
		AvatarURL: m.Author.AvatarURL(),
		AllowedMentions: &api.AllowedMentions{
			Parse: []api.AllowedMentionType{api.AllowUserMention},
		},
		Embeds: embeds,
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
		UserID:    discord.UserID(pkMsg.Sender),
	})
	if err != nil {
		bot.Sugar.Errorf("Error inserting message: %v", err)
		return
	}
}

func hasAnySuffix(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}

	return false
}
