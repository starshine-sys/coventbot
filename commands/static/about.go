package static

import (
	"fmt"
	"runtime"
	"time"

	"github.com/starshine-sys/coventbot/etc"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

const botVersion = 1

func (bot *Bot) about(ctx *bcr.Context) (err error) {
	embed := discord.Embed{
		Title: "About",
		Color: ctx.Router.EmbedColor,

		Thumbnail: &discord.EmbedThumbnail{
			URL: ctx.Bot.AvatarURL(),
		},

		Fields: []discord.EmbedField{
			{
				Name:   "Version",
				Value:  fmt.Sprintf("v%d (bcr v%s)", botVersion, bcr.Version()),
				Inline: true,
			},
			{
				Name:   "Go version",
				Value:  runtime.Version(),
				Inline: true,
			},
			{
				Name:   "Creator",
				Value:  "<@!694563574386786314> / starshine system 🌠✨#0001",
				Inline: false,
			},
			{
				Name:   "Uptime",
				Value:  fmt.Sprintf("%v (since %v)", etc.HumanizeDuration(etc.DurationPrecisionSeconds, time.Since(bot.start)), bot.start.Format("Jan _2 2006, 15:04:05 MST")),
				Inline: false,
			},
			{
				Name:   "Source code",
				Value:  "https://github.com/starshine-sys/coventbot",
				Inline: false,
			},
		},

		Timestamp: discord.NowTimestamp(),
		Footer:    &discord.EmbedFooter{Text: "Made with Arikawa"},
	}

	_, err = ctx.Send("", &embed)
	return
}
