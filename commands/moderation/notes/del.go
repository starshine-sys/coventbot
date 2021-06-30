package notes

import (
	"strconv"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) delNote(ctx *bcr.Context) (err error) {
	id, err := strconv.ParseUint(ctx.RawArgs, 0, 0)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "Couldn't parse your input (``%v``) as a number.", ctx.RawArgs)
		return
	}

	err = bot.DB.DelNote(ctx.Guild.ID, id)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "No note with that ID found.")
		return
	}

	_, err = ctx.Reply("Note #%v deleted!", id)
	return
}
