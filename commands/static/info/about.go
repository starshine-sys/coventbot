// SPDX-License-Identifier: AGPL-3.0-only
package info

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/starshine-sys/bcr"
	"gitlab.com/1f320/x/duration"
)

var GitVer string

func init() {
	if GitVer == "" {
		log.Println("Warning: GitVer is empty, falling back to checking at runtime.")
		git := exec.Command("git", "rev-parse", "--short", "HEAD")
		// ignoring errors *should* be fine? if there's no output we just fall back to "unknown"
		b, _ := git.Output()
		GitVer = string(b)
		if GitVer == "" {
			GitVer = "unknown"
		}
	}
}

func (bot *Bot) about(ctx *bcr.Context) (err error) {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	creator := "<@!694563574386786314>"
	u, err := ctx.State.User(694563574386786314)
	if err == nil {
		creator = fmt.Sprintf("<@!%v> / %v", u.ID, u.Tag())
	}

	bot.Counters.Mu.Lock()
	defer bot.Counters.Mu.Unlock()
	embed := discord.Embed{
		Title: "About",
		Color: ctx.Router.EmbedColor,

		Thumbnail: &discord.EmbedThumbnail{
			URL: ctx.Bot.AvatarURL(),
		},

		Fields: []discord.EmbedField{
			{
				Name:   "Version",
				Value:  fmt.Sprintf("%v (bcr v%s)", GitVer, bcr.Version()),
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
				Value:  fmt.Sprintf("%v (since <t:%v>)", duration.Format(time.Since(bot.start)), bot.start.Unix()),
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

	_, err = ctx.Send("", embed)
	return
}
