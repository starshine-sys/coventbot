package config

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

// UserConfig ...
type UserConfig struct {
	UserID discord.UserID

	DisableLevelupMessages bool
	RemindersInDM          bool
	UsernamesOptOut        bool
}

func (bot *Bot) userCfg(ctx *bcr.Context) (err error) {
	if len(ctx.Args) == 0 {
		_, err = ctx.Send("", &discord.Embed{
			Title: "User configuration",
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
					Name:  "`usernames_opt_out`",
					Value: "Disables username logging. Note that nickname changes will still be logged, but are limited to that specific server's moderators.",
				},
				{
					Name:  "`embedless_reminders`",
					Value: "Sends reminder messages without embeds (except for the jump link), as long as the text fits in a normal message.",
				},
			},
			Color: ctx.Router.EmbedColor,
		})
		return
	}

	if strings.EqualFold(ctx.RawArgs, "show") {
		uc := UserConfig{}

		pgxscan.Get(context.Background(), bot.DB.Pool, &uc, "select * from user_config where user_id = $1", ctx.Author.ID)

		_, err = ctx.Send("", &discord.Embed{
			Title:       "User configuration",
			Description: fmt.Sprintf("`disable_levelup_messages`: %v\n`reminders_in_dm`: %v\n`usernames_opt_out`: %v", uc.DisableLevelupMessages, uc.RemindersInDM, uc.UsernamesOptOut),
			Color:       ctx.Router.EmbedColor,
		})
		return
	}

	if len(ctx.Args) != 2 {
		_, err = ctx.Send("Too few or too many arguments given.", nil)
		return
	}

	switch strings.ToLower(ctx.Args[0]) {
	case "disable_levelup_messages", "reminders_in_dm", "usernames_opt_out", "embedless_reminders":
		b, err := strconv.ParseBool(ctx.Args[1])
		if err != nil {
			_, err = ctx.Send("Couldn't parse your input as a boolean (true or false)", nil)
			return err
		}

		_, err = bot.DB.Pool.Exec(context.Background(), "insert into user_config (user_id, "+ctx.Args[0]+") values ($1, $2) on conflict (user_id) do update set "+ctx.Args[0]+" = $2", ctx.Author.ID, b)
		if err != nil {
			return bot.Report(ctx, err)
		}

		_, err = ctx.Sendf("Set `%v` to `%v`!", ctx.Args[0], b)
	default:
		_, err = ctx.Send("I don't recognise that config key.", nil)
	}
	return
}
