package todos

import (
	"context"
	"strconv"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) delete(ctx *bcr.Context) (err error) {
	id, err := strconv.ParseInt(ctx.RawArgs, 0, 0)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "You didn't give a valid ID.")
	}

	t, err := bot.getTodo(id, ctx.Author.ID)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "Couldn't find a todo with that ID.")
		return
	}

	ctx.State.DeleteMessage(t.ChannelID, t.MID)

	ct, err := bot.DB.Pool.Exec(context.Background(), "delete from todos where id = $1 and user_id = $2", id, ctx.Author.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if ct.RowsAffected() == 0 {
		_, err = ctx.Replyc(bcr.ColourRed, "Couldn't find a todo with that ID.")
	} else {
		_, err = ctx.Reply("Todo #%v deleted!", id)
	}
	return
}
