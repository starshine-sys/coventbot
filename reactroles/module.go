package reactroles

import (
	"context"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

// Bot ...
type Bot struct {
	*bot.Bot
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Reaction roles"

	b := &Bot{bot}

	rr := bot.Router.AddCommand(&bcr.Command{
		Name:    "reactroles",
		Aliases: []string{"rr"},
		Summary: "Create or edit reaction roles",

		Command: func(ctx *bcr.Context) (err error) {
			err = ctx.Help([]string{"reactroles"})
			return
		},
	})

	rr.AddSubcommand(&bcr.Command{
		Name:    "update",
		Summary: "Add or update reaction roles for the given message.",
		Usage:   "<message> <emote/role pairs...>",
		Args:    bcr.MinArgs(3),

		Command: b.update,
	})

	rr.AddSubcommand(&bcr.Command{
		Name:    "new",
		Summary: "Create a new reaction role message in the given channel.",
		Usage:   "<channel> <name> <emote/role pairs...>",
		Args:    bcr.MinArgs(4),

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.BoolP("mention", "m", false, "Show roles as @mentions")
			fs.StringP("description", "d", "", "A description to show before the role list")

			return fs
		},

		Command: b.new,
	})

	simple := rr.AddSubcommand(&bcr.Command{
		Name:    "simple",
		Summary: "Create a new message with the given roles.",
		Usage:   "<channel> <name> <roles...>",
		Args:    bcr.MinArgs(3),

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.BoolP("mention", "m", false, "Show roles as @mentions")
			fs.StringP("description", "d", "", "A description to show before the role list")

			return fs
		},

		Command: b.simple,
	})

	simple.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add more roles to the given message.",
		Usage:   "<message link|ID> <roles...>",
		Args:    bcr.MinArgs(2),

		Command: b.simpleAdd,
	})

	simple.AddSubcommand(&bcr.Command{
		Name:    "update",
		Summary: "Replace the roles on the given message with the given roles.",
		Usage:   "<message link|ID> <roles...>",
		Args:    bcr.MinArgs(2),

		Command: b.simpleUpdate,
	})

	rr.AddSubcommand(&bcr.Command{
		Name:    "clear",
		Summary: "Clear react roles from the given message or the entire server.",
		Usage:   "[message]",

		Command: b.clear,
	})

	// add handlers
	bot.Router.AddHandler(b.reactionAdd)
	bot.Router.AddHandler(b.reactionRemove)

	// add cleanup handlers
	bot.Router.AddHandler(b.channelDelete)
	bot.Router.AddHandler(b.messageDelete)

	return
}

func (bot *Bot) channelDelete(ev *gateway.ChannelDeleteEvent) {
	ct, err := bot.DB.Pool.Exec(context.Background(), "delete from react_roles where channel_id = $1", ev.ID)
	if err != nil {
		bot.Sugar.Errorf("Error cleaning up reaction roles: %v", err)
	}
	if n := ct.RowsAffected(); n != 0 {
		bot.Sugar.Infof("Removed %v reaction role entries as the channel they were in was deleted.", n)
	}
}

func (bot *Bot) messageDelete(ev *gateway.MessageDeleteEvent) {
	ct, err := bot.DB.Pool.Exec(context.Background(), "delete from react_roles where message_id = $1", ev.ID)
	if err != nil {
		bot.Sugar.Errorf("Error cleaning up reaction roles: %v", err)
	}
	if n := ct.RowsAffected(); n != 0 {
		bot.Sugar.Infof("Removed %v reaction role entries as the message they were on was deleted.", n)
	}
}
