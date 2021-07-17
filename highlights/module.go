package highlights

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

// Bot ...
type Bot struct {
	*bot.Bot

	UserExpiration   map[userExpirationKey]time.Time
	UserExpirationMu sync.Mutex

	WordExpiration   map[wordExpirationKey]time.Time
	WordExpirationMu sync.Mutex
}

type userExpirationKey struct {
	UserID  discord.UserID
	GuildID discord.GuildID
}

type wordExpirationKey struct {
	UserID  discord.UserID
	GuildID discord.GuildID
	Word    string
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Highlights"

	b := &Bot{
		UserExpiration: map[userExpirationKey]time.Time{},
		WordExpiration: map[wordExpirationKey]time.Time{},
		Bot:            bot,
	}

	b.Router.AddHandler(b.messageCreate)

	hl := bot.Router.AddCommand(&bcr.Command{
		Name:      "hl",
		Aliases:   []string{"highlight"},
		Summary:   "Show your highlighted words.",
		Usage:     "[user]",
		GuildOnly: true,
		Command:   b.listHl,
	})

	hl.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a highlight.",
		Usage:   "<word>",
		Args:    bcr.MinArgs(1),

		GuildOnly: true,
		Command:   b.addHl,
	})

	hl.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Aliases: []string{"del", "delete", "rm"},
		Summary: "Remove a highlight.",
		Usage:   "<word>",
		Args:    bcr.MinArgs(1),

		GuildOnly: true,
		Command:   b.delHl,
	})

	hl.AddSubcommand(&bcr.Command{
		Name:    "block",
		Summary: "Block a user or channel from your highlights.",
		Usage:   "<user|channel>",
		Args:    bcr.MinArgs(1),

		GuildOnly: true,
		Command:   b.hlBlock,
	})

	hl.AddSubcommand(&bcr.Command{
		Name:    "unblock",
		Summary: "Unblock a user or channel from your highlights.",
		Usage:   "<user|channel>",
		Args:    bcr.MinArgs(1),

		GuildOnly: true,
		Command:   b.hlUnblock,
	})

	hl.AddSubcommand(&bcr.Command{
		Name:    "test",
		Aliases: []string{"match", "matches"},
		Summary: "Test whether a given string matches any of your highlights.",
		Usage:   "<test message>",
		Args:    bcr.MinArgs(1),

		GuildOnly: true,
		Command:   b.hlTest,
	})

	conf := hl.AddSubcommand(&bcr.Command{
		Name:    "config",
		Aliases: []string{"conf", "cfg"},
		Summary: "Configure highlights.",

		CustomPermissions: b.ModRole,
		Command:           b.showHlConfig,
	})

	conf.AddSubcommand(&bcr.Command{
		Name:    "toggle",
		Summary: "Enable or disable highlights.",

		OwnerOnly:         true,
		CustomPermissions: b.ModRole,
		Command:           b.toggleHl,
	})

	conf.AddSubcommand(&bcr.Command{
		Name:    "block",
		Summary: "Block a channel or category from highlights.",
		Usage:   "<channel>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: b.ModRole,
		Command:           b.modBlockHl,
	})

	conf.AddSubcommand(&bcr.Command{
		Name:    "unblock",
		Summary: "Unblock a channel or category from highlights.",
		Usage:   "<channel>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: b.ModRole,
		Command:           b.modHlUnblock,
	})

	state, _ := bot.Router.StateFromGuildID(0)

	var o sync.Once
	state.AddHandler(func(_ *gateway.ReadyEvent) {
		o.Do(func() {
			go b.deleteQueue(state)
		})
	})

	return s, append(list, hl)
}

type msg struct {
	MessageID discord.MessageID
	ChannelID discord.ChannelID
}

func (bot *Bot) deleteQueue(s *state.State) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case <-sc:
			break
		default:
		}

		msgs := []msg{}

		err := pgxscan.Select(context.Background(), bot.DB.Pool, &msgs, "select * from highlight_delete_queue where message_id < $1 limit 10", discord.NewSnowflake(time.Now().Add(-24*time.Hour)))
		if err != nil {
			bot.Sugar.Errorf("Error getting messages: %v", err)
			time.Sleep(time.Second)
			continue
		}

		for _, m := range msgs {
			_, err = bot.DB.Pool.Exec(context.Background(), "delete from highlight_delete_queue where message_id = $1", m.MessageID)
			if err != nil {
				bot.Sugar.Errorf("Error removing message from db: %v", err)
			}

			err = s.DeleteMessage(m.ChannelID, m.MessageID)
			if err != nil {
				bot.Sugar.Errorf("Error deleting message: %v", err)
			}
		}

		time.Sleep(time.Second)
	}
}
