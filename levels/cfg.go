package levels

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) config(ctx *bcr.Context) (err error) {
	if len(ctx.Args) == 0 {
		e := discord.Embed{
			Title: "Level config for " + ctx.Guild.Name,

			Color: bcr.ColourBlurple,
		}

		sc, err := bot.getGuildConfig(ctx.Message.GuildID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		e.Description = fmt.Sprintf("`levels_enabled`: %v\n`leaderboard_mod_only`: %v\n`show_next_reward`: %v\n`between_xp`: %v", sc.LevelsEnabled, sc.LeaderboardModOnly, sc.ShowNextReward, sc.BetweenXP)

		if sc.RewardLog.IsValid() {
			e.Description += "\n`reward_log`: " + sc.RewardLog.Mention()
		} else {
			e.Description += "\n`reward_log`: None"
		}

		if sc.NolevelsLog.IsValid() {
			e.Description += "\n`nolevels_log`: " + sc.NolevelsLog.Mention()
		} else {
			e.Description += "\n`nolevels_log`: None"
		}

		rewards, err := bot.getAllRewards(ctx.Message.GuildID)
		if err == nil && len(rewards) > 0 {
			var buf string
			for _, r := range rewards {
				buf += fmt.Sprintf("%v: %v\n", r.Lvl, r.RoleReward.Mention())
			}
			e.Fields = append(e.Fields, discord.EmbedField{
				Name:  "Rewards",
				Value: buf,
			})
		}

		if sc.RewardText != "" {
			e.Fields = append(e.Fields, discord.EmbedField{
				Name:  "Reward message",
				Value: "```" + sc.RewardText + "```",
			})
		}

		_, err = ctx.Send("", &e)
		return err
	}

	if len(ctx.Args) < 2 {
		_, err = ctx.Send("Not enough arguments: you must give both a key and a new value.", nil)
		return
	}

	switch strings.ToLower(ctx.Args[0]) {
	case "levels_enabled", "leaderboard_mod_only", "show_next_reward":
		b, err := strconv.ParseBool(ctx.Args[1])
		if err != nil {
			_, err = ctx.Send("You must give either `true` or `false` for the new value.", nil)
			return err
		}
		_, err = bot.DB.Pool.Exec(
			context.Background(),
			"update server_levels set "+ctx.Args[0]+" = $1 where id = $2",
			b, ctx.Message.GuildID,
		)
		if err != nil {
			return bot.Report(ctx, err)
		}
		_, err = ctx.SendEmbed(bcr.SED{
			Message: fmt.Sprintf("Set `%v` to `%v`.", ctx.Args[0], b),
		})
		return err
	case "between_xp":
		t, err := time.ParseDuration(ctx.Args[1])
		if err != nil {
			_, err = ctx.Send("Couldn't parse your input as a valid duration.", nil)
			return err
		}
		_, err = bot.DB.Pool.Exec(
			context.Background(),
			"update server_levels set between_xp = $1 where id = $2",
			t, ctx.Message.GuildID,
		)
		if err != nil {
			return bot.Report(ctx, err)
		}
		_, err = ctx.SendEmbed(bcr.SED{
			Message: fmt.Sprintf("Set `%v` to `%v`.", ctx.Args[0], t),
		})
		return err
	case "reward_text", "reward":
		text := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))
		if len(text) >= 1024 {
			_, err = ctx.Send("Input too long, maximum of 1024 characters.", nil)
			return
		}

		if text == "clear" || text == "-clear" || text == "--clear" {
			text = ""
		}

		_, err = bot.DB.Pool.Exec(
			context.Background(),
			"update server_levels set reward_text = $1 where id = $2",
			text, ctx.Message.GuildID,
		)
		if err != nil {
			return bot.Report(ctx, err)
		}
		_, err = ctx.SendEmbed(bcr.SED{
			Title:   "Reward message updated",
			Message: fmt.Sprintf("```%v```", text),
		})
		return err
	case "reward_log", "nolevels_log":
		id := discord.ChannelID(0)
		if ctx.Args[1] != "clear" {
			ch, err := ctx.ParseChannel(ctx.Args[1])
			if err != nil || ch.Type != discord.GuildText || ch.GuildID != ctx.Message.GuildID {
				_, err = ctx.Send("I couldn't find that channel.", nil)
				return err
			}
			id = ch.ID
		}

		_, err = bot.DB.Pool.Exec(
			context.Background(),
			"update server_levels set "+ctx.Args[0]+" = $1 where id = $2",
			id, ctx.Message.GuildID,
		)
		if err != nil {
			return bot.Report(ctx, err)
		}

		if id == 0 {
			_, err = ctx.SendEmbed(bcr.SED{
				Message: fmt.Sprintf("Cleared `%v`.", ctx.Args[0]),
			})
			return
		}

		_, err = ctx.SendEmbed(bcr.SED{
			Message: fmt.Sprintf("Set `%v` to %v.", ctx.Args[0], id.Mention()),
		})
		return
	}

	return
}
