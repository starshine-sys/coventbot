// SPDX-License-Identifier: AGPL-3.0-only
package keyroles

import (
	"context"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) list(ctx *bcr.Context) (err error) {
	var keyRoles []uint64
	err = bot.DB.Pool.QueryRow(context.Background(), "select keyroles from servers where id = $1", ctx.Guild.ID).Scan(&keyRoles)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(keyRoles) == 0 {
		_, err = ctx.Reply("There are no key roles set.")
		return
	}

	s := []string{}

	for _, r := range keyRoles {
		s = append(s, fmt.Sprintf("<@&%v>\n", r))
	}

	_, err = bot.PagedEmbed(ctx,
		bcr.StringPaginator("Key roles", bcr.ColourBlurple, s, 20), 10*time.Minute,
	)
	return err
}

func (bot *Bot) add(ctx *bcr.Context) (err error) {
	var keyRoles []uint64
	err = bot.DB.Pool.QueryRow(context.Background(), "select keyroles from servers where id = $1", ctx.Guild.ID).Scan(&keyRoles)
	if err != nil {
		return bot.Report(ctx, err)
	}

	r, err := ctx.ParseRole(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Reply("Role not found.")
		return
	}

	for _, kr := range keyRoles {
		if discord.RoleID(kr) == r.ID {
			_, err = ctx.Reply("%v is already a key role.", r.Mention())
			return
		}
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set keyroles = array_append(keyroles, $1) where id = $2", r.ID, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Added key role %v.", r.Mention())
	return
}

func (bot *Bot) remove(ctx *bcr.Context) (err error) {
	var keyRoles []uint64
	err = bot.DB.Pool.QueryRow(context.Background(), "select keyroles from servers where id = $1", ctx.Guild.ID).Scan(&keyRoles)
	if err != nil {
		return bot.Report(ctx, err)
	}

	r, err := ctx.ParseRole(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Reply("Role not found.")
		return
	}

	var isKeyRole bool
	for _, kr := range keyRoles {
		if discord.RoleID(kr) == r.ID {
			isKeyRole = true
			break
		}
	}

	if !isKeyRole {
		_, err = ctx.Reply("%v is not a key role.", r.Mention())
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set keyroles = array_remove(keyroles, $1) where id = $2", r.ID, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Removed key role %v.", r.Mention())
	return
}
