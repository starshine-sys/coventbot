package keyroles

import (
	"context"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) channel(ctx *bcr.Context) (err error) {
	var id discord.ChannelID

	if ctx.RawArgs == "-clear" || ctx.RawArgs == "--clear" {
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

	var logChannel discord.ChannelID
	err = bot.DB.Pool.QueryRow(context.Background(), "select keyrole_channel from servers where id = $1", ctx.Guild.ID).Scan(&logChannel)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if logChannel == id {
		_, err = ctx.Send("The given channel is already the key role channel.", nil)
		return err
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set keyrole_channel = $1 where id = $2", id, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if id == 0 {
		_, err = ctx.Send("Key role channel reset.", nil)
		return
	}
	_, err = ctx.Sendf("Key role channel changed to %v.", id.Mention())
	return
}
