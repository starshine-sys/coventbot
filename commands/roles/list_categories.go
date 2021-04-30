package roles

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/etc"
)

// Category is a single category
type Category struct {
	ID       uint64
	ServerID discord.GuildID
	Name     string

	RequireRole discord.RoleID
	Roles       []uint64
}

func (bot *Bot) listCategories(ctx *bcr.Context) (err error) {
	var cats []*Category

	err = pgxscan.Select(context.Background(), bot.DB.Pool, &cats, "select id, name, server_id, require_role, roles from roles where server_id = $1 order by name", ctx.Message.GuildID)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			_, err = ctx.Send("This server has no role categories.", nil)
			return
		}
	}

	e := discord.Embed{
		Title: "Roles",
		Color: etc.ColourBlurple,
	}

	for _, c := range cats {
		e.Description += fmt.Sprintf("%v: %v role", c.Name, len(c.Roles))
		if len(c.Roles) != 1 {
			e.Description += "s"
		}
		e.Description += "\n"
	}

	e.Description += fmt.Sprintf("\n\nUse `%vroles [category name]` to see the roles in a category.", ctx.Prefix)

	_, err = ctx.Send("", &e)
	return
}
