package reminders

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) reminders(ctx *bcr.Context) (err error) {
	rms := []Reminder{}

	err = pgxscan.Select(context.Background(), bot.DB.Pool, &rms, "select * from reminders where user_id = $1 order by expires asc", ctx.Author.ID)
	if err != nil {
		bot.Report(ctx, err)
	}

	if len(rms) == 0 {
		_, err = ctx.Sendf("You have no reminders. Set some with `%vremindme`!", ctx.Prefix)
		return
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

	_, err = bot.PagedEmbed(ctx, bcr.StringPaginator(fmt.Sprintf("Reminders (%v)", len(rms)), bcr.ColourBlurple, slice, 5), 10*time.Minute)
	return
}
