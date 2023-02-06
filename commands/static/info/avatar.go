// SPDX-License-Identifier: AGPL-3.0-only
package info

import (
	"fmt"
	"image/png"
	"net/http"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	bcr2 "github.com/starshine-sys/bcr/v2"
	"github.com/starshine-sys/tribble/etc"
)

func (bot *Bot) avatar(ctx *bcr.Context) (err error) {
	u := ctx.Author

	if len(ctx.Args) > 0 {
		m, err := ctx.ParseMember(ctx.RawArgs)
		if err == nil {
			u = m.User
		} else {
			user, err := ctx.ParseUser(ctx.RawArgs)
			if err == nil {
				u = *user
			}
		}
	}

	resp, err := http.Get(u.AvatarURLWithType(discord.PNGImage))
	if err != nil {
		return bot.Report(ctx, err)
	}
	defer resp.Body.Close()

	pfp, err := png.Decode(resp.Body)

	r, g, b, _ := etc.AverageColour(pfp)

	var clr discord.Color = bcr.ColourBlurple
	if r != 0 || g != 0 || b != 0 {
		clr = discord.Color(r)<<16 + discord.Color(g)<<8 + discord.Color(b)
	}

	_, err = ctx.Send("", discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: u.Username + "#" + u.Discriminator,
			Icon: u.AvatarURL(),
		},

		Description: fmt.Sprintf("[jpg](%v?size=1024) | [png](%v?size=1024) | [webp](%v?size=1024)", u.AvatarURLWithType(discord.JPEGImage), u.AvatarURLWithType(discord.PNGImage), u.AvatarURLWithType(discord.WebPImage)),

		Image: &discord.EmbedImage{
			URL: u.AvatarURL() + "?size=1024",
		},

		Color: clr,

		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("Average colour: #%02X%02X%02X", r, g, b),
		},
	})
	return
}

func (bot *Bot) avatarMenu(ctx *bcr2.CommandContext) (err error) {
	u := ctx.FirstUser()

	embeds := []discord.Embed{
		{
			Title:       "Avatar for " + u.Tag(),
			Description: fmt.Sprintf("[Link](%v)", u.AvatarURL()+"?size=1024"),
			Image: &discord.EmbedImage{
				URL: u.AvatarURL() + "?size=1024",
			},
			Color: bcr.ColourBlurple,
		},
	}
	if m, err := bot.Member(ctx.Event.GuildID, u.ID); err == nil {
		if m.Avatar != "" {
			embeds = append(embeds, discord.Embed{
				Title:       "Server-specific avatar",
				Description: fmt.Sprintf("[Link](%v)", m.AvatarURL(ctx.Event.GuildID)+"?size=1024"),
				Image: &discord.EmbedImage{
					URL: m.AvatarURL(ctx.Event.GuildID) + "?size=1024",
				},
				Color: bcr.ColourBlurple,
			})
		}
	}

	return ctx.ReplyEphemeral("", embeds...)
}
