package info

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) idtime(ctx *bcr.Context) (err error) {
	sf, err := discord.ParseSnowflake(ctx.Args[0])
	if err != nil {
		_, err = ctx.Sendf("Couldn't parse ``%v`` as a Discord ID.", bcr.EscapeBackticks(ctx.Args[0]))
		return err
	}

	_, err = ctx.Send("", &discord.Embed{
		Title:       fmt.Sprintf("Timestamp for `%v`", sf),
		Description: sf.Time().UTC().Format("January 02 2006, 15:04:05.000 MST"),
		Color:       ctx.Router.EmbedColor,
	})
	return
}
