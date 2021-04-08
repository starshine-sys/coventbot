package admin

import (
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

// Bot ...
type Bot struct {
	*bot.Bot
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Bot owner commands"

	b := &Bot{Bot: bot}

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "guild",
		Summary: "Show info for the given server ID.",
		Usage:   "<ID>",
		Args:    bcr.MinArgs(1),

		Hidden:    true,
		OwnerOnly: true,
		Command:   b.serverInfo,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "status",
		Summary: "Set the bot's status.",
		Usage:   "[new status]",

		Hidden:    true,
		OwnerOnly: true,
		Command:   b.status,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "activity",
		Summary: "Set the bot's activity.",
		Usage:   "[type] [new activity]",

		Hidden:    true,
		OwnerOnly: true,
		Command:   b.activity,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "dm",
		Summary: "DM the given user a text-only message.",
		Usage:   "<user> <message>",

		Hidden:    true,
		OwnerOnly: true,
		Command:   b.dm,
	}))

	b.State.AddHandler(b.ready)

	return
}
