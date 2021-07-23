package reminders

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) reminders(ctx bcr.Contexter) (err error) {
	rms := []Reminder{}

	err = pgxscan.Select(context.Background(), bot.DB.Pool, &rms, "select * from reminders where user_id = $1 order by expires asc", ctx.User().ID)
	if err != nil {
		bot.Report(ctx, err)
	}

	title := "Reminders"

	if v, ok := ctx.(*bcr.Context); ok {
		limitChannel, _ := v.Flags.GetBool("channel")
		if limitChannel {
			title = "Reminders in #" + v.Channel.Name
			prev := rms
			rms = nil
			for _, r := range prev {
				if r.ChannelID == v.Channel.ID {
					rms = append(rms, r)
				}
			}
		}

		limitServer, _ := v.Flags.GetBool("server")
		if limitServer && v.Guild != nil {
			title = "Reminders in " + v.Guild.Name
			prev := rms
			rms = nil
			for _, r := range prev {
				if r.ServerID == v.Guild.ID {
					rms = append(rms, r)
				}
			}
		}
	}

	if len(rms) == 0 {
		prefix := "/"
		if v, ok := ctx.(*bcr.Context); ok {
			prefix = v.Prefix
		}

		return ctx.SendfX("You have no reminders. Set some with `%vremindme`!", prefix)
	}

	var slice []string

	for _, r := range rms {
		text := r.Reminder
		if len(text) > 100 {
			text = text[:100] + "..."
		}

		linkServer := r.ServerID.String()
		if !r.ServerID.IsValid() {
			linkServer = "@me"
		}

		slice = append(slice, fmt.Sprintf(`**#%v**: %v
%v UTC ([link](https://discord.com/channels/%v/%v/%v))

`, r.ID, text, r.Expires.Format("2006-01-02 | 15:04"), linkServer, r.ChannelID, r.MessageID))
	}

	embeds := bcr.StringPaginator(fmt.Sprintf("%v (%v)", title, len(rms)), bcr.ColourBlurple, slice, 5)

	if v, ok := ctx.(*bcr.Context); ok {
		_, err = bot.PagedEmbed(v, embeds, 10*time.Minute)
		return
	}

	_, _, err = ctx.ButtonPages(embeds, 10*time.Minute)
	return
}
