package levels

import (
	"fmt"
	"math"

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
			currentLevel(l.XP),
		))
	}

	name := "Leaderboard"
	g, err := ctx.State.Guild(ctx.Message.GuildID)
	if err == nil {
		name = "Leaderboard for " + g.Name
	}

	_, err = ctx.PagedEmbed(
		StringPaginator(name, bcr.ColourBlurple, strings, 15),
		true,
	)
	return err
}

// StringPaginator paginates strings, for use in ctx.PagedEmbed
func StringPaginator(title string, colour discord.Color, slice []string, perPage int) []discord.Embed {
	var (
		embeds []discord.Embed
		count  int

		pages = 1
		buf   = discord.Embed{
			Title: title,
			Color: colour,
			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("Page 1/%v", math.Ceil(float64(len(slice))/float64(perPage))),
			},
		}
	)

	for _, s := range slice {
		if count >= perPage {
			embeds = append(embeds, buf)
			buf = discord.Embed{
				Title: title,
				Color: colour,
				Footer: &discord.EmbedFooter{
					Text: fmt.Sprintf("Page %v/%v", pages+1, math.Ceil(float64(len(slice))/float64(perPage))),
				},
			}
			count = 0
			pages++
		}
		buf.Description += s
		count++
	}

	embeds = append(embeds, buf)

	return embeds
}
