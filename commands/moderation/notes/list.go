// SPDX-License-Identifier: AGPL-3.0-only
package notes

import (
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) list(ctx *bcr.Context) (err error) {
	u, err := ctx.ParseUser(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "User not found.")
		return
	}

	notes, err := bot.DB.UserNotes(ctx.Guild.ID, u.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(notes) == 0 {
		_, err = ctx.Reply("%v has no notes.", u.Mention())
		return
	}

	fields := []discord.EmbedField{}

	for _, n := range notes {
		fields = append(fields, discord.EmbedField{
			Name:  fmt.Sprintf("Note #%v", n.ID),
			Value: fmt.Sprintf("From %v at <t:%v>:\n%v", n.Moderator.Mention(), n.Created.Unix(), n.Note),
		})
	}

	embeds := bcr.FieldPaginator("", "", bcr.ColourBlurple, fields, 5)
	for i := range embeds {
		embeds[i].Author = &discord.EmbedAuthor{
			Name: fmt.Sprintf("%v (%v)", u.Tag(), u.ID),
			Icon: u.AvatarURL(),
		}
	}

	_, err = bot.PagedEmbed(ctx, embeds, 10*time.Minute)
	return
}
