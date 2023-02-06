// SPDX-License-Identifier: AGPL-3.0-only
package config

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/common"
)

func (bot *Bot) showNodes(ctx *bcr.Context) (err error) {
	edited, _ := ctx.Flags.GetBool("edited")
	if edited {
		nodes, err := bot.DB.Permissions(ctx.Guild.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		if len(nodes) == 0 {
			return ctx.SendX("There are no nodes with custom permission levels.")
		}

		strs := make([]string, len(nodes))
		for i := range nodes {
			strs[i] = fmt.Sprintf("`%v`: `%v`\n", nodes[i].Name, nodes[i].Level.String())
		}

		_, err = bot.PagedEmbed(ctx, bcr.StringPaginator("Node overrides for "+ctx.Guild.Name, bot.EmbedColour, strs, 15), 15*time.Minute)
		return err
	}

	strs := make([]string, len(bot.ValidNodes))
	for i := range bot.ValidNodes {
		strs[i] = "`" + bot.ValidNodes[i] + "`\n"
	}

	_, err = bot.PagedEmbed(ctx, bcr.StringPaginator("Valid permission nodes", bot.EmbedColour, strs, 15), 15*time.Minute)
	return err
}

func (bot *Bot) setNode(ctx *bcr.Context) (err error) {
	node := strings.ToLower(ctx.Args[0])
	if !contains(bot.ValidNodes, node) {
		return ctx.SendfX("%v is not a valid permission node.", bcr.AsCode(node))
	}

	level := common.DisabledLevel
	switch strings.ToLower(ctx.Args[1]) {
	case "disabled", "0":
		level = common.DisabledLevel
	case "everyone", "1":
		level = common.EveryoneLevel
	case "moderator", "2":
		level = common.ModeratorLevel
	case "manager", "3":
		level = common.ManagerLevel
	case "admin", "4":
		level = common.AdminLevel
	default:
		return ctx.SendfX("%v is not a valid permission level.", bcr.AsCode(ctx.Args[1]))
	}

	err = bot.DB.SetPermissions(ctx.Guild.ID, node, level)
	if err != nil {
		return bot.Report(ctx, err)
	}

	return ctx.SendX("", discord.Embed{
		Title:       "Permission node updated",
		Description: fmt.Sprintf("Commands under %v can now be used by %v and higher.", bcr.AsCode(node), bcr.AsCode(level.String())),
		Color:       bot.EmbedColour,
	})
}

func (bot *Bot) resetNode(ctx *bcr.Context) (err error) {
	node := strings.ToLower(ctx.RawArgs)

	var exists bool
	err = bot.DB.Pool.QueryRow(context.Background(),
		"select exists(select * from permission_nodes where guild_id = $1 and name = $2)",
		ctx.Guild.ID, node,
	).Scan(&exists)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if !exists {
		return ctx.SendfX("No node named %v exists, or it already doesn't have an override.", bcr.AsCode(node))
	}

	err = bot.DB.ResetPermissions(ctx.Guild.ID, node)
	if err != nil {
		return bot.Report(ctx, err)
	}

	return ctx.SendX("", discord.Embed{
		Title:       "Permission node updated",
		Description: fmt.Sprintf("%v's permissions have been reset to the default value.", bcr.AsCode(node)),
		Color:       bot.EmbedColour,
	})
}

func contains(slice []string, s string) bool {
	for _, i := range slice {
		if i == s {
			return true
		}
	}
	return false
}
