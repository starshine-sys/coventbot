package info

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) colour(ctx *bcr.Context) (err error) {
	ctx.RawArgs = strings.TrimPrefix(ctx.RawArgs, "#")

	clr, err := strconv.ParseUint(ctx.RawArgs, 16, 0)
	if err != nil {
		_, err = ctx.Send("You didn't give a valid hex code colour.")
		return
	}

	url := fmt.Sprintf("https://fakeimg.pl/256x256/%06X/?text=%%20", clr)
	r, g, b := discord.Color(clr).RGB()
	_, err = ctx.Send("", discord.Embed{
		Thumbnail: &discord.EmbedThumbnail{
			URL: url,
		},
		Title:       fmt.Sprintf("#%06X", clr),
		Color:       discord.Color(clr),
		Description: fmt.Sprintf("**R:** %v **G:** %v **B:** %v", r, g, b),
	})
	return
}
