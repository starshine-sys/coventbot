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
		g, err := ctx.State.Guild(ctx.Message.GuildID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		e := discord.Embed{
			Title: "Level config for " + g.Name,

			Color: bcr.ColourBlurple,
		}

		sc, err := bot.getGuildConfig(ctx.Message.GuildID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		e.Description = fmt.Sprintf("`levels_enabled`: %v\n`leaderboard_mod_only`: %v\n`show_next_reward`: %v\n`between_xp`: %v", sc.LevelsEnabled, sc.LeaderboardModOnly, sc.ShowNextReward, sc.BetweenXP)

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

	switch strings.ToLower(ctx.Args[0]) {
	case "levels_enabled", "leaderboard_mod_only", "show_next_reward":
		b, err := strconv.ParseBool(ctx.Args[1])
		if err != nil {
			_, err = ctx.Send("You must give either `true` or `false` for the new value.", nil)
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
	case "reward_text", "reward":
		text := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))
		if len(text) >= 1024 {
			_, err = ctx.Send("Input too long, maximum of 1024 characters.", nil)
			return
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
	}

	return
}
