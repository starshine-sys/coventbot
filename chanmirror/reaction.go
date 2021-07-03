package chanmirror

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) reactionAdd(ev *gateway.MessageReactionAddEvent) {
	msgInfo, err := bot.message(ev.MessageID)
	if err != nil {
		return
	}

	s, _ := bot.Router.StateFromGuildID(ev.GuildID)

	if ev.Emoji.Name == "❌" && ev.UserID == msgInfo.UserID {
		err = s.DeleteMessage(ev.ChannelID, ev.MessageID)
		if err != nil {
			bot.Sugar.Errorf("Error deleting message: %v", err)
		}
		return
	}

	if ev.Emoji.Name == "❓" || ev.Emoji.Name == "❔" {
		s.DeleteUserReaction(ev.ChannelID, ev.MessageID, ev.UserID, ev.Emoji.APIString())

		ch, err := s.CreatePrivateChannel(ev.UserID)
		if err != nil {
			return
		}

		e := discord.Embed{
			Color: bcr.ColourBlurple,
			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("ID: %v", msgInfo.MessageID),
			},
			Timestamp: discord.NewTimestamp(msgInfo.MessageID.Time()),
		}

		msg, err := s.Message(msgInfo.ChannelID, msgInfo.MessageID)
		if err == nil {
			e.Description = msg.Content
		}

		mr, err := bot.mirrorTo(msgInfo.ChannelID)
		if err != nil {
			e.Fields = append(e.Fields, discord.EmbedField{
				Name:   "Original message",
				Value:  msgInfo.Original.String(),
				Inline: true,
			})
		} else {
			e.Fields = append(e.Fields, discord.EmbedField{
				Name:   "Original message",
				Value:  fmt.Sprintf("ID: %v\nChannel: %v", msgInfo.Original, mr.FromChannel.Mention()),
				Inline: true,
			})
		}

		u, err := s.User(msgInfo.UserID)
		if err != nil {
			e.Author = &discord.EmbedAuthor{
				Name: fmt.Sprintf("[unknown user %v]", msgInfo.UserID),
			}

			e.Fields = append(e.Fields, discord.EmbedField{
				Name:   "Original sender",
				Value:  fmt.Sprintf("*[unknown user %v]*", msgInfo.UserID),
				Inline: true,
			})
		} else {
			e.Author = &discord.EmbedAuthor{
				Name: fmt.Sprintf("%v#%v", u.Username, u.Discriminator),
				Icon: u.AvatarURL(),
			}

			e.Fields = append(e.Fields, discord.EmbedField{
				Name:   "Original sender",
				Value:  fmt.Sprintf("%v#%v\n%v", u.Username, u.Discriminator, u.Mention()),
				Inline: true,
			})
		}

		_, err = s.SendEmbeds(ch.ID, e)
		if err != nil {
			bot.Sugar.Errorf("Error sending message: %v", err)
		}
		return
	}
}
