package reminders

import (
	"context"
	"fmt"
	"strings"
	"time"

	"codeberg.org/eviedelta/detctime/durationparser"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) remindme(ctx *bcr.Context) (err error) {
	t, i, err := ParseTime(ctx.Args)
	if err != nil {
		dur, err := durationparser.Parse(ctx.Args[0])
		if err != nil {
			_, err = ctx.Replyc(bcr.ColourRed, "I couldn't parse your input as a valid time or duration.")
			return err
		}
		i = 0
		t = time.Now().UTC().Add(dur)
	}

	if t.Before(time.Now().UTC()) {
		_, err = ctx.Replyc(bcr.ColourRed, "That time is in the past.")
		return
	}

	rm := "N/A"

	if len(ctx.Args) > i+1 {
		rm = ctx.RawArgs
		for n := 0; n <= i; n++ {
			rm = strings.TrimSpace(strings.TrimPrefix(rm, ctx.Args[n]))
		}
	}

	if rm == ctx.RawArgs {
		rm = strings.Join(ctx.Args[i+1:], " ")
	}

	var id uint64
	err = bot.DB.Pool.QueryRow(context.Background(), `insert into reminders
	(user_id, message_id, channel_id, server_id, reminder, expires)
	values
	($1, $2, $3, $4, $5, $6) returning id`, ctx.Author.ID, ctx.Message.ID, ctx.Channel.ID, ctx.Message.GuildID, rm, t).Scan(&id)
	if err != nil {
		return bot.Report(ctx, err)
	}

	dur := t.Sub(time.Now().UTC()) + time.Minute

	e := discord.Embed{
		Color:       bcr.ColourGreen,
		Description: fmt.Sprintf("Reminder #%v set for %v from now.\n(%v UTC)", id, bcr.HumanizeDuration(bcr.DurationPrecisionMinutes, dur), t.Format("2006-01-02 15:04:05")),
	}

	if t.After(time.Now().Add(100 * 365 * 24 * time.Hour)) {
		e.Description += "\n\nthe time isn't correct, we know, no need to report it, that's what you get for setting it to something so big that Go doesn't support it ([see here](https://golang.org/pkg/time/#Duration))"
	}

	_, err = ctx.Send("", e)
	return err
}
