package config

import (
	"context"
	"strings"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) prefix(ctx *bcr.Context) (err error) {
	prefixes, err := bot.DB.Prefixes(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(prefixes) == 0 {
		prefixes = bot.Router.Prefixes
	}

	prefixes = append([]string{ctx.Bot.Mention()}, prefixes...)

	_, err = ctx.SendEmbed(bcr.SED{
		Title:   "Prefixes",
		Message: strings.Join(prefixes, "\n"),
	})
	return err
}

func (bot *Bot) prefixAdd(ctx *bcr.Context) (err error) {
	prefixes, err := bot.DB.Prefixes(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(prefixes) > 20 {
		_, err = bot.Send(ctx, "This server already has the maximum number of prefixes (20).")
		return
	}

	if strings.Contains(ctx.RawArgs, "@") {
		_, err = bot.Send(ctx, "Prefix can't include a mention.")
		return
	}

	for _, p := range prefixes {
		if strings.EqualFold(p, ctx.RawArgs) {
			_, err = bot.Send(ctx, "``%v`` is already a prefix for this server.", bcr.EscapeBackticks(ctx.RawArgs))
			return
		}
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set prefixes = array_append(prefixes, $1) where id = $2", ctx.RawArgs, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = bot.Send(ctx, "Added prefix ``%v``", bcr.EscapeBackticks(ctx.RawArgs))
	return
}

func (bot *Bot) prefixRemove(ctx *bcr.Context) (err error) {
	prefixes, err := bot.DB.Prefixes(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if strings.Contains(ctx.RawArgs, "@") {
		_, err = bot.Send(ctx, "You can't remove the <@!%v> prefix.", ctx.Bot.ID)
		return
	}

	var isPrefix bool
	for _, p := range prefixes {
		if strings.EqualFold(p, ctx.RawArgs) {
			isPrefix = true
			break
		}
	}

	if !isPrefix {
		_, err = bot.Send(ctx, "``%v`` is not a prefix for this server.", bcr.EscapeBackticks(ctx.RawArgs))
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set prefixes = array_remove(prefixes, $1) where id = $2", ctx.RawArgs, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = bot.Send(ctx, "Removed prefix ``%v``", bcr.EscapeBackticks(ctx.RawArgs))
	return
}
