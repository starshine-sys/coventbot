package info

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/common"
)

func (bot *Bot) roleInfo(ctx *bcr.Context) (err error) {
	r, err := ctx.ParseRole(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("The specified role could not be found.")
		return
	}

	colour := r.Color
	if colour == 0 {
		colour = ctx.Router.EmbedColor
	}

	e := discord.Embed{
		Title:       "Role info",
		Description: "`" + r.Mention() + "`",
		Color:       colour,
		Fields: []discord.EmbedField{
			{
				Name:   "ID",
				Value:  r.ID.String(),
				Inline: true,
			},
			{
				Name:   "Name",
				Value:  r.Name,
				Inline: true,
			},
			{
				Name:   "Colour",
				Value:  fmt.Sprintf("#%v", r.Color.String()),
				Inline: true,
			},
			{
				Name:   "Position",
				Value:  fmt.Sprint(r.Position),
				Inline: true,
			},
			{
				Name:   "Mentionable",
				Value:  fmt.Sprint(r.Mentionable),
				Inline: true,
			},
			{
				Name:   "Hoisted",
				Value:  fmt.Sprint(r.Hoist),
				Inline: true,
			},
			{
				Name:   "Created",
				Value:  fmt.Sprintf("<t:%v>\n(%v)", r.ID.Time().Unix(), common.FormatTime(r.ID.Time().UTC())),
				Inline: false,
			},
			{
				Name:   "Permissions",
				Value:  fmt.Sprintf("%v", FirstNonZero(strings.Join(bcr.PermStrings(r.Permissions), ", "), "None")),
				Inline: false,
			},
		},
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v | Created at", r.ID),
		},
		Timestamp: discord.Timestamp(r.ID.Time()),
	}

	_, err = ctx.Send("", e)
	return
}
