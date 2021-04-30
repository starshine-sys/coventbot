package roles

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) addCategory(ctx *bcr.Context) (err error) {
	name := ctx.RawArgs

	var exists bool
	err = bot.DB.Pool.QueryRow(context.Background(), "select exists (select * from roles where server_id = $1 and name ilike $2)", ctx.Message.GuildID, name).Scan(&exists)
	if err != nil {
		bot.Sugar.Errorf("Error: %v", err)
		return bot.Report(ctx, err)
	}
	if exists {
		_, err = ctx.Send("A category with that name already exists.", nil)
		return
	}

	c := Category{
		Name:     name,
		ServerID: ctx.Message.GuildID,
	}

	err = bot.DB.Pool.QueryRow(context.Background(), "insert into roles (server_id, name) values ($1, $2) returning id", c.ServerID, c.Name).Scan(&c.ID)
	if err != nil {
		bot.Sugar.Errorf("Error: %v", err)
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("", &discord.Embed{
		Title:       "Category created",
		Description: fmt.Sprintf("Role category \"%v\" created with ID %v.", c.Name, c.ID),
		Color:       bcr.ColourGreen,
	})
	return
}
