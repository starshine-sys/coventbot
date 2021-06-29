package quotes

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) cmdQuoteDelete(ctx *bcr.Context) (err error) {
	if !bot.GlobalPerms(ctx).Has(discord.PermissionManageMessages) {
		_, err = ctx.Send("Only users with the Manage Messages permission can delete quotes.", nil)
		return
	}

	q, err := bot.getQuote(ctx.RawArgs, ctx.Guild.ID)
	if err != nil {
		_, err = ctx.Send("No quote with that ID found.", nil)
		return
	}

	e := q.Embed()
	msg, err := ctx.Send("Are you sure you want to delete this quote?", &e)
	if err != nil {
		return
	}

	yes, timeout := ctx.YesNoHandler(*msg, ctx.Author.ID)
	if !yes || timeout {
		_, err = ctx.Send(":x: Cancelled.", nil)
		return
	}

	err = bot.delQuote(ctx.Guild.ID, q.HID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	err = bot.State.DeleteMessage(ctx.Message.ChannelID, msg.ID)
	if err != nil {
		bot.Sugar.Errorf("Error deleting quote message: %v", err)
	}

	_, err = ctx.Sendf("Quote `%v` deleted!", q.HID)
	return
}
