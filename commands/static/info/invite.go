package info

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) invite(ctx *bcr.Context) (err error) {
	perms := discord.PermissionViewChannel |
		discord.PermissionReadMessageHistory |
		discord.PermissionAddReactions |
		discord.PermissionAttachFiles |
		discord.PermissionCreateInstantInvite |
		discord.PermissionUseExternalEmojis |
		discord.PermissionEmbedLinks |
		discord.PermissionManageEmojis |
		discord.PermissionManageMessages |
		discord.PermissionManageRoles |
		discord.PermissionSendMessages |
		discord.PermissionViewAuditLog |
		discord.PermissionBanMembers |
		discord.PermissionManageWebhooks |
		discord.PermissionManageChannels |
		discord.PermissionMentionEveryone

	invite := func(u discord.UserID, p discord.Permissions) string {
		return fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%v&permissions=%v&scope=applications.commands%%20bot", u, p)
	}

	if bot.Config.Branding.Private {
		if bot.Config.Branding.PublicID.IsValid() {
			u, err := ctx.State.User(bot.Config.Branding.PublicID)
			if err == nil {
				_, err = ctx.Send(fmt.Sprintf("This instance of the bot is private, but you can invite %v, the public version of this bot.", u.Username), discord.Embed{
					Title: "Invite",
					Description: fmt.Sprintf("[Invite link (recommended)](%v)\n\n[Invite link (admin)](%v)",
						invite(u.ID, perms), invite(u.ID, discord.PermissionAdministrator)),
					Color: ctx.Router.EmbedColor,
				})
				return err
			}
		}

		_, err = ctx.Send("This instance of the bot is private, please DM the bot's owner for details.")
		return
	}

	_, err = ctx.Send("", discord.Embed{
		Title:       "Invite",
		Description: fmt.Sprintf("[Invite link (recommended)](%v)\n\n[Invite link (admin)](%v)", invite(ctx.Bot.ID, perms), invite(ctx.Bot.ID, discord.PermissionAdministrator)),
		Color:       ctx.Router.EmbedColor,
	})
	return
}
