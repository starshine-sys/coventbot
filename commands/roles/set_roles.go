package roles

import (
	"strconv"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) setRoles(ctx *bcr.Context) (err error) {
	id, err := strconv.ParseInt(ctx.Args[0], 0, 0)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "Couldn't parse your input as an ID.")
		return
	}

	cat, err := bot.categoryID(ctx.Guild.ID, id)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "No category with that ID found.")
	}

	roles, _ := ctx.GreedyRoleParser(ctx.Args[1:])
	var ids []uint64
	for _, r := range roles {
		ids = append(ids, uint64(r.ID))
	}

	err = bot.categoryRoles(cat.ID, ids)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Roles in %v updated!", cat.Name)
	return
}
