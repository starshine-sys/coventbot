package info

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

// Bot ...
type Bot struct {
	*bot.Bot

	start time.Time
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Info commands"

	b := &Bot{
		Bot:   bot,
		start: time.Now().UTC(),
	}

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "ping",
		Summary: "Show the bot's latency.",

		Command: b.ping,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "about",
		Summary: "Show some info about the bot.",

		Command: b.about,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "invite",
		Summary: "Get an invite link for the bot.",

		Command: b.invite,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "colour",
		Aliases: []string{"color"},
		Summary: "Preview a colour.",
		Usage:   "<hex code>",
		Args:    bcr.MinArgs(1),

		Command: b.colour,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "userinfo",
		Aliases: []string{"i", "profile", "whois"},
		Summary: "Show information about a user or yourself.",
		Usage:   "[user]",

		Command: b.memberInfo,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "roleinfo",
		Aliases: []string{"ri"},
		Summary: "Show information about a role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		Command: b.roleInfo,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "serverinfo",
		Aliases: []string{"si", "guildinfo"},
		Summary: "Show information about the server.",

		GuildOnly: true,
		Command:   b.serverInfo,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "help",
		Summary: "Show info about the bot, or info about a specific command.",
		Usage:   "[command]",
		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.BoolP("all", "a", false, "Show all commands, not just the ones you have access to.")
			return fs
		},

		Command: b.help,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "avatar",
		Aliases: []string{"pfp", "a"},
		Summary: "Show a user's avatar.",
		Usage:   "[user]",

		Command: b.avatar,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "idtime",
		Aliases: []string{"snowflake"},
		Summary: "Get the timestamp for a Discord ID.",
		Usage:   "<ID>",
		Args:    bcr.MinArgs(1),

		Command: b.idtime,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "getinvite",
		Aliases: []string{"inviteinfo", "invite-info"},
		Summary: "Get basic information from an invite link.",
		Usage:   "<invite link/ID>",
		Args:    bcr.MinArgs(1),

		Command: b.inviteInfo,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:      "message",
		Aliases:   []string{"msg", "m"},
		Summary:   "Quote a message across channels.",
		Usage:     "<message ID|link>",
		Args:      bcr.MinArgs(1),
		GuildOnly: true,

		Command: b.message,
	}))

	b.Router.AddHandler(b.messageCreate)

	return
}
