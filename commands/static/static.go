package static

import (
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
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
	s = "Utility commands"
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
		Name:    "info",
		Aliases: []string{"i", "userinfo", "profile", "whois"},
		Summary: "Show information about a user or yourself.",
		Usage:   "[user]",

		Command: b.memberInfo,
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
		Summary: "Show a list of commands, or info about a specific command.",
		Usage:   "[command]",

		Command: b.CommandList,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "addemoji",
		Aliases: []string{"addemote", "steal"},
		Summary: "Add an emoji.",
		Usage:   "-h",

		Permissions: discord.PermissionManageEmojis,

		Command: b.addEmoji,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "avatar",
		Aliases: []string{"pfp", "a"},
		Summary: "Show a user's avatar.",
		Usage:   "[user]",

		Command: b.avatar,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "bubble",
		Summary: "Bubble wrap!",
		Usage:   "[-prepop] [-size 1-13]",

		Command: b.bubble,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "enlarge",
		Aliases: []string{"e"},
		Summary: "Enlarge a custom emoji.",
		Usage:   "<emoji>",
		Args:    bcr.MinArgs(1),

		Command: b.enlarge,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "idtime",
		Aliases: []string{"snowflake"},
		Summary: "Get the timestamp for a Discord ID.",
		Usage:   "<IDs...>",
		Args:    bcr.MinArgs(1),

		Command: b.idtime,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "makeinvite",
		Aliases: []string{"createinvite"},
		Summary: "Make an invite for the current channel, or the given channel.",
		Usage:   "[channel]",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionCreateInstantInvite,
		Command:     b.makeInvite,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "invite",
		Summary: "Get an invite link for the bot.",

		Command: b.invite,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "poll",
		Summary: "Make a poll using an embed.",
		Usage:   "<question> <option 1> <option 2> [options...]",
		Args:    bcr.MinArgs(3),

		GuildOnly: true,
		Command:   b.poll,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "quickpoll",
		Aliases: []string{"qp"},
		Summary: "Make a poll on the originating message.",
		Usage:   "[--options/-o num]",

		Command: b.quickpoll,
	}))

	echo := b.Router.AddCommand(&bcr.Command{
		Name:        "echo",
		Aliases:     []string{"say"},
		Summary:     "Make the bot say something.",
		Description: "To echo something into a different channel, use the `echo to` subcommand.",

		Permissions: discord.PermissionManageMessages,

		Command: b.echo,
	})

	echo.AddSubcommand(&bcr.Command{
		Name:    "to",
		Summary: "Echo something to the specified channel.",
		Usage:   "<channel>",

		Permissions: discord.PermissionManageMessages,

		Command: b.echoTo,
	})

	list = append(list, echo)
	return s, list
}
