package static

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) enlarge(ctx *bcr.Context) (err error) {
	if !emojiMatch.MatchString(ctx.RawArgs) {
		_, err = ctx.Send("You didn't give a __custom__ emoji to enlarge.", nil)
		return err
	}

	extension := ".png"
	groups := emojiMatch.FindStringSubmatch(ctx.RawArgs)
	if groups[1] == "a" {
		extension = ".gif"
	}
	url := fmt.Sprintf("https://cdn.discordapp.com/emojis/%v%v", groups[3], extension)

	_, err = ctx.Send("", &discord.Embed{
		Image: &discord.EmbedImage{
			URL: url,
		},
		Color: ctx.Router.EmbedColor,
	})
	return err
}
