package bot

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

var _ bcr.CustomPerms = (*AdminRole)(nil)
var _ bcr.CustomPerms = (*ManagerRole)(nil)
var _ bcr.CustomPerms = (*ModeratorRole)(nil)

// AdminRole checks if the user has a role with the Administrator permission, or a role with the ADMIN perm level
type AdminRole struct {
	*Bot
}

func (bot *AdminRole) String(ctx bcr.Contexter) string {
	return "Admin"
}

// Check ...
func (bot *AdminRole) Check(ctx bcr.Contexter) (bool, error) {
	if ctx.GetMember() == nil || ctx.GetGuild() == nil {
		return false, nil
	}

	if ctx.User().ID == ctx.GetGuild().OwnerID {
		return true, nil
	}

	for _, id := range ctx.GetMember().RoleIDs {
		for _, r := range ctx.GetGuild().Roles {
			if r.ID == id {
				if r.Permissions.Has(discord.PermissionAdministrator) {
					return true, nil
				}
			}
		}
	}

	var roles []uint64
	err := bot.DB.Pool.QueryRow(context.Background(), "select admin_roles from servers where id = $1", ctx.GetGuild().ID).Scan(&roles)
	if err != nil {
		return false, err
	}

	for _, r := range roles {
		for _, id := range ctx.GetMember().RoleIDs {
			if r == uint64(id) {
				return true, nil
			}
		}
	}

	return false, nil
}

// ManagerRole checks if the user has a role with the Manage Server permission, or a role with the MODERATOR perm level
type ManagerRole struct {
	*Bot
}

func (bot *ManagerRole) String(ctx bcr.Contexter) string {
	return "Moderator"
}

// Check ...
func (bot *ManagerRole) Check(ctx bcr.Contexter) (bool, error) {
	if ctx.GetMember() == nil || ctx.GetGuild() == nil {
		return false, nil
	}

	if ctx.User().ID == ctx.GetGuild().OwnerID {
		return true, nil
	}

	for _, id := range ctx.GetMember().RoleIDs {
		for _, r := range ctx.GetGuild().Roles {
			if r.ID == id {
				if r.Permissions.Has(discord.PermissionManageGuild) {
					return true, nil
				}
			}
		}
	}

	var roles []uint64
	err := bot.DB.Pool.QueryRow(context.Background(), "select manager_roles || admin_roles from servers where id = $1", ctx.GetGuild().ID).Scan(&roles)
	if err != nil {
		return false, err
	}

	for _, r := range roles {
		for _, id := range ctx.GetMember().RoleIDs {
			if r == uint64(id) {
				return true, nil
			}
		}
	}

	return false, nil
}

// ModeratorRole checks if the user has a role with the Manage Messages permission, or a role with the HELPER perm level
type ModeratorRole struct {
	*Bot
}

func (bot *ModeratorRole) String(ctx bcr.Contexter) string {
	return "Helper"
}

// Check ...
func (bot *ModeratorRole) Check(ctx bcr.Contexter) (bool, error) {
	if ctx.GetMember() == nil || ctx.GetGuild() == nil {
		return false, nil
	}

	if ctx.User().ID == ctx.GetGuild().OwnerID {
		return true, nil
	}

	for _, id := range ctx.GetMember().RoleIDs {
		for _, r := range ctx.GetGuild().Roles {
			if r.ID == id {
				if r.Permissions.Has(discord.PermissionManageMessages) {
					return true, nil
				}
			}
		}
	}

	var roles []uint64
	err := bot.DB.Pool.QueryRow(context.Background(), "select moderator_roles || manager_roles || admin_roles from servers where id = $1", ctx.GetGuild().ID).Scan(&roles)
	if err != nil {
		return false, err
	}

	for _, r := range roles {
		for _, id := range ctx.GetMember().RoleIDs {
			if r == uint64(id) {
				return true, nil
			}
		}
	}

	return false, nil
}
