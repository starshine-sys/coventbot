package bot

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/api/webhook"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) guildDelete(ev *gateway.GuildDeleteEvent) {
	if ev.Unavailable {
		return
	}

	var e discord.Embed

	if g, err := bot.State.Guild(ev.ID); err != nil {
		e = discord.Embed{
			Color:       bcr.ColourBlurple,
			Title:       "Left unknown server",
			Description: fmt.Sprintf("Left server **%v**", ev.ID),
		}
	} else {
		owner := g.OwnerID.Mention()
		if o, err := bot.State.User(g.OwnerID); err == nil {
			owner = fmt.Sprintf("%v#%v (%v)", o.Username, o.Discriminator, o.Mention())
		}

		e = discord.Embed{
			Title: "Left server",
			Color: bcr.ColourBlurple,
			Thumbnail: &discord.EmbedThumbnail{
				URL: g.IconURL(),
			},

			Description: fmt.Sprintf("Left server **%v**", g.Name),

			Fields: []discord.EmbedField{{
				Name:  "Owner",
				Value: owner,
			}},

			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("ID: %v", g.ID),
			},
			Timestamp: discord.NowTimestamp(),
		}
	}

	if bot.GuildLogWebhook != nil {
		botUser, err := bot.State.Me()
		if err != nil {
			bot.Sugar.Errorf("Error getting bot user: %v", err)
		}

		bot.GuildLogWebhook.Execute(webhook.ExecuteData{
			Username:  fmt.Sprintf("%v server leave", botUser.Username),
			AvatarURL: botUser.AvatarURL(),

			Embeds: []discord.Embed{e},
		})
	}
}
