package reminders

import (
	"context"
	"fmt"
	"strings"
	"time"

	"codeberg.org/eviedelta/detctime/durationparser"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) remindme(ctx *bcr.Context) (err error) {
	dur, err := durationparser.Parse(ctx.Args[0])
	if err != nil {
		_, err = ctx.Sendf("I couldn't parse ``%v`` as a valid duration.", bcr.EscapeBackticks(ctx.Args[0]))
	}

	rm := "N/A"
	if len(ctx.Args) > 1 {
		rm = strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))
		// if we couldn't trim the time off we just fall back to strings.Join
		if rm == ctx.RawArgs {
			rm = strings.Join(ctx.Args[1:], " ")
		}
	}

	var id uint64
	err = bot.DB.Pool.QueryRow(context.Background(), `insert into reminders
	(user_id, message_id, channel_id, server_id, reminder, expires)
	values
	($1, $2, $3, $4, $5, $6) returning id`, ctx.Author.ID, ctx.Message.ID, ctx.Channel.ID, ctx.Message.GuildID, rm, time.Now().UTC().Add(dur)).Scan(&id)
	if err != nil {
		return bot.Report(ctx, err)
	}

	e := discord.Embed{
		Color:       bcr.ColourGreen,
		Description: fmt.Sprintf("Reminder #%v set for %v from now.\n(%v UTC)", id, bcr.HumanizeDuration(bcr.DurationPrecisionMinutes, dur), time.Now().UTC().Add(dur).Format("2006-01-02 15:04:05")),
	}

	_, err = ctx.Send("", &e)
	return err
}
