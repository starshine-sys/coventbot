// SPDX-License-Identifier: AGPL-3.0-only
package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/db"
)

func (bot *Bot) watchlist(ctx *bcr.Context) (err error) {
	b, err := bot.DB.Watchlist(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	var x string
	// append all user IDs (as mentions) to x
	for _, c := range b {
		x += fmt.Sprintf("<@%v> (%v)\n", c, c)
	}
	if len(b) == 0 {
		x = "No users are on the watchlist."
	}
	_, err = ctx.Send("", discord.Embed{
		Title:       "Watchlist",
		Description: x,
		Color:       ctx.Router.EmbedColor,
	})
	return err
}

func (bot *Bot) watchlistRemove(ctx *bcr.Context) (err error) {
	if ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("You need to provide a channel.")
		return err
	}

	u, err := ctx.ParseUser(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("The user you gave was not found.")
		return err
	}

	err = bot.DB.RemoveFromWatchlist(ctx.Message.GuildID, u.ID)
	if err != nil {
		if err == db.ErrorNotBlacklisted {
			_, err = ctx.Send("That user isn't on the watchlist.")
			return err
		}

		return bot.Report(ctx, err)
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "delete from watch_list_reasons where user_id = $1 and server_id = $2", u.ID, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Sendf("Removed %v / **%v#%v** from the watchlist.", u.Mention(), u.Username, u.Discriminator)
	return
}

func (bot *Bot) watchlistAdd(ctx *bcr.Context) (err error) {
	if ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("You need to provide a user.")
		return err
	}

	u, err := ctx.ParseUser(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("The user you gave was not found.")
		return err
	}

	err = bot.DB.AddToWatchlist(ctx.Message.GuildID, u.ID)
	if err != nil {
		if err == db.ErrorAlreadyBlacklisted {
			_, err = ctx.Send("That user is already on the watchlist.")
			return err
		}

		return bot.Report(ctx, err)
	}

	_, err = ctx.Sendf("Added %v / **%v#%v** to the watchlist.", u.Mention(), u.Username, u.Discriminator)
	return
}

func (bot *Bot) watchlistReason(ctx *bcr.Context) (err error) {
	u, err := ctx.ParseUser(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("I could not find that user.")
		return
	}

	if len(ctx.Args) == 1 {
		var reason string
		bot.DB.Pool.QueryRow(context.Background(), "select reason from watch_list_reasons where user_id = $1 and server_id = $2", u.ID, ctx.Message.GuildID).Scan(&reason)
		if reason == "" {
			_, err = ctx.Send("There is no reason set for that user.")
			return
		}

		_, err = ctx.Sendf("Reason for %v#%v:\n> %v", u.Username, u.Discriminator, reason)
		return
	}

	// if the user isn't on the watch list, return
	if !bot.DB.IsWatchlisted(ctx.Message.GuildID, u.ID) {
		_, err = ctx.Sendf("That user isn't on the watchlist.")
		return
	}

	reason := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))

	_, err = bot.DB.Pool.Exec(context.Background(), "insert into watch_list_reasons (user_id, server_id, reason) values ($1, $2, $3) on conflict (user_id, server_id) do update set reason = $3", u.ID, ctx.Message.GuildID, reason)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Sendf("Updated watchlist reason for %v#%v.", u.Username, u.Discriminator)
	return
}
