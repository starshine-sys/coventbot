package static

import (
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/coventbot/bot"
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
		Name:    "enlarge",
		Aliases: []string{"e"},
		Summary: "Enlarge a custom emoji.",
		Usage:   "<emoji>",
		Args:    bcr.MinArgs(1),

		Command: b.enlarge,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "invite",
		Summary: "Get an invite link for the bot.",

		Command: b.invite,
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
