package config

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) adminRoles(ctx *bcr.Context) (err error) {
	var roles []uint64
	err = bot.DB.Pool.QueryRow(context.Background(), "select admin_roles from servers where id = $1", ctx.Guild.ID).Scan(&roles)
	if err != nil {
		return bot.Report(ctx, err)
	}

	s := ""
	for _, r := range roles {
		s += fmt.Sprintf("<@&%v>\n", r)
	}

	_, err = ctx.Send("", discord.Embed{
		Title:       "Admin roles for " + ctx.Guild.Name,
		Description: s,
		Color:       bcr.ColourBlurple,
	})
	return
}

func (bot *Bot) adminAddRole(ctx *bcr.Context) (err error) {
	r, err := ctx.ParseRole(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "I couldn't find that role.")
		return
	}

	var roles []uint64
	err = bot.DB.Pool.QueryRow(context.Background(), "select admin_roles from servers where id = $1", ctx.Guild.ID).Scan(&roles)
	if err != nil {
		return bot.Report(ctx, err)
	}

	for _, id := range roles {
		if r.ID == discord.RoleID(id) {
			_, err = ctx.Replyc(bcr.ColourRed, "%v is already a admin role.", r.Mention())
			return
		}
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set admin_roles = array_append(admin_roles, $1) where id = $2", r.ID, ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Added %v as a admin role!", r.Mention())
	return
}

func (bot *Bot) adminRemoveRole(ctx *bcr.Context) (err error) {
	r, err := ctx.ParseRole(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "I couldn't find that role.")
		return
	}

	var roles []uint64
	err = bot.DB.Pool.QueryRow(context.Background(), "select admin_roles from servers where id = $1", ctx.Guild.ID).Scan(&roles)
	if err != nil {
		return bot.Report(ctx, err)
	}

	var isSet bool
	for _, id := range roles {
		if r.ID == discord.RoleID(id) {
			isSet = true
		}
	}

	if !isSet {
		_, err = ctx.Replyc(bcr.ColourRed, "%v already isn't a admin role.", r.Mention())
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set admin_roles = array_remove(admin_roles, $1) where id = $2", r.ID, ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Removed %v as a admin role!", r.Mention())
	return
}
