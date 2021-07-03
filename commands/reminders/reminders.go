package reminders

import (
	"context"
	"fmt"
	"math"

	"github.com/diamondburned/arikawa/v3/discord"
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

	_, err = ctx.PagedEmbed(StringPaginator(fmt.Sprintf("Reminders (%v)", len(rms)), bcr.ColourBlurple, slice, 5), false)
	return
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
