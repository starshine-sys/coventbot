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
		Name:    "help",
		Summary: "Show a list of commands, or info about a specific command.",
		Usage:   "[command]",

		Command: b.CommandList,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "addemoji",
		Aliases: []string{"addemote", "steal"},
		Summary: "Add an emoji",
		Usage:   "-h",

		Permissions: discord.PermissionManageEmojis,

		Command: b.addEmoji,
	}))

	return s, list
}
