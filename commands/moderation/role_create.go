package moderation

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) roleCreate(ctx *bcr.Context) (err error) {
	if p, _ := ctx.State.Permissions(ctx.Channel.ID, ctx.Bot.ID); !p.Has(discord.PermissionManageRoles) {
		return ctx.SendX("I do not have the Manage Roles permission in this server, and cannot create any roles.")
	}

	fs := flag.NewFlagSet("", flag.ContinueOnError)

	var hoist, mentionable bool
	fs.BoolVar(&hoist, "h", false, "Whether to hoist the new role")
	fs.BoolVar(&mentionable, "m", false, "Whether to make the new role mentionable")

	_ = fs.Parse(ctx.Args)
	ctx.Args = fs.Args()

	var color = discord.Color(-1)

	if len(ctx.Args) == 0 {
		_, err = ctx.Sendf("No name specified!")
		return
	}
	name := ctx.Args[0]

	if len(ctx.Args) > 1 {
		c, err := strconv.ParseUint(ctx.Args[1], 16, 32)
		if err != nil {
			_, err = ctx.Sendf("Couldn't parse your input (``%v``) as a colour.", bcr.EscapeBackticks(ctx.Args[1]))
			return err
		}

		color = discord.Color(c)
	}

	r, err := ctx.State.CreateRole(ctx.Message.GuildID, api.CreateRoleData{
		Name:        name,
		Permissions: 0,
		Color:       color,
		Hoist:       hoist,
		Mentionable: mentionable,
	})
	if err != nil {
		_, err = ctx.Send("I could not create a new role. Is this server at the 250 role limit?")
		return err
	}

	embedColor := r.Color
	if r.Color == 0 {
		embedColor = ctx.Router.EmbedColor
	}

	e := discord.Embed{
		Title: "Success!",
		Description: fmt.Sprintf(`The role **%v** has been created.
**Colour:** %s
**Mentionable:** %v
**Shown separately:** %v`, r.Name, r.Color, r.Mentionable, r.Hoist),
		Color: embedColor,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v", r.ID),
		},
	}

	_, err = ctx.Send("", e)
	return
}
