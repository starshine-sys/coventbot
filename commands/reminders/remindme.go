package reminders

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"codeberg.org/eviedelta/detctime/durationparser"
	"github.com/diamondburned/arikawa/v3/bot/extras/shellwords"
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

	var embedless bool
	bot.DB.Pool.QueryRow(context.Background(), "select embedless_reminders from user_config where user_id = $1", ctx.Author.ID).Scan(&embedless)

	content := ""
	e := []discord.Embed{}

	if embedless {
		if len(rm) > 128 {
			rm = rm[:128] + "..."
		}
		if rm == "N/A" {
			rm = "something"
		} else {
			rm = "**" + rm + "**"
		}

		content = fmt.Sprintf("Okay %v, I'll remind you about %v %v. (#%v)", ctx.DisplayName(), rm, bcr.HumanizeTime(bcr.DurationPrecisionSeconds, t.Add(time.Second)), id)
	} else {
		e = []discord.Embed{{
			Color:       bcr.ColourGreen,
			Description: fmt.Sprintf("Reminder #%v set for %v from now.\n(%v UTC)", id, bcr.HumanizeDuration(bcr.DurationPrecisionSeconds, t.Sub(time.Now())+time.Second), t.Format("2006-01-02 15:04:05")),
		}}

		// only show this "ad" every few reminders
		if rand.Intn(2) == 1 {
			e[0].Footer = &discord.EmbedFooter{
				Text: "Did you know? You can use `" + ctx.Prefix + "usercfg` to make reminders more compact!",
			}
		}
	}

	_, err = ctx.Send(content, e...)
	return err
}

func (bot *Bot) remindmeSlash(ctx bcr.Contexter) (err error) {
	var t time.Time
	rm := ctx.GetStringFlag("text")
	if rm == "" {
		rm = "N/A"
	}

	when := ctx.GetStringFlag("when")
	args, err := shellwords.Parse(when)
	if err != nil {
		args = strings.Fields(when)
	}
	t, _, err = ParseTime(args)
	if err != nil {
		dur, err := durationparser.Parse(when)
		if err != nil {
			return ctx.SendEphemeral("I couldn't parse your input as a valid time or duration.")
		}
		t = time.Now().UTC().Add(dur)
	}

	guildID := discord.GuildID(0)
	if ctx.GetGuild() != nil {
		guildID = ctx.GetGuild().ID
	}
	// as there isn't a message associated with a slash command, we just use an approximate message ID
	// it'll still link to the correct(ish) time
	msgID := discord.NewSnowflake(time.Now())

	var id uint64
	err = bot.DB.Pool.QueryRow(context.Background(), `insert into reminders
	(user_id, message_id, channel_id, server_id, reminder, expires)
	values
	($1, $2, $3, $4, $5, $6) returning id`, ctx.User().ID, msgID, ctx.GetChannel().ID, guildID, rm, t).Scan(&id)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(rm) > 128 {
		rm = rm[:128] + "..."
	}
	if rm == "N/A" {
		rm = "something"
	} else {
		rm = "**" + rm + "**"
	}

	name := ctx.User().Username
	if v, ok := ctx.(*bcr.SlashContext); ok {
		if v.Member != nil && v.Member.Nick != "" {
			name = v.Member.Nick
		}
	}

	msg, err := ctx.Send(fmt.Sprintf("Okay %v, I'll remind you about %v %v. (#%v)", name, rm, bcr.HumanizeTime(bcr.DurationPrecisionSeconds, t.Add(time.Second)), id), discord.Embed{
		Description: t.Format("2006-01-02 15:04:05") + " UTC",
		Color:       bcr.ColourGreen,
	})
	if err != nil {
		return err
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update reminders set message_id = $1 where id = $2", msg.ID, id)
	return
}
