package levels

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/dustin/go-humanize"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) level(ctx *bcr.Context) (err error) {
	sc, err := bot.getGuildConfig(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}
	if !sc.LevelsEnabled {
		_, err = ctx.Send("Levels are not enabled on this server.", nil)
		return
	}

	u := &ctx.Author
	if len(ctx.Args) > 0 {
		u, err = ctx.ParseUser(ctx.RawArgs)
		if err != nil {
			_, err = ctx.Send("User not found.", nil)
			return
		}
	}

	uc, err := bot.getUser(ctx.Message.GuildID, u.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	lvl := currentLevel(uc.XP)
	xpForNext := expForNextLevel(lvl)
	xpForPrev := expForNextLevel(lvl - 1)

	var percent int64
	{
		progress := uc.XP - xpForPrev
		needed := xpForNext - xpForPrev

		p := float64(progress) / float64(needed)

		percent = int64(p * 100)
	}

	// get leaderboard (for rank)
	var rank int
	lb, err := bot.getLeaderboard(ctx.Message.GuildID)
	if err == nil {
		for i, uc := range lb {
			if uc.UserID == u.ID {
				rank = i + 1
				break
			}
		}
	}

	// get user colour
	clr := uc.Colour
	if clr == 0 {
		clr, err = ctx.State.MemberColor(ctx.Message.GuildID, u.ID)
		if err != nil || clr == 0 {
			clr = bcr.ColourBlurple
		}
	}

	e := discord.Embed{
		Thumbnail: &discord.EmbedThumbnail{
			URL: u.AvatarURLWithType(discord.PNGImage),
		},
		Title: fmt.Sprintf("Level %v - Rank #%v", lvl, rank),
		Description: fmt.Sprintf(
			"**%v%%**\n%v/%v XP", percent,
			humanize.Comma(uc.XP), humanize.Comma(xpForNext),
		),
		Color: clr,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("%v#%v", u.Username, u.Discriminator),
		},
		Timestamp: discord.NowTimestamp(),
	}

	// get next reward
	reward := bot.getNextReward(ctx.Message.GuildID, lvl)
	if reward != nil && sc.ShowNextReward {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Next reward",
			Value: fmt.Sprintf("%v\nat level %v", reward.RoleReward.Mention(), reward.Lvl),
		})
	} else if sc.ShowNextReward {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Next reward",
			Value: "No more rewards",
		})
	}

	_, err = ctx.Send("", &e)
	return
}
