package config

import (
	"strings"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) prefix(ctx *bcr.Context) (err error) {
	if len(ctx.Args) == 0 {
		prefixes, err := bot.DB.Prefixes(ctx.Message.GuildID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		if len(prefixes) == 0 {
			_, err = ctx.Send("This server has no custom prefixes.", nil)
			return err
		}

		_, err = ctx.Sendf("This server's prefixes: ``%v``", bcr.EscapeBackticks(strings.Join(prefixes, ", ")))
		return err
	}

	if len(ctx.Args) > 10 {
		_, err = ctx.Send("Too many prefixes, maximum of 10.", nil)
		return
	}

	if ctx.RawArgs == "-clear" {
		ctx.Args = []string{}
	}

	err = bot.DB.SetPrefixes(ctx.Message.GuildID, ctx.Args)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(ctx.Args) > 0 {
		_, err = ctx.Sendf("Prefixes updated! New prefixes: ``%v``", bcr.EscapeBackticks(strings.Join(ctx.Args, ", ")))
		return
	}

	_, err = ctx.Send("Prefixes reset.", nil)
	return err
}
