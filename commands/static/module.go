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
		Description: `Adds an emoji. Source is optional if a file is attached.
Source can be either a link to an emote, an existing emote, or a link to a message (with the ` + "`-msg`" + ` flag).

If a message link is given as input, and the message has multiple emotes in it, a menu will pop up allowing you to choose the specific emote.`,
		Usage: "<source> [name]",

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

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "meow",
		Summary: "Send a random meowmoji.",

		Command: b.meow,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "roll",
		Aliases: []string{"dice"},
		Summary: "Roll dice, defaults to 1d20.",
		Usage:   "[int?]d?num",

		Command: b.roll,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "embedsource",
		Aliases: []string{"embed-source"},
		Summary: "Show the source for a message's embed(s).",
		Usage:   "<message link>",
		Args:    bcr.MinArgs(1),

		Command: b.embedSource,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "hello",
		Aliases: []string{"hi", "hey", "heya"},
		Summary: "Say hi!",
		Hidden:  true,

		Command: b.hello,
	}))

	return s, list
}
