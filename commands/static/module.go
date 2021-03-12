package static

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
	"github.com/starshine-sys/tribble/commands/static/info"
)

// Bot ...
type Bot struct {
	*bot.Bot
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Utility commands"
	b := &Bot{
		Bot: bot,
	}

	bot.Add(info.Init)

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "addemoji",
		Aliases: []string{"addemote", "steal"},
		Summary: "Add an emoji.",
		Usage:   "-h",

		Permissions: discord.PermissionManageEmojis,

		Command: b.addEmoji,
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
		Name:    "makeinvite",
		Aliases: []string{"createinvite"},
		Summary: "Make an invite for the current channel, or the given channel.",
		Usage:   "[channel]",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionCreateInstantInvite,
		Command:     b.makeInvite,
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
