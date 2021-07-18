package config

import (
	"context"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

type channelStats struct {
	ChannelID discord.ChannelID
	Count     int64
}

func (bot *Bot) starboardStats(ctx *bcr.Context) (err error) {
	var stats []channelStats

	err = pgxscan.Select(context.Background(), bot.DB.Pool, &stats, "select distinct channel_id, count(*) as count from starboard_messages where server_id = $1 group by channel_id order by count desc, channel_id asc", ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(stats) == 0 {
		_, err = ctx.Replyc(bcr.ColourRed, "%v has no starboarded messages.", ctx.Guild.Name)
		return
	}

	var total int64
	var strings []string
	for i, st := range stats {
		total += st.Count
		s := fmt.Sprintf("%v. %v: %v message", i+1, st.ChannelID.Mention(), humanize.Comma(st.Count))
		if st.Count != 1 {
			s += "s"
		}

		strings = append(strings, s+"\n")
	}

	t := humanize.Comma(total)
	embeds := bcr.StringPaginator("Starboard statistics for "+ctx.Guild.Name, ctx.Router.EmbedColor, strings, 25)
	for i := range embeds {
		embeds[i].Footer.Text += " | " + t + " total message"
		if total != 1 {
			embeds[i].Footer.Text += "s"
		}
	}

	_, err = bot.PagedEmbed(ctx, embeds, 15*time.Minute)
	return
}
