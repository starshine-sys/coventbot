package admin

import (
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) activity(ctx *bcr.Context) (err error) {
	s := bot.Settings()

	if ctx.RawArgs == "" {
		e := discord.Embed{
			Title: "Status",
			Fields: []discord.EmbedField{
				{
					Name:  "Type",
					Value: s.ActivityType,
				},
				{
					Name: "Activity",
					Value: func(b bool) string {
						if b {
							return "<none>"
						}
						return "`" + s.Activity + "`"
					}(s.Activity == ""),
				},
			},
			Color: ctx.Router.EmbedColor,
		}

		_, err = ctx.Send("", &e)
		return
	}

	if ctx.RawArgs == "-clear" {
		s.ActivityType = "playing"
		s.Activity = ""

		err = bot.SetSettings(s)
		if err != nil {
			return bot.Report(ctx, err)
		}
		bot.updateStatus()

		_, err = ctx.Send("Cleared activity.", nil)
		return
	}

	if len(ctx.Args) < 2 {
		_, err = ctx.Send("You didn't give both an activity type and activity.", nil)
		return
	}

	if strings.HasPrefix(strings.ToLower(ctx.RawArgs), "playing") {
		s.ActivityType = "playing"
		s.Activity = ctx.RawArgs[len("playing")+1:]
	} else if strings.HasPrefix(strings.ToLower(ctx.RawArgs), "listening") {
		s.ActivityType = "listening to"
		s.Activity = ctx.RawArgs[len("listening")+1:]
		if strings.HasPrefix(strings.ToLower(s.Activity), "to") {
			s.Activity = s.Activity[len("to")+1:]
		}
	} else if strings.HasPrefix(strings.ToLower(ctx.RawArgs), "watching") {
		s.ActivityType = "watching"
		s.Activity = ctx.RawArgs[len("watching")+1:]
	}

	s.Activity = strings.TrimSpace(s.Activity)

	err = bot.SetSettings(s)
	if err != nil {
		return bot.Report(ctx, err)
	}

	bot.updateStatus()

	_, err = ctx.Sendf("Updated status to `%v %v`!", s.ActivityType, s.Activity)
	return
}
