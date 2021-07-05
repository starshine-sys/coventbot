package bot

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

var _ bcr.CustomPerms = (*AdminRole)(nil)
var _ bcr.CustomPerms = (*ModRole)(nil)
var _ bcr.CustomPerms = (*HelperRole)(nil)

// AdminRole checks if the user has a role with the Administrator permission, or a role with the ADMIN perm level
type AdminRole struct {
	*Bot
}

func (bot *AdminRole) String() string {
	return "Admin"
}

// Check ...
func (bot *AdminRole) Check(ctx *bcr.Context) (bool, error) {
	if ctx.Member == nil || ctx.Guild == nil {
		return false, nil
	}

	if ctx.Author.ID == ctx.Guild.OwnerID {
		return true, nil
	}

	for _, id := range ctx.Member.RoleIDs {
		for _, r := range ctx.Guild.Roles {
			if r.ID == id {
				if r.Permissions.Has(discord.PermissionAdministrator) {
					return true, nil
				}
			}
		}
	}

	var roles []uint64
	err := bot.DB.Pool.QueryRow(context.Background(), "select admin_roles from servers where id = $1", ctx.Guild.ID).Scan(&roles)
	if err != nil {
		return false, err
	}

	for _, r := range roles {
		for _, id := range ctx.Member.RoleIDs {
			if r == uint64(id) {
				return true, nil
			}
		}
	}

	return false, nil
}

// ModRole checks if the user has a role with the Manage Server permission, or a role with the MODERATOR perm level
type ModRole struct {
	*Bot
}

func (bot *ModRole) String() string {
	return "Moderator"
}

// Check ...
func (bot *ModRole) Check(ctx *bcr.Context) (bool, error) {
	if ctx.Member == nil || ctx.Guild == nil {
		return false, nil
	}

	if ctx.Author.ID == ctx.Guild.OwnerID {
		return true, nil
	}

	for _, id := range ctx.Member.RoleIDs {
		for _, r := range ctx.Guild.Roles {
			if r.ID == id {
				if r.Permissions.Has(discord.PermissionManageGuild) {
					return true, nil
				}
			}
		}
	}

	var roles []uint64
	err := bot.DB.Pool.QueryRow(context.Background(), "select mod_roles || admin_roles from servers where id = $1", ctx.Guild.ID).Scan(&roles)
	if err != nil {
		return false, err
	}

	for _, r := range roles {
		for _, id := range ctx.Member.RoleIDs {
			if r == uint64(id) {
				return true, nil
			}
		}
	}

	return false, nil
}

// HelperRole checks if the user has a role with the Manage Messages permission, or a role with the HELPER perm level
type HelperRole struct {
	*Bot
}

func (bot *HelperRole) String() string {
	return "Helper"
}

// Check ...
func (bot *HelperRole) Check(ctx *bcr.Context) (bool, error) {
	if ctx.Member == nil || ctx.Guild == nil {
		return false, nil
	}

	if ctx.Author.ID == ctx.Guild.OwnerID {
		return true, nil
	}

	for _, id := range ctx.Member.RoleIDs {
		for _, r := range ctx.Guild.Roles {
			if r.ID == id {
				if r.Permissions.Has(discord.PermissionManageMessages) {
					return true, nil
				}
			}
		}
	}

	var roles []uint64
	err := bot.DB.Pool.QueryRow(context.Background(), "select helper_roles || mod_roles || admin_roles from servers where id = $1", ctx.Guild.ID).Scan(&roles)
	if err != nil {
		return false, err
	}

	for _, r := range roles {
		for _, id := range ctx.Member.RoleIDs {
			if r == uint64(id) {
				return true, nil
			}
		}
	}

	return false, nil
}
