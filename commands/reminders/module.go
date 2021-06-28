package reminders

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

// Reminder ...
type Reminder struct {
	ID        uint64
	UserID    discord.UserID
	MessageID discord.MessageID
	ChannelID discord.ChannelID
	ServerID  discord.GuildID

	Reminder string

	SetTime time.Time
	Expires time.Time

	ReminderInDM bool
}

// Bot ...
type Bot struct {
	*bot.Bot
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Reminders"

	b := &Bot{Bot: bot}

	rm := bot.Router.AddCommand(&bcr.Command{
		Name:    "remindme",
		Aliases: []string{"remind", "reminder"},

		Summary: "Set a reminder for yourself.",
		Usage:   "<time or duration> [reason]",
		Args:    bcr.MinArgs(1),

		Command: b.remindme,
	})

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name: "reminders",

		Summary: "Show your reminders.",

		Command: b.reminders,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "delreminder",
		Aliases: []string{"deletereminder", "delete-reminder", "delrm"},

		Summary: "Delete one of your reminders.",
		Usage:   "<id>",
		Args:    bcr.MinArgs(1),

		Command: b.delreminder,
	}))

	rm.AddSubcommand(b.Router.AliasMust("list", nil, []string{"reminders"}, nil))
	rm.AddSubcommand(b.Router.AliasMust("delete", []string{"remove", "rm", "del"}, []string{"delreminder"}, nil))

	var o sync.Once
	bot.State.AddHandler(func(_ *gateway.ReadyEvent) {
		o.Do(func() {
			go b.doReminders()
		})
	})

	return s, append(list, rm)
}

// doReminders gets 5 reminders at a time and executes them, then sleeps for 1 second.
// this should be fine unless we get >5 reminders a second,
// at which point we have bigger problems tbh
func (bot *Bot) doReminders() {
	for {
		rms := []Reminder{}

		err := pgxscan.Select(context.Background(), bot.DB.Pool, &rms, "select * from reminders where expires < $1 limit 5", time.Now().UTC())
		if err != nil {
			bot.Sugar.Errorf("Error getting reminders: %v", err)
			time.Sleep(time.Second)
			continue
		}

		for _, r := range rms {
			reminder := " something"
			if r.Reminder != "N/A" {
				reminder = fmt.Sprintf("\n%v", r.Reminder)
			}

			linkServer := r.ServerID.String()
			if !r.ServerID.IsValid() {
				linkServer = "@me"
			}

			desc := fmt.Sprintf("%v you asked to be reminded about%v", bcr.HumanizeTime(bcr.DurationPrecisionSeconds, r.SetTime), reminder)
			if len(desc) > 2048 {
				desc = desc[:2040] + "..."
			}

			e := discord.Embed{
				Title:       fmt.Sprintf("Reminder #%v", r.ID),
				Description: desc,

				Color:     bcr.ColourBlurple,
				Timestamp: discord.NewTimestamp(r.SetTime),

				Fields: []discord.EmbedField{{
					Name:  "Link",
					Value: fmt.Sprintf("[Jump to message](https://discord.com/channels/%v/%v/%v)", linkServer, r.ChannelID, r.MessageID),
				}},
			}

			if r.ServerID.IsValid() {
				var shouldDM bool
				bot.DB.Pool.QueryRow(context.Background(), "select reminders_in_dm from user_config where user_id = $1", r.UserID).Scan(&shouldDM)
				if !shouldDM {
					_, err = bot.State.SendMessage(r.ChannelID, r.UserID.Mention(), &e)
					if err == nil {
						bot.DB.Pool.Exec(context.Background(), "delete from reminders where id = $1", r.ID)
						continue
					}
				}
			}

			ch, err := bot.State.CreatePrivateChannel(r.UserID)
			if err != nil {
				bot.Sugar.Errorf("Error sending reminder %v: %v", r.ID, err)
				bot.DB.Pool.Exec(context.Background(), "delete from reminders where id = $1", r.ID)
				continue
			}

			_, err = bot.State.SendEmbed(ch.ID, e)
			if err != nil {
				bot.Sugar.Errorf("Error sending reminder %v: %v", r.ID, err)
				bot.DB.Pool.Exec(context.Background(), "delete from reminders where id = $1", r.ID)
				continue
			}
			bot.DB.Pool.Exec(context.Background(), "delete from reminders where id = $1", r.ID)
		}

		time.Sleep(time.Second)
	}
}
