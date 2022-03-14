package quotes

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) cmdQuoteDelete(ctx *bcr.Context) (err error) {
	q, err := bot.getQuote(ctx.RawArgs, ctx.Guild.ID)
	if err != nil {
		_, err = ctx.Send("No quote with that ID found.")
		return
	}

	yes, timeout := ctx.ConfirmButton(ctx.Author.ID, bcr.ConfirmData{
		Embeds: []discord.Embed{{
			Color:       bcr.ColourBlurple,
			Description: fmt.Sprintf("Are you sure you want to delete the quote `%v` by %v?", q.HID, q.UserID.Mention()),
		}},
		YesPrompt: "Delete",
		YesStyle:  discord.DangerButtonStyle(),
	})
	if !yes || timeout {
		_, err = ctx.Send(":x: Cancelled.")
		return
	}

	err = bot.delQuote(ctx.Guild.ID, q.HID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Sendf("Quote `%v` deleted!", q.HID)
	return
}
