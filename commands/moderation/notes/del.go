package notes

import (
	"strconv"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) delNote(ctx *bcr.Context) (err error) {
	if !bot.globalPerms(ctx).Has(discord.PermissionManageRoles) {
		_, err = ctx.Replyc(bcr.ColourRed, "You're not allowed to use this command.")
		return
	}

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
