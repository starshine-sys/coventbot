package levels

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) leaderboard(ctx *bcr.Context) (err error) {
	sc, err := bot.getGuildConfig(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	noRanks, err := bot.DB.GuildBoolGet(ctx.Message.GuildID, "levels:disable_ranks")
	if err != nil {
		return bot.Report(ctx, err)
	}
	if noRanks {
		return ctx.SendX("Ranks are disabled on this server.")
	}

	if sc.LeaderboardModOnly || !sc.LevelsEnabled {
		perm, _ := bot.HelperRole.Check(ctx)

		if !perm {
			_, err = ctx.Sendf("You don't have permission to use this command, you need the **Manage Messages** permission to use it.")
			return
		}
	}

	full, _ := ctx.Flags.GetBool("full")
	lb, err := bot.getLeaderboard(ctx.Message.GuildID, full)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(lb) == 0 {
		_, err = ctx.Sendf("There doesn't seem to be anyone on the leaderboard...")
		return
	}

	var strings []string
	for i, l := range lb {
		strings = append(strings, fmt.Sprintf(
			"%v. %v: `%v` XP, level `%v`\n",
			i+1,
			l.UserID.Mention(),
			humanize.Comma(l.XP),
			sc.CalculateLevel(l.XP),
		))
	}

	name := "Leaderboard for " + ctx.Guild.Name

	embeds := bcr.StringPaginator(name, bcr.ColourBlurple, strings, 20)
	if bot.Config.HTTPBaseURL != "" {
		for i := range embeds {
			embeds[i].URL = fmt.Sprintf("%v/leaderboard/%v", bot.Config.HTTPBaseURL, ctx.Guild.ID)
		}
	}

	_, err = bot.PagedEmbed(ctx,
		embeds,
		10*time.Minute,
	)
	return err
}
