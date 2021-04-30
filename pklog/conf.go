package pklog

import (
	"context"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) setChannel(ctx *bcr.Context) (err error) {
	id := discord.ChannelID(0)

	if ctx.RawArgs != "-clear" {
		ch, err := ctx.ParseChannel(ctx.RawArgs)
		if err != nil {
			_, err = ctx.Send("Channel not found.", nil)
			return err
		}

		if ctx.Message.GuildID != ch.GuildID {
			_, err = ctx.Send("The given channel isn't in this server.", nil)
			return err
		}

		id = ch.ID
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set pk_log_channel = $1 where id = $2", id, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if id == 0 {
		_, err = ctx.Send("PluralKit message logging disabled.", nil)
	} else {
		_, err = ctx.Sendf("PluralKit messages are now being logged to %v.", id.Mention())
	}
	return
}

func (bot *Bot) resetCache(ctx *bcr.Context) (err error) {
	bot.ResetCache(ctx.Message.GuildID)
	_, err = ctx.Send("Webhook cache reset.", nil)
	return
}
