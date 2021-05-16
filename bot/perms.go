package bot

import "github.com/diamondburned/arikawa/v2/discord"

// GlobalPermissions ...
func (bot *Bot) GlobalPermissions(guildID discord.GuildID, userID discord.UserID) (perms discord.Permissions, err error) {
	g, err := bot.State.Guild(guildID)
	if err != nil {
		return 0, err
	}
	u, err := bot.State.Member(guildID, userID)
	if err != nil {
		return 0, err
	}

	if g.OwnerID == userID {
		return discord.PermissionAll, nil
	}

	// get role permissions
	for _, role := range g.Roles {
		for _, id := range u.RoleIDs {
			if role.ID == id {
				perms |= role.Permissions
			}
		}
	}

	if perms.Has(discord.PermissionAdministrator) {
		return discord.PermissionAll, nil
	}

	return perms, nil
}
