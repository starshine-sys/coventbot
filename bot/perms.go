// SPDX-License-Identifier: AGPL-3.0-only
package bot

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/common"
)

func (bot *Bot) CheckPermissions(ctx *bcr.Context, routing bool) (name string, allowed bool, data api.SendMessageData) {
	rawPath := ctx.FullCommandPath
	// if we're not actually checking permissions (and only running this function for the permission name)
	// check the help command's arguments
	if !routing && strings.EqualFold(ctx.FullCommandPath[0], "help") && len(ctx.Args) > 0 {
		rawPath = ctx.Args
	}

	var path []string
	c := bot.Router.GetCommand(rawPath[0])
	if c == nil {
		return common.DisabledLevel.String(), false, api.SendMessageData{
			Content: fmt.Sprintf("❌ No command named %v found.", bcr.AsCode(rawPath[0])),
		}
	}

	path = append(path, c.Name)
	if len(rawPath) > 1 {
		for _, p := range rawPath[1:] {
			c = c.GetCommand(p)
			if c == nil {
				return common.DisabledLevel.String(), false, api.SendMessageData{
					Content: fmt.Sprintf("❌ No command named %v found.", bcr.AsCode(
						strings.Join(rawPath, " "),
					)),
				}
			}

			path = append(path, c.Name)
		}
	}
	commandPath := strings.Join(path, ".")

	if ctx.Guild == nil {
		node := common.DefaultPermissions.NodeFor(commandPath)
		if node.Level == common.EveryoneLevel {
			return "`" + common.EveryoneLevel.String() + "`", true, api.SendMessageData{}
		}
		return "`" + node.Level.String() + "`", false, api.SendMessageData{Content: "❌ This command cannot be run in DMs."}
	}

	userLevel := bot.UserBotPermissions(ctx.Author, ctx.Member, ctx.Guild)
	node := bot.NodeLevel(ctx.Guild.ID, commandPath)

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

func (bot *Bot) UserBotPermissions(user discord.User, member *discord.Member, guild *discord.Guild) common.PermissionLevel {
	if user.ID == guild.OwnerID {
		return common.AdminLevel
	}

	// check admin perms
	for _, id := range member.RoleIDs {
		for _, r := range guild.Roles {
			if r.ID == id {
				if r.Permissions.Has(discord.PermissionAdministrator) {
					return common.AdminLevel
				}
			}
		}
	}

	// check manage guild
	for _, id := range member.RoleIDs {
		for _, r := range guild.Roles {
			if r.ID == id {
				if r.Permissions.Has(discord.PermissionManageGuild) {
					return common.ManagerLevel
				}
			}
		}
	}

	// check manage messages
	for _, id := range member.RoleIDs {
		for _, r := range guild.Roles {
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
		guild.ID,
	).Scan(&moderator, &manager, &admin)
	if err != nil {
		bot.Sugar.Errorf("error geting role overrides for guild %v: %v", guild.ID, err)
		return common.EveryoneLevel
	}

	for _, r := range admin {
		for _, id := range member.RoleIDs {
			if r == uint64(id) {
				return common.AdminLevel
			}
		}
	}

	for _, r := range manager {
		for _, id := range member.RoleIDs {
			if r == uint64(id) {
				return common.ManagerLevel
			}
		}
	}

	for _, r := range moderator {
		for _, id := range member.RoleIDs {
			if r == uint64(id) {
				return common.ModeratorLevel
			}
		}
	}

	return common.EveryoneLevel
}

func (bot *Bot) NodeLevel(guildID discord.GuildID, node string) common.Node {
	if !guildID.IsValid() {
		return common.DefaultPermissions.NodeFor(node)
	}

	nodes, err := bot.DB.Permissions(guildID)
	if err != nil {
		bot.Sugar.Errorf("getting permission nodes for %v: %v", guildID, err)
		return common.DefaultPermissions.NodeFor(node)
	}
	return nodes.NodeFor(node)
}

func (bot *Bot) InitValidPermissionNodes() {
	var nodes []string

	cmds := bot.Router.Commands()
	for _, c := range cmds {
		if c.Hidden {
			continue
		}

		nodes = append(nodes, bot.recurseSubcommandNodes(c, "", 0)...)
	}

	nnodes := make(common.Nodes, len(nodes))
	for i := range nodes {
		nnodes[i] = common.Node{Name: nodes[i]}
	}

	sort.Sort(nnodes)
	for i := range nnodes {
		nodes[i] = nnodes[i].Name
	}

	bot.ValidNodes = nodes
}

func (bot *Bot) recurseSubcommandNodes(c *bcr.Command, prefix string, i int) (nodes []string) {
	if i > 10 {
		bot.Sugar.Warn("recurseSubCommandNodes: recursion exceeded 10")
		return nil
	}

	subCmds := c.Subcommands()
	if len(subCmds) == 0 {
		return []string{prefix + c.Name}
	}

	for _, sub := range subCmds {
		if c.Hidden {
			continue
		}

		nodes = append(nodes, bot.recurseSubcommandNodes(sub, prefix+c.Name+".", i+1)...)
	}

	nodes = append(nodes, prefix+c.Name, prefix+c.Name+".*")
	return nodes
}
