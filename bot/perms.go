package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/common"
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

func (bot *Bot) CheckPermissions(ctx *bcr.Context) (name string, allowed bool, data api.SendMessageData) {
	rawPath := ctx.FullCommandPath
	if strings.EqualFold(ctx.FullCommandPath[0], "help") {
		rawPath = ctx.Args
	}

	var path []string
	c := bot.Router.GetCommand(rawPath[0])
	path = append(path, c.Name)
	if len(rawPath) > 1 {
		for _, p := range rawPath[1:] {
			c = c.GetCommand(p)
			path = append(path, c.Name)
		}
	}
	commandPath := strings.Join(path, ".")

	if ctx.Guild == nil {
		node := common.DefaultPermissions.NodeFor(commandPath)
		switch node.Level {
		case common.EveryoneLevel:
			return "`" + common.EveryoneLevel.String() + "`", true, api.SendMessageData{}
		default:
			return "`" + node.Level.String() + "`", false, api.SendMessageData{Content: "❌ This command cannot be run in DMs."}
		}
	}

	nodes, err := bot.DB.Permissions(ctx.Guild.ID)
	if err != nil {
		bot.Sugar.Errorf("getting permission nodes for guild %v: %v", ctx.Guild.ID, err)
		return "`" + common.DisabledLevel.String() + "`", false, api.SendMessageData{
			Content: "Error checking your permission level.",
		}
	}

	userLevel := bot.userPermissions(ctx)
	node := nodes.NodeFor(commandPath)

	if node.Level == common.DisabledLevel && userLevel != common.AdminLevel {
		return "`" + common.DisabledLevel.String() + "`", false, api.SendMessageData{Content: "This command is disabled."}
	}

	if userLevel < node.Level {
		return "`" + node.Level.String() + "`", false, bot.permError(node.Name, node.Level, userLevel)
	}

	return "`" + node.Level.String() + "`", true, api.SendMessageData{}
}

func (bot *Bot) permError(
	effectiveNode string,
	requiredLevel common.PermissionLevel,
	userLevel common.PermissionLevel,
) api.SendMessageData {
	return api.SendMessageData{
		Embeds: []discord.Embed{{
			Title: "Missing permissions",
			Description: fmt.Sprintf(
				"You're not allowed to use this command. This command needs `%s` permissions, but you only have `%s`.",
				requiredLevel, userLevel,
			),
			Color: bcr.ColourRed,
			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("Effective permission node: %v", effectiveNode),
			},
		}},
	}
}

func (bot *Bot) userPermissions(ctx *bcr.Context) common.PermissionLevel {
	if ctx.Author.ID == ctx.Guild.OwnerID {
		return common.AdminLevel
	}

	// check admin perms
	for _, id := range ctx.Member.RoleIDs {
		for _, r := range ctx.Guild.Roles {
			if r.ID == id {
				if r.Permissions.Has(discord.PermissionAdministrator) {
					return common.AdminLevel
				}
			}
		}
	}

	// check manage guild
	for _, id := range ctx.Member.RoleIDs {
		for _, r := range ctx.Guild.Roles {
			if r.ID == id {
				if r.Permissions.Has(discord.PermissionManageGuild) {
					return common.ManagerLevel
				}
			}
		}
	}

	// check manage messages
	for _, id := range ctx.Member.RoleIDs {
		for _, r := range ctx.Guild.Roles {
			if r.ID == id {
				if r.Permissions.Has(discord.PermissionManageMessages) {
					return common.ModeratorLevel
				}
			}
		}
	}

	var moderator, manager, admin []uint64
	err := bot.DB.Pool.QueryRow(
		context.Background(),
		"select moderator_roles, manager_roles, admin_roles from servers where id = $1",
		ctx.GetGuild().ID,
	).Scan(&moderator, &manager, &admin)
	if err != nil {
		bot.Sugar.Errorf("error geting role overrides for guild %v: %v", ctx.Guild.ID, err)
		return common.EveryoneLevel
	}

	for _, r := range admin {
		for _, id := range ctx.Member.RoleIDs {
			if r == uint64(id) {
				return common.AdminLevel
			}
		}
	}

	for _, r := range manager {
		for _, id := range ctx.Member.RoleIDs {
			if r == uint64(id) {
				return common.ManagerLevel
			}
		}
	}

	for _, r := range moderator {
		for _, id := range ctx.Member.RoleIDs {
			if r == uint64(id) {
				return common.ModeratorLevel
			}
		}
	}

	return common.EveryoneLevel
}
