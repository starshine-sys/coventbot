package tags

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) tagModerator(ctx *bcr.Context) (err error) {
	if len(ctx.Args) == 0 {
		var id discord.RoleID
		err = bot.DB.Pool.QueryRow(context.Background(), "select tag_mod_role from servers where id = $1", ctx.Message.GuildID).Scan(&id)
		if err != nil {
			_, err = ctx.Sendf("Internal error occurred: %v", err)
			return
		}

		s := "There is currently no tag mod role set. Anyone can create, edit, and delete tags."
		if r, err := ctx.State.Role(ctx.Message.GuildID, id); err == nil {
			s = fmt.Sprintf("The current tag mod role is %v. To reset this role, call this command with `-clear`.", r.Mention())
		}

		_, err = ctx.Send("", &discord.Embed{
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
			_, err = ctx.Send("I could not find that role.", nil)
			return err
		}
		id = r.ID
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set tag_mod_role = $1 where id = $2", id, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	desc := fmt.Sprintf("Tag mod role set to %v. Note that users without this role will not be able to create, edit, or delete tags.", id.Mention())
	if !id.IsValid() {
		desc = fmt.Sprintf("Tag mod role reset. Note that anyone will be able to create, edit, or delete tags.")
	}

	_, err = ctx.Send("", &discord.Embed{
		Description: desc,
		Color:       ctx.Router.EmbedColor,
	})
	return
}
