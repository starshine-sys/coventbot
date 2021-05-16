package moderation

import (
	"time"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) cmdSetSlowmode(ctx *bcr.Context) (err error) {
	ch, err := ctx.ParseChannel(ctx.Args[0])
	if err != nil {
		_, err = ctx.Sendf("Invalid channel given: ``%v``", bcr.EscapeBackticks(ctx.Args[0]))
		return
	}

	if ch.GuildID != ctx.Channel.GuildID {
		_, err = ctx.Sendf("The channel must be in this server.")
		return
	}

	clear, _ := ctx.Flags.GetBool("clear")
	if clear {
		err = bot.clearSlowmode(ch.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		_, err = ctx.Sendf("Cleared the slowmode for %v.", ch.Mention())
		return
	}

	if len(ctx.Args) < 2 {
		_, err = ctx.Sendf("You must give a duration.")
		return
	}

	duration, err := time.ParseDuration(ctx.Args[1])
	if err != nil {
		_, err = ctx.Sendf("Invalid duration given: ``%v``", bcr.EscapeBackticks(ctx.Args[1]))
		return
	}

	err = bot.setSlowmode(ch.GuildID, ch.ID, duration)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Sendf("Set the slowmode for %v to %s.", ch.Mention(), duration)
	return
}
