package levels

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/dustin/go-humanize"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) leaderboard(ctx *bcr.Context) (err error) {
	sc, err := bot.getGuildConfig(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if sc.LeaderboardModOnly || !sc.LevelsEnabled {
		if perms, _ := bot.State.Permissions(ctx.Channel.ID, ctx.Author.ID); !perms.Has(discord.PermissionManageMessages) {
			_, err = ctx.Sendf("You don't have permission to use this command, you need the **Manage Messages** permission to use it.")
			return
		}
	}

	lb, err := bot.getLeaderboard(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	var display []Levels
	full, _ := ctx.Flags.GetBool("full")
	if full {
		display = lb
	} else {
		gm := bot.Members(ctx.Message.GuildID)
		for _, l := range lb {
			for _, m := range gm {
				if m.User.ID == l.UserID {
					display = append(display, l)
					break
				}
			}
		}
	}

	if len(display) == 0 {
		_, err = ctx.Sendf("There doesn't seem to be anyone on the leaderboard...")
		return
	}

	var strings []string
	for i, l := range display {
		strings = append(strings, fmt.Sprintf(
			"%v. %v: `%v` XP, level `%v`\n",
			i+1,
			l.UserID.Mention(),
			humanize.Comma(l.XP),
			currentLevel(l.XP),
		))
	}

	name := "Leaderboard for " + ctx.Guild.Name

	_, err = ctx.PagedEmbed(
		bcr.StringPaginator(name, bcr.ColourBlurple, strings, 15),
		true,
	)
	return err
}
