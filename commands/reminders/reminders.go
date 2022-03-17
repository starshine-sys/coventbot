package reminders

import (
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
	bcr2 "github.com/starshine-sys/bcr/v2"
)

func (bot *Bot) remindersSlash(ctx *bcr2.CommandContext) (err error) {
	rms := []Reminder{}

	err = pgxscan.Select(context.Background(), bot.DB.Pool, &rms, "select * from reminders where user_id = $1 order by expires asc", ctx.User.ID)
	if err != nil {
		bot.ReportInteraction(ctx, err)
	}

	title := "Reminders"

	if len(rms) == 0 {
		return ctx.Reply("You have no reminders. Set some with `/remindme`!")
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
<t:%v> ([link](https://discord.com/channels/%v/%v/%v))

`, r.ID, text, r.Expires.Unix(), linkServer, r.ChannelID, r.MessageID))
	}

	embeds := bcr.StringPaginator(fmt.Sprintf("%v (%v)", title, len(rms)), bcr.ColourBlurple, slice, 5)

	_, _, err = ctx.Paginate(bcr2.PaginateEmbeds(embeds...), 10*time.Minute)
	return
}

func (bot *Bot) reminders(ctx *bcr.Context) (err error) {
	rms := []Reminder{}

	err = pgxscan.Select(context.Background(), bot.DB.Pool, &rms, "select * from reminders where user_id = $1 order by expires asc", ctx.Author.ID)
	if err != nil {
		bot.Report(ctx, err)
	}

	title := "Reminders"

	if len(rms) == 0 {
		return ctx.SendX("You have no reminders. Set some with `/remindme`!")
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
<t:%v> ([link](https://discord.com/channels/%v/%v/%v))

`, r.ID, text, r.Expires.Unix(), linkServer, r.ChannelID, r.MessageID))
	}

	embeds := bcr.StringPaginator(fmt.Sprintf("%v (%v)", title, len(rms)), bcr.ColourBlurple, slice, 5)

	_, err = bot.PagedEmbed(ctx, embeds, 10*time.Minute)
	return
}
