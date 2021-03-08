package static

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
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

	_, err = ctx.Send("", &discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: u.Username + "#" + u.Discriminator,
			Icon: u.AvatarURL(),
		},
		Image: &discord.EmbedImage{
			URL: u.AvatarURL() + "?size=1024",
		},
		Color: ctx.Router.EmbedColor,
	})
	return
}
