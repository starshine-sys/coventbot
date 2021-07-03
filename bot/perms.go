package bot

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

// GlobalPerms ...
func (bot *Bot) GlobalPerms(ctx *bcr.Context) (perms discord.Permissions) {
	if ctx.Guild == nil || ctx.Member == nil {
		return discord.PermissionViewChannel | discord.PermissionSendMessages | discord.PermissionAddReactions | discord.PermissionReadMessageHistory
	}

	for _, id := range ctx.Member.RoleIDs {
		for _, r := range ctx.Guild.Roles {
			if id == r.ID {
				perms |= r.Permissions
				break
			}
		}
	}

	return perms
}
