package tags

import (
	"context"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) isModerator(ctx *bcr.Context) (isMod bool) {
	var id discord.RoleID
	bot.DB.Pool.QueryRow(context.Background(), "select tag_mod_role from servers where id = $1", ctx.Message.GuildID).Scan(&id)
	if id.IsValid() {
		if ctx.Member == nil {
			return false
		}

		for _, r := range ctx.Member.RoleIDs {
			if r == id {
				return true
			}
		}

		return false
	}

	p, _ := ctx.State.Permissions(ctx.Channel.ID, ctx.Author.ID)
	return p.Has(discord.PermissionManageGuild)
}

func (bot *Bot) editTag(ctx *bcr.Context) (err error) {
	args := strings.Split(ctx.RawArgs, "\n")
	if len(args) < 2 {
		_, err = ctx.Send("Not enough arguments given: need at least 2, separated by a newline.", nil)
		return
	}

	t, err := bot.DB.GetTag(ctx.Message.GuildID, args[0])
	if err != nil {
		_, err = ctx.Send("No tag with that name found.", nil)
		return
	}

	if t.CreatedBy != ctx.Author.ID && !bot.isModerator(ctx) {
		_, err = ctx.Send("You don't have permission to edit this tag.", nil)
		return
	}

	_, err = bot.DB.Pool.Exec(
		context.Background(), "update tags set response = $1 where id = $2",
		strings.TrimSpace(strings.Join(args[1:], "\n")),
		t.ID,
	)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("Tag updated!", nil)
	return
}
