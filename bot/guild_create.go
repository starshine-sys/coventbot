// SPDX-License-Identifier: AGPL-3.0-only
package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/starshine-sys/bcr"
)

// GuildCreate logs the bot joining a server, and creates a database entry if one doesn't exist
func (bot *Bot) GuildCreate(g *gateway.GuildCreateEvent) {
	// create the server if it doesn't exist
	exists, err := bot.DB.CreateServerIfNotExists(g.ID)
	if err != nil {
		bot.Sugar.Errorf("Error creating database entry for server: %v", err)
	}

	err = bot.setRolePerms(g)
	if err != nil {
		bot.Sugar.Errorf("Error setting permissions for server: %v", err)
	}

	// if we already joined the server, don't log the join
	if exists && g.Joined.Time().Before(time.Now().Add(-1*time.Minute)) {
		return
	}

	bot.Sugar.Infof("Joined server %v (%v).", g.Name, g.ID)

	s, _ := bot.Router.StateFromGuildID(g.ID)

	botUser, err := s.Me()
	if err != nil {
		bot.Sugar.Errorf("Error getting bot user: %v", err)
	}

	owner := g.OwnerID.Mention()
	if o, err := s.User(g.OwnerID); err == nil {
		owner = fmt.Sprintf("%v#%v (%v)", o.Username, o.Discriminator, o.Mention())
	}

	if bot.GuildLogWebhook != nil {
		bot.GuildLogWebhook.Execute(webhook.ExecuteData{
			Username:  fmt.Sprintf("%v server join", botUser.Username),
			AvatarURL: botUser.AvatarURL(),

			Embeds: []discord.Embed{{
				Title: "Joined new server",
				Color: bcr.ColourBlurple,
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

func (bot *Bot) setRolePerms(g *gateway.GuildCreateEvent) (err error) {
	var rolesSetUp bool
	err = bot.DB.Pool.QueryRow(context.Background(), "select roles_set_up from servers where id = $1", g.ID).Scan(&rolesSetUp)
	if err != nil || rolesSetUp {
		return
	}

	var (
		helperRoles = []uint64{}
		modRoles    = []uint64{}
		adminRoles  = []uint64{}
	)

	for _, r := range g.Roles {
		if r.Managed {
			continue
		}

		if r.Permissions.Has(discord.PermissionAdministrator) {
			adminRoles = append(adminRoles, uint64(r.ID))
		} else if r.Permissions.Has(discord.PermissionManageGuild) {
			modRoles = append(modRoles, uint64(r.ID))
		} else if r.Permissions.Has(discord.PermissionManageMessages) {
			helperRoles = append(helperRoles, uint64(r.ID))
		}
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set moderator_roles = $1, manager_roles = $2, admin_roles = $3, roles_set_up = true where id = $4", helperRoles, modRoles, adminRoles, g.ID)
	if err == nil {
		bot.Sugar.Infof("Set helper/mod/admin roles for server %v (%v)", g.Name, g.ID)
	}
	return
}
