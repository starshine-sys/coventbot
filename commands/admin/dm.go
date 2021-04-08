package admin

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) dm(ctx *bcr.Context) (err error) {
	u, err := ctx.ParseUser(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("I could not find that user.", nil)
		return
	}

	msg := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))

	m, err := ctx.NewDM(u.ID).Embed(&discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: ctx.Bot.Username + " Admin",
			Icon: ctx.Bot.AvatarURL(),
		},
		Description: msg,
		Timestamp:   discord.NowTimestamp(),
		Color:       ctx.Router.EmbedColor,
	}).Send()
	if err != nil {
		_, err = ctx.Sendf("Error sending message to %v: %v", u.ID, err)
		return
	}

	_, err = ctx.Send("> DM sent!", &discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: ctx.DisplayName(),
			Icon: ctx.Author.AvatarURL(),
		},
		Description: msg,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v", m.ID),
		},
		Timestamp: discord.NowTimestamp(),
		Color:     ctx.Router.EmbedColor,
	})
	return
}
