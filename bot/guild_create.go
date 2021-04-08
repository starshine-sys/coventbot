package bot

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/api/webhook"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/starshine-sys/tribble/etc"
)

// GuildCreate logs the bot joining a server, and creates a database entry if one doesn't exist
func (bot *Bot) GuildCreate(g *gateway.GuildCreateEvent) {
	// create the server if it doesn't exist
	exists, err := bot.DB.CreateServerIfNotExists(g.ID)
	// if the server exists, don't log the join
	if exists {
		return
	}
	if err != nil {
		bot.Sugar.Errorf("Error creating database entry for server: %v", err)
		return
	}

	bot.Sugar.Infof("Joined server %v (%v).", g.Name, g.ID)

	botUser, _ := bot.State.Me()

	owner := g.OwnerID.Mention()
	if o, err := bot.State.User(g.OwnerID); err == nil {
		owner = fmt.Sprintf("%v#%v (%v)", o.Username, o.Discriminator, o.Mention())
	}

	if bot.GuildLogWebhook.ID.IsValid() {
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

				Fields: []discord.EmbedField{{
					Name:  "Owner",
					Value: owner,
				}},

				Footer: &discord.EmbedFooter{
					Text: fmt.Sprintf("ID: %v", g.ID),
				},
				Timestamp: discord.NowTimestamp(),
			}},
		})
	}
	return
}
