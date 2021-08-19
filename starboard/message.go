package starboard

import (
	"context"
	"fmt"
	"regexp"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/starshine-sys/tribble/db"
	"github.com/starshine-sys/tribble/etc"
)

func (bot *Bot) deleteMessage(
	state *state.State, channelID discord.ChannelID, messageID discord.MessageID, settings db.StarboardSettings, s *db.StarboardMessage) {
	err := state.DeleteMessage(settings.StarboardChannel, s.StarboardMessageID, "")
	if err != nil {
		bot.Sugar.Errorf("Error deleting starboard message: %v", err)
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "delete from starboard_messages where message_id = $1", messageID)
	if err != nil {
		bot.Sugar.Errorf("Error deleting database entry for starboard message: %v", err)
	}
}

func (bot *Bot) starboardMessage(state *state.State, m discord.Message, settings db.StarboardSettings, s *db.StarboardMessage, count int) {
	embed := bot.embed(m)
	msgContent := fmt.Sprintf("**%v** %v <#%v>", count, settings.StarboardEmoji, m.ChannelID)

	// if s is nil, this is a new message
	if s == nil || s.MessageID == 0 {
		msg, err := state.SendMessage(settings.StarboardChannel, msgContent, embed)
		if err != nil {
			bot.Sugar.Errorf("Error sending starboard message: %v", err)
			return
		}

		err = bot.DB.SaveStarboardMessage(db.StarboardMessage{
			MessageID:          m.ID,
			ChannelID:          m.ChannelID,
			ServerID:           m.GuildID,
			StarboardMessageID: msg.ID,
		})
		if err != nil {
			bot.Sugar.Errorf("Error saving starboard message: %v", err)
		}
	} else {
		// otherwise, edit the existing message
		_, err := state.EditMessage(settings.StarboardChannel, s.StarboardMessageID, msgContent, embed)
		if err != nil {
			bot.Sugar.Errorf("Error editing starboard message: %v", err)
			return
		}
	}
}

// embed creates a starboard embed for the given message object
func (bot *Bot) embed(m discord.Message) (e discord.Embed) {
	name := m.Author.Username
	if !m.WebhookID.IsValid() {
		member, err := bot.Member(m.GuildID, m.Author.ID)
		if err == nil && member.Nick != "" {
			name = member.Nick
		}
	}

	var attachmentLink string
	if len(m.Attachments) > 0 {
		match, _ := regexp.MatchString("\\.(png|jpg|jpeg|gif|webp)$", m.Attachments[0].URL)
		if match {
			attachmentLink = m.Attachments[0].URL
		}
	}

	return discord.Embed{
		Description: m.Content,
		Author: &discord.EmbedAuthor{
			Name: name,
			Icon: m.Author.AvatarURL(),
		},
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v", m.ID),
		},
		Fields: []discord.EmbedField{
			{
				Name:  "Source",
				Value: fmt.Sprintf("[Jump to message](https://discord.com/channels/%v/%v/%v)", m.GuildID, m.ChannelID, m.ID),
			},
		},
		Timestamp: discord.Timestamp(m.Timestamp.Time()),
		Color:     etc.ColourGold,
		Image: &discord.EmbedImage{
			URL: attachmentLink,
		},
	}
}
