// SPDX-License-Identifier: AGPL-3.0-only
package levels

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) blacklistAdd(ctx *bcr.Context) (err error) {
	// try parse channel
	ch, err := ctx.ParseChannel(ctx.RawArgs)
	if err == nil {
		if ch.GuildID != ctx.Guild.ID || (ch.Type != discord.GuildText && ch.Type != discord.GuildNews && ch.Type != discord.GuildCategory) {
			return ctx.SendX("Invalid channel provided, must be in this guild and be a text or category channel.")
		}

		if ch.Type == discord.GuildCategory {
			_, err := bot.DB.Pool.Exec(context.Background(), "update server_levels set blocked_categories = array_append(blocked_categories, $1) where id = $2", ch.ID, ctx.Guild.ID)
			if err != nil {
				return bot.Report(ctx, err)
			}

			return ctx.SendfX("Added %v to the list of blocked categories!", ch.Name)
		} else {
			_, err := bot.DB.Pool.Exec(context.Background(), "update server_levels set blocked_channels = array_append(blocked_channels, $1) where id = $2", ch.ID, ctx.Guild.ID)
			if err != nil {
				return bot.Report(ctx, err)
			}

			return ctx.SendfX("Added %v/#%v to the list of blocked channels!", ch.Mention(), ch.Name)
		}
	}

	// else try parsing role
	r, err := ctx.ParseRole(ctx.RawArgs)
	if err != nil {
		return ctx.SendX("Input is not a valid role or channel.")
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update server_levels set blocked_roles = array_append(blocked_roles, $1) where id = $2", r.ID, ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	return ctx.SendfX("Added @%v to the list of blocked roles!", r.Name)
}

func (bot *Bot) blacklistRemove(ctx *bcr.Context) (err error) {
	// try parse channel
	ch, err := ctx.ParseChannel(ctx.RawArgs)
	if err == nil {
		if ch.GuildID != ctx.Guild.ID || (ch.Type != discord.GuildText && ch.Type != discord.GuildNews && ch.Type != discord.GuildCategory) {
			return ctx.SendX("Invalid channel provided, must be in this guild and be a text or category channel.")
		}

		if ch.Type == discord.GuildCategory {
			_, err := bot.DB.Pool.Exec(context.Background(), "update server_levels set blocked_categories = array_remove(blocked_categories, $1) where id = $2", ch.ID, ctx.Guild.ID)
			if err != nil {
				return bot.Report(ctx, err)
			}

			return ctx.SendfX("Removed %v from the list of blocked categories!", ch.Name)
		} else {
			_, err := bot.DB.Pool.Exec(context.Background(), "update server_levels set blocked_channels = array_remove(blocked_channels, $1) where id = $2", ch.ID, ctx.Guild.ID)
			if err != nil {
				return bot.Report(ctx, err)
			}

			return ctx.SendfX("Removed %v/#%v from the list of blocked channels!", ch.Mention(), ch.Name)
		}
	}

	// else try parsing role
	r, err := ctx.ParseRole(ctx.RawArgs)
	if err != nil {
		return ctx.SendX("Input is not a valid role or channel.")
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update server_levels set blocked_roles = array_remove(blocked_roles, $1) where id = $2", r.ID, ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	return ctx.SendfX("Removed @%v from the list of blocked roles!", r.Name)
}
