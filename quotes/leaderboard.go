package quotes

import (
	"context"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

type userLeaderboard struct {
	UserID discord.UserID
	Quotes int64
}

type channelLeaderboard struct {
	ChannelID discord.ChannelID
	Quotes    int64
}

func (bot *Bot) leaderboard(ctx *bcr.Context) (err error) {
	if !bot.quotesEnabled(ctx.Guild.ID) {
		_, err = ctx.Send("Quotes aren't enabled on this server, sorry :(\nAsk a server admin to enable it!")
		return
	}

	var s []string
	var title string

	switch strings.ToLower(ctx.RawArgs) {
	case "chan", "ch", "channel", "channels":
		title = "Channel"
		var lb []channelLeaderboard
		err = pgxscan.Select(context.Background(), bot.DB.Pool, &lb, "select distinct channel_id, count(*) as quotes from quotes where server_id = $1 group by channel_id order by quotes desc, channel_id asc", ctx.Guild.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		for i, c := range lb {
			s = append(s, fmt.Sprintf("%d. <#%v>: `%v` quotes\n", i+1, c.ChannelID, humanize.Comma(c.Quotes)))
		}
	case "user", "users", "members":
		title = "User"
		var lb []userLeaderboard
		err = pgxscan.Select(context.Background(), bot.DB.Pool, &lb, "select distinct user_id, count(*) as quotes from quotes where server_id = $1 group by user_id order by quotes desc, user_id asc", ctx.Guild.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		for i, u := range lb {
			s = append(s, fmt.Sprintf("%d. <@!%v>: `%v` quotes\n", i+1, u.UserID, humanize.Comma(u.Quotes)))
		}
	}

	_, err = ctx.PagedEmbed(bcr.StringPaginator(title+" quote leaderboard", bcr.ColourBlurple, s, 15), false)
	return
}
