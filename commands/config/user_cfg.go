// SPDX-License-Identifier: AGPL-3.0-only
package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) userCfg(ctx *bcr.Context) (err error) {
	if len(ctx.Args) == 0 {
		_, err = ctx.Send("", discord.Embed{
			Title:       "User configuration",
			Description: fmt.Sprintf("To enable any of these, use `%vusercfg` with the name and `true`; for example: `%vusercfg embedless_reminders true`. To disable them, run the same command but with `false` instead of `true`.\n\nTo show your current configuration, use `%vusercfg show`", ctx.Prefix, ctx.Prefix, ctx.Prefix),
			Fields: []discord.EmbedField{
				{
					Name:  "`disable_levelup_messages`",
					Value: "Disables level-up DMs, even if the server has them enabled.",
				},
				{
					Name:  "`reminders_in_dm`",
					Value: "If this is enabled, your reminders are always sent in a DM, even if the bot can send messages in the source channel.",
				},
				{
					Name:  "`embedless_reminders`",
					Value: "Sends reminder messages without embeds (except for the jump link), as long as the text fits in a normal message.",
				},
				{
					Name:  "`reaction_pages`",
					Value: "Makes paginated messages use reactions instead of buttons.",
				},
				{
					Name:  "`timezone`",
					Value: "Set your time zone, used in `remindme`.",
				},
			},
			Color: ctx.Router.EmbedColor,
		})
		return
	}

	if strings.EqualFold(ctx.RawArgs, "show") {
		lvlMsg, err := bot.DB.UserBoolGet(ctx.Author.ID, "disable_levelup_messages")
		if err != nil {
			return bot.Report(ctx, err)
		}

		dmReminders, err := bot.DB.UserBoolGet(ctx.Author.ID, "reminders_in_dm")
		if err != nil {
			return bot.Report(ctx, err)
		}

		embedless, err := bot.DB.UserBoolGet(ctx.Author.ID, "embedless_reminders")
		if err != nil {
			return bot.Report(ctx, err)
		}

		reactionPages, err := bot.DB.UserBoolGet(ctx.Author.ID, "reaction_pages")
		if err != nil {
			return bot.Report(ctx, err)
		}

		timezone, err := bot.DB.UserStringGet(ctx.Author.ID, "timezone")
		if err != nil {
			return bot.Report(ctx, err)
		}

		if timezone == "" {
			timezone = "UTC"
		}

		_, err = ctx.Send("", discord.Embed{
			Title: "User configuration",
			Description: fmt.Sprintf(
				"`disable_levelup_messages`: %v\n`reminders_in_dm`: %v\n`embedless_reminders`: %v\n`reaction_pages`: %v\n`timezone`: %v",
				lvlMsg, dmReminders, embedless, reactionPages, timezone,
			),
			Color: ctx.Router.EmbedColor,
		})
		return err
	}

	if len(ctx.Args) != 2 {
		_, err = ctx.Send("Too few or too many arguments given.")
		return
	}

	switch strings.ToLower(ctx.Args[0]) {
	case "disable_levelup_messages", "reminders_in_dm", "embedless_reminders", "reaction_pages":
		b, err := strconv.ParseBool(ctx.Args[1])
		if err != nil {
			_, err = ctx.Send("Couldn't parse your input as a boolean (true or false)")
			return err
		}

		err = bot.DB.UserBoolSet(ctx.Author.ID, strings.ToLower(ctx.Args[0]), b)
		if err != nil {
			return bot.Report(ctx, err)
		}

		_, err = ctx.Sendf("Set `%v` to `%v`!", ctx.Args[0], b)
		return err
	case "timezone":
		loc, err := time.LoadLocation(ctx.Args[1])
		if err != nil {
			_, err = ctx.Replyc(bcr.ColourRed, "I couldn't find a timezone named ``%v``.\nTimezone should be in `Continent/City` format; to find your timezone, use a tool such as <https://xske.github.io/tz/>.", bcr.EscapeBackticks(ctx.Args[1]))
			return err
		}

		err = bot.DB.UserStringSet(ctx.Author.ID, "timezone", loc.String())
		if err != nil {
			return bot.Report(ctx, err)
		}

		_, err = ctx.Replyc(bcr.ColourGreen, "Set your timezone to %v.", loc.String())
		return err
	default:
		_, err = ctx.Send("I don't recognise that config key.")
		return err
	}
}
