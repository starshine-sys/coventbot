// SPDX-License-Identifier: AGPL-3.0-only
package tickets

import (
	"context"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) cfg(ctx *bcr.Context) (err error) {
	cat, err := ctx.ParseChannel(ctx.Args[0])
	if err != nil || cat.GuildID != ctx.Message.GuildID || cat.Type != discord.GuildCategory {
		_, err = ctx.Sendf("Category not found.")
		return
	}

	name := ctx.Args[1]

	logChannel, err := ctx.ParseChannel(ctx.Args[2])
	if err != nil {
		_, err = ctx.Send("Given log channel not found.")
		return
	}

	limit, _ := ctx.Flags.GetInt("limit")
	count, _ := ctx.Flags.GetUint("count")
	creatorClose, _ := ctx.Flags.GetBool("creator-close")

	if count != 0 {
		count = count - 1
	}

	sql := `insert into ticket_categories
(category_id, server_id, per_user_limit, log_channel, count, can_creator_close, name)
values ($1, $2, $3, $4, $5, $6, $7)
on conflict (category_id) do update
set per_user_limit = $3, log_channel = $4, count = $5, can_creator_close = $6, name = $7`

	_, err = bot.DB.Pool.Exec(context.Background(), sql, cat.ID, ctx.Message.GuildID, limit, logChannel.ID, count, creatorClose, name)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Updated the configuration for %v: %v", cat.Name, name)
	return
}

func (bot *Bot) mention(ctx *bcr.Context) (err error) {
	cat, err := ctx.ParseChannel(ctx.Args[0])
	if err != nil || cat.GuildID != ctx.Message.GuildID || cat.Type != discord.GuildCategory {
		_, err = ctx.Sendf("Category not found.")
		return
	}

	mention := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))

	if mention == "-clear" || mention == "--clear" {
		mention = ""
	}

	var exists bool
	err = bot.DB.Pool.QueryRow(context.Background(), "select exists(select * from ticket_categories where category_id = $1)", cat.ID).Scan(&exists)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if !exists {
		_, err = ctx.Sendf("The channel you gave (%v) isn't a ticket category.", cat.Name)
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update ticket_categories set mention = $1", mention)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Updated mention to ``%v``.", bcr.EscapeBackticks(mention))
	return err
}

func (bot *Bot) description(ctx *bcr.Context) (err error) {
	cat, err := ctx.ParseChannel(ctx.Args[0])
	if err != nil || cat.GuildID != ctx.Message.GuildID || cat.Type != discord.GuildCategory {
		_, err = ctx.Sendf("Category not found.")
		return
	}

	description := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))

	if description == "-clear" || description == "--clear" {
		description = ""
	}

	var exists bool
	err = bot.DB.Pool.QueryRow(context.Background(), "select exists(select * from ticket_categories where category_id = $1)", cat.ID).Scan(&exists)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if !exists {
		_, err = ctx.Sendf("The channel you gave (%v) isn't a ticket category.", cat.Name)
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update ticket_categories set description = $1", description)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Updated description to\n```%v```", description)
	return err
}
