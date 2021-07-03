package moderation

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) slowmodeRole(ctx *bcr.Context) (err error) {
	if len(ctx.Args) == 0 {
		var id discord.RoleID
		err = bot.DB.Pool.QueryRow(context.Background(), "select slowmode_ignore_role from servers where id = $1", ctx.Message.GuildID).Scan(&id)
		if err != nil {
			return bot.Report(ctx, err)
		}

		s := "No role currently ignores slowmode."
		if r, err := ctx.State.Role(ctx.Message.GuildID, id); err == nil {
			s = fmt.Sprintf("The currently ignored role is %v. To reset this role, call this command with `-clear`.", r.Mention())
		}

		_, err = ctx.Send("", discord.Embed{
			Description: s,
			Color:       ctx.Router.EmbedColor,
		})
		return
	}

	var id discord.RoleID
	if ctx.RawArgs == "-clear" || ctx.RawArgs == "--clear" {
		id = 0
	} else {
		r, err := ctx.ParseRole(ctx.RawArgs)
		if err != nil {
			_, err = ctx.Send("I could not find that role.")
			return err
		}
		id = r.ID
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set slowmode_ignore_role = $1 where id = $2", id, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	desc := fmt.Sprintf("%v will now ignore slowmode.", id.Mention())
	if !id.IsValid() {
		desc = fmt.Sprintf("Slowmode ignore role reset.")
	}

	_, err = ctx.Send("", discord.Embed{
		Description: desc,
		Color:       ctx.Router.EmbedColor,
	})
	return
}
