package roles

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) config(ctx *bcr.Context) (err error) {
	name := ctx.Args[0]
	var colour discord.Color
	var requireRole discord.RoleID
	var existing bool

	{
		cat, err := bot.categoryName(ctx.Guild.ID, ctx.Args[0])
		if err == nil {
			name = cat.Name
			existing = true
		}
	}

	roleStr, _ := ctx.Flags.GetString("require-role")
	desc, _ := ctx.Flags.GetString("desc")
	clr, _ := ctx.Flags.GetString("colour")

	if clr != "" {
		i, err := strconv.ParseUint(strings.TrimPrefix(clr, "#"), 16, 0)
		if err != nil {
			_, err = ctx.Replyc(bcr.ColourRed, "Couldn't parse your input ``%v`` as a colour.", clr)
			return err
		}
		colour = discord.Color(i)
	}

	if roleStr != "" {
		r, err := ctx.ParseRole(roleStr)
		if err != nil {
			_, err = ctx.Replyc(bcr.ColourRed, "Couldn't parse your input ``%v`` as a role.", roleStr)
			return err
		}
		requireRole = r.ID
	}

	if len(desc) > 1000 {
		_, err = ctx.Replyc(bcr.ColourRed, "Description too long (%v > 1000 characters).", len(desc))
		return
	}

	c, err := bot.newCategory(ctx.Guild.ID, name, desc, requireRole, colour)
	if err != nil {
		return bot.Report(ctx, err)
	}

	e := discord.Embed{
		Title:       "Category " + c.Name + " created",
		Description: fmt.Sprintf("**Colour:** #%06X\n", colour),
		Color:       colour,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v", c.ID),
		},
	}

	if e.Color == 0 {
		e.Color = bcr.ColourBlurple
	}
	if requireRole.IsValid() {
		e.Description += "**Require role:** " + requireRole.Mention()
	}
	if desc != "" {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Description",
			Value: desc,
		})
	}
	if existing {
		e.Title = "Category " + c.Name + " updated"
	}

	_, err = ctx.Send("", &e)
	return
}

func (bot *Bot) delete(ctx *bcr.Context) (err error) {
	id, err := strconv.ParseInt(ctx.RawArgs, 0, 0)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "Couldn't parse your input as an ID.")
		return
	}

	ct, err := bot.DB.Pool.Exec(context.Background(), "delete from role_categories where id = $1 and server_id = $2", id, ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}
	if ct.RowsAffected() == 0 {
		_, err = ctx.Replyc(bcr.ColourRed, "Couldn't find a category with that ID.")
	}

	_, err = ctx.Reply("Category deleted.")
	return
}
