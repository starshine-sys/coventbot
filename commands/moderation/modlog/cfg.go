package modlog

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *ModLog) setchannel(ctx *bcr.Context) (err error) {
	id := discord.ChannelID(0)

	if ctx.RawArgs != "-clear" {
		ch, err := ctx.ParseChannel(ctx.RawArgs)
		if err != nil {
			_, err = ctx.Send("Channel not found.")
			return err
		}

		if ctx.Message.GuildID != ch.GuildID {
			_, err = ctx.Send("The given channel isn't in this server.")
			return err
		}

		id = ch.ID
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set mod_log_channel = $1 where id = $2", id, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if id == 0 {
		_, err = ctx.Send("Mod log channel cleared.")
	} else {
		_, err = ctx.Sendf("Mod log channel set to %v.", id.Mention())
	}
	return
}
