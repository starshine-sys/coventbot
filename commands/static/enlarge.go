package static

import (
	"fmt"
	"image"

	// to decode emoji
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"net/http"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/etc"
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

	resp, err := http.Get(url)
	if err != nil {
		return bot.Report(ctx, err)
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)

	r, g, b, _ := etc.AverageColour(img)

	clr := discord.Color(r)<<16 + discord.Color(g)<<8 + discord.Color(b)
	if clr == 0 {
		clr = bcr.ColourBlurple
	}

	_, err = ctx.Send("", &discord.Embed{
		Image: &discord.EmbedImage{
			URL: url,
		},
		Color: clr,
	})
	return err
}
