package bot

import (
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v2/api/webhook"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/starshine-sys/tribble/etc"
)

// GuildCreate logs the bot joining a server, and creates a database entry if one doesn't exist
func (bot *Bot) GuildCreate(g *gateway.GuildCreateEvent) {
	// create the server if it doesn't exist
	exists, err := bot.DB.CreateServerIfNotExists(g.ID)
	if err != nil {
		bot.Sugar.Errorf("Error creating database entry for server: %v", err)
	}

	// if we already joined the server, don't log the join
	if exists && g.Joined.Time().Before(time.Now().Add(-1*time.Minute)) {
		return
	}

	bot.Sugar.Infof("Joined server %v (%v).", g.Name, g.ID)

	botUser, err := bot.State.Me()
	if err != nil {
		bot.Sugar.Errorf("Error getting bot user: %v", err)
	}

	owner := g.OwnerID.Mention()
	if o, err := bot.State.User(g.OwnerID); err == nil {
		owner = fmt.Sprintf("%v#%v (%v)", o.Username, o.Discriminator, o.Mention())
	}

	if bot.GuildLogWebhook != nil {
		bot.GuildLogWebhook.Execute(webhook.ExecuteData{
			Username:  fmt.Sprintf("%v server join", botUser.Username),
			AvatarURL: botUser.AvatarURL(),

			Embeds: []discord.Embed{{
				Title: "Joined new server",
				Color: etc.ColourBlurple,
				Thumbnail: &discord.EmbedThumbnail{
					URL: g.IconURL(),
				},

				Description: fmt.Sprintf("Joined new server **%v**", g.Name),

				Fields: []discord.EmbedField{
					{
						Name:   "Owner",
						Value:  owner,
						Inline: true,
					},
					{
						Name:   "Members",
						Value:  fmt.Sprintf("%v", g.MemberCount),
						Inline: true,
					},
				},

				Footer: &discord.EmbedFooter{
					Text: fmt.Sprintf("ID: %v", g.ID),
				},
				Timestamp: discord.NowTimestamp(),
			}},
		})
	}
	return
}
