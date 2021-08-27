package reminders

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/spf13/pflag"
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

		SlashCommand: b.remindmeSlash,
		Options: &[]discord.CommandOption{
			{
				Name:        "when",
				Type:        discord.StringOption,
				Description: "When or in how long to remind you.",
				Required:    true,
			},
			{
				Name:        "text",
				Type:        discord.StringOption,
				Description: "What to remind you of.",
			},
		},
	})

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name: "reminders",

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.BoolP("channel", "c", false, "Only show reminders in this channel.")
			fs.BoolP("server", "s", false, "Only show reminders in this server.")

			return fs
		},

		Summary: "Show your reminders.",

		SlashCommand: b.reminders,
		Options:      &[]discord.CommandOption{},
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

	state, _ := bot.Router.StateFromGuildID(0)

	var o sync.Once
	state.AddHandler(func(_ *gateway.ReadyEvent) {
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
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case <-sc:
			break
		default:
		}

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

			state, _ := bot.Router.StateFromGuildID(r.ServerID)

			desc := fmt.Sprintf("%v you asked to be reminded about%v", bcr.HumanizeTime(bcr.DurationPrecisionSeconds, r.SetTime), reminder)
			if len(desc) > 2048 {
				desc = desc[:2040] + "..."
			}

			var shouldDM, embedless bool

			bot.DB.Pool.QueryRow(context.Background(), "select reminders_in_dm, embedless_reminders from user_config where user_id = $1", r.UserID).Scan(&shouldDM, &embedless)

			shouldDM = shouldDM || !r.ServerID.IsValid()

			if r.ServerID.IsValid() {
				shouldDM = true

				// this is Ugly™ but it should work
				// basically we need to get All of them to check perms
				m, err := bot.Member(r.ServerID, r.UserID)
				if err == nil {
					g, err := state.Guild(r.ServerID)
					if err == nil {
						ch, err := state.Channel(r.ChannelID)
						if err == nil {
							perms := discord.CalcOverwrites(*g, *ch, m)
							if perms.Has(discord.PermissionSendMessages | discord.PermissionViewChannel) {
								shouldDM = false
							}
						}
					}
				}
			}

			bot.Sugar.Debugf("Executing reminder #%v, should DM: %v, embedless: %v", r.ID, shouldDM, embedless)

			e := []discord.Embed{{
				Title:       fmt.Sprintf("Reminder #%v", r.ID),
				Description: desc,

				Color:     bcr.ColourBlurple,
				Timestamp: discord.NewTimestamp(r.SetTime),

				Fields: []discord.EmbedField{{
					Name:  "Link",
					Value: fmt.Sprintf("[Jump to message](https://discord.com/channels/%v/%v/%v)", linkServer, r.ChannelID, r.MessageID),
				}},
			}}

			data := api.SendMessageData{
				Content: r.UserID.Mention(),
				Embeds:  e,
				AllowedMentions: &api.AllowedMentions{
					Parse: []api.AllowedMentionType{api.AllowUserMention},
				},
			}

			if embedless {
				s := fmt.Sprintf("%v: %v (%v)", r.UserID.Mention(), r.Reminder, bcr.HumanizeTime(bcr.DurationPrecisionSeconds, r.SetTime))

				if len(s) <= 2000 {
					data.Content = s
					data.Embeds = nil
					data.Components = []discord.Component{
						discord.ActionRowComponent{
							Components: []discord.Component{
								discord.ButtonComponent{
									Label: "Jump to message",
									Style: discord.LinkButton,
									URL:   fmt.Sprintf("https://discord.com/channels/%v/%v/%v", linkServer, r.ChannelID, r.MessageID),
								},
							},
						},
					}
				}
			}

			switch shouldDM {
			case false:
				_, err = state.SendMessageComplex(r.ChannelID, data)
				if err == nil {
					bot.DB.Pool.Exec(context.Background(), "delete from reminders where id = $1", r.ID)
					continue
				}

				fallthrough
			case true:
				if data.Content == r.UserID.Mention() {
					data.Content = ""
				}

				ch, err := state.CreatePrivateChannel(r.UserID)
				if err != nil {
					bot.Sugar.Errorf("Error sending reminder %v: %v", r.ID, err)
					bot.DB.Pool.Exec(context.Background(), "delete from reminders where id = $1", r.ID)
					continue
				}

				_, err = state.SendMessageComplex(ch.ID, data)
				if err != nil {
					bot.Sugar.Errorf("Error sending reminder %v: %v", r.ID, err)
				}

				bot.DB.Pool.Exec(context.Background(), "delete from reminders where id = $1", r.ID)
				continue
			}
		}

		time.Sleep(time.Second)
	}
}
