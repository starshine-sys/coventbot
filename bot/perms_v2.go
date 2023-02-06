// SPDX-License-Identifier: AGPL-3.0-only
package bot

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr/v2"
	"github.com/starshine-sys/tribble/common"
)

func (bot *Bot) CheckAdmin(ctx *bcr.CommandContext) (err error) {
	if ctx.Member == nil || ctx.Guild == nil {
		return bcr.NewCheckError[*bcr.CommandContext]("This command requires an admin role.")
	}

	if ctx.User.ID == ctx.Guild.OwnerID {
		return nil
	}

	for _, id := range ctx.Member.RoleIDs {
		for _, r := range ctx.Guild.Roles {
			if r.ID == id {
				if r.Permissions.Has(discord.PermissionAdministrator) {
					return nil
				}
			}
		}
	}

	var roles []uint64
	err = bot.DB.Pool.QueryRow(context.Background(), "select admin_roles from servers where id = $1", ctx.Guild.ID).Scan(&roles)
	if err != nil {
		return err
	}

	for _, r := range roles {
		for _, id := range ctx.Member.RoleIDs {
			if r == uint64(id) {
				return nil
			}
		}
	}

	return bcr.NewCheckError[*bcr.CommandContext]("This command requires an admin role.")
}

func (bot *Bot) CheckMod(ctx *bcr.CommandContext) (err error) {
	if ctx.Member == nil || ctx.Guild == nil {
		return bcr.NewCheckError[*bcr.CommandContext]("This command requires a moderator role.")
	}

	if ctx.User.ID == ctx.Guild.OwnerID {
		return nil
	}

	for _, id := range ctx.Member.RoleIDs {
		for _, r := range ctx.Guild.Roles {
			if r.ID == id {
				if r.Permissions.Has(discord.PermissionManageGuild) {
					return nil
				}
			}
		}
	}

	var roles []uint64
	err = bot.DB.Pool.QueryRow(context.Background(), "select manager_roles || admin_roles from servers where id = $1", ctx.Guild.ID).Scan(&roles)
	if err != nil {
		return err
	}

	for _, r := range roles {
		for _, id := range ctx.Member.RoleIDs {
			if r == uint64(id) {
				return nil
			}
		}
	}

	return bcr.NewCheckError[*bcr.CommandContext]("This command requires a moderator role.")
}

func (bot *Bot) CheckHelper(ctx *bcr.CommandContext) (err error) {
	if ctx.Member == nil || ctx.Guild == nil {
		return bcr.NewCheckError[*bcr.CommandContext]("This command requires a helper role.")
	}

	if ctx.User.ID == ctx.Guild.OwnerID {
		return nil
	}

	for _, id := range ctx.Member.RoleIDs {
		for _, r := range ctx.Guild.Roles {
			if r.ID == id {
				if r.Permissions.Has(discord.PermissionManageMessages) {
					return nil
				}
			}
		}
	}

	var roles []uint64
	err = bot.DB.Pool.QueryRow(context.Background(), "select moderator_roles || manager_roles || admin_roles from servers where id = $1", ctx.Guild.ID).Scan(&roles)
	if err != nil {
		return err
	}

	for _, r := range roles {
		for _, id := range ctx.Member.RoleIDs {
			if r == uint64(id) {
				return nil
			}
		}
	}

	return bcr.NewCheckError[*bcr.CommandContext]("This command requires a helper role.")
}

func (bot *Bot) RequireNode(required string) bcr.Check[*bcr.CommandContext] {
	return func(ctx *bcr.CommandContext) (err error) {
		if ctx.Guild == nil {
			node := common.DefaultPermissions.NodeFor(required)
			if node.Level == common.EveryoneLevel {
				return nil
			}
			return bcr.NewCheckError[*bcr.CommandContext]("❌ This command cannot be run in DMs.")
		}

		userLevel := bot.UserBotPermissions(ctx.User, ctx.Member, ctx.Guild)
		node := bot.NodeLevel(ctx.Guild.ID, required)

		if node.Level == common.DisabledLevel && userLevel != common.AdminLevel {
			return bcr.NewCheckError[*bcr.CommandContext]("❌ This command is disabled.")
		}

		if userLevel < node.Level {
			return bcr.NewCheckError[*bcr.CommandContext]("", bot.permError(required, node.Level, userLevel).Embeds...)
		}

		return nil
	}
}
