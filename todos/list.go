package todos

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) list(ctx *bcr.Context) (err error) {
	todos := []Todo{}

	err = pgxscan.Select(context.Background(), bot.DB.Pool, &todos, "select * from todos where user_id = $1 and complete = false order by id asc", ctx.Author.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(todos) == 0 {
		_, err = ctx.Reply("You have no incomplete todos. Set some with `%vtodo`!", ctx.Prefix)
		return
	}

	var slice []string

	for _, todo := range todos {
		text := todo.Description
		if len(text) > 200 {
			text = text[:200] + "..."
		}

		slice = append(slice, fmt.Sprintf(`**#%v**: %v
%v UTC ([link](https://discord.com/channels/%v/%v/%v), [original](https://discord.com/channels/%v/%v/%v))

`, todo.ID, text, todo.Created.Format("2006-01-02 | 15:04"), todo.ServerID, todo.ChannelID, todo.MID, todo.OrigServerID, todo.OrigChannelID, todo.OrigMID))
	}

	_, err = ctx.PagedEmbed(bcr.StringPaginator(fmt.Sprintf("Todos (%v)", len(todos)), bcr.ColourBlurple, slice, 5), false)
	return
}
