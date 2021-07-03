package quotes

import (
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) cmdQuoteDelete(ctx *bcr.Context) (err error) {
	q, err := bot.getQuote(ctx.RawArgs, ctx.Guild.ID)
	if err != nil {
		_, err = ctx.Send("No quote with that ID found.")
		return
	}

	e := q.Embed(bot.PK)
	msg, err := ctx.Send("Are you sure you want to delete this quote?", e)
	if err != nil {
		return
	}

	yes, timeout := ctx.YesNoHandler(*msg, ctx.Author.ID)
	if !yes || timeout {
		_, err = ctx.Send(":x: Cancelled.")
		return
	}

	err = bot.delQuote(ctx.Guild.ID, q.HID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	err = ctx.State.DeleteMessage(ctx.Message.ChannelID, msg.ID)
	if err != nil {
		bot.Sugar.Errorf("Error deleting quote message: %v", err)
	}

	_, err = ctx.Sendf("Quote `%v` deleted!", q.HID)
	return
}
