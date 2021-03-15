package config

import (
	"context"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) watchlistChannel(ctx *bcr.Context) (err error) {
	var id discord.ChannelID

	if ctx.RawArgs == "-clear" {
		id = 0
	} else {
		ch, err := ctx.ParseChannel(ctx.RawArgs)
		if err != nil {
			_, err = ctx.Send("Channel not found.", nil)
			return err
		}

		if ch.GuildID != ctx.Message.GuildID {
			_, err = ctx.Send("The given channel isn't in this server.", nil)
			return err
		}

		id = ch.ID
	}

	current, err := bot.DB.WatchlistChannel(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	if current == id {
		_, err = ctx.Send("The given channel is already the watch list channel.", nil)
		return err
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set watch_list_channel = $1 where id = $2", id, ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	if id == 0 {
		_, err = ctx.Send("Watchlist channel reset.", nil)
		return
	}
	_, err = ctx.Sendf("Watchlist channel changed to %v.", id.Mention())
	return
}
