package info

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/dustin/go-humanize"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/etc"
)

const botVersion = 5

var gitVer string

func init() {
	git := exec.Command("git", "rev-parse", "--short", "HEAD")
	// ignoring errors *should* be fine? if there's no output we just fall back to "unknown"
	b, _ := git.Output()
	gitVer = string(b)
	if gitVer == "" {
		gitVer = "unknown"
	}
}

func (bot *Bot) about(ctx *bcr.Context) (err error) {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	creator := "<@!694563574386786314>"
	u, err := ctx.State.User(694563574386786314)
	if err == nil {
		creator = fmt.Sprintf("<@!%v> / %v#%v", u.ID, u.Username, u.Discriminator)
	}

	bot.Counters.Mu.Lock()
	embed := discord.Embed{
		Title: "About",
		Color: ctx.Router.EmbedColor,

		Thumbnail: &discord.EmbedThumbnail{
			URL: ctx.Bot.AvatarURL(),
		},

		Fields: []discord.EmbedField{
			{
				Name:   "Version",
				Value:  fmt.Sprintf("v%d-%v (bcr v%s)", botVersion, gitVer, bcr.Version()),
				Inline: true,
			},
			{
				Name:   "Go version",
				Value:  runtime.Version(),
				Inline: true,
			},
			{
				Name:   "Creator",
				Value:  creator,
				Inline: false,
			},
			{
				Name:   "Uptime",
				Value:  fmt.Sprintf("%v (since %v)", etc.HumanizeDuration(etc.DurationPrecisionSeconds, time.Since(bot.start)), bot.start.Format("Jan _2 2006, 15:04:05 MST")),
				Inline: false,
			},
			{
				Name:   "Source code",
				Value:  "https://github.com/starshine-sys/tribble",
				Inline: false,
			},
			{
				Name: "Stats since last restart",
				Value: fmt.Sprintf(`Bot mentions: %v
Messages: %v

Memory used: %v / %v (%v garbage collected)
Goroutines: %v
Member cache size: %v members`, bot.Counters.Mentions, bot.Counters.Messages,
					humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys), humanize.Bytes(stats.TotalAlloc), humanize.Comma(int64(runtime.NumGoroutine())), humanize.Comma(bot.CacheLen())),
				Inline: false,
			},
		},

		Timestamp: discord.NowTimestamp(),
		Footer:    &discord.EmbedFooter{Text: "Made with Arikawa"},
	}
	bot.Counters.Mu.Unlock()

	_, err = ctx.Send("", &embed)
	return
}
