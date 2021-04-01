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
		Summary: "Show info for the given server ID",
		Usage:   "<ID>",
		Args:    bcr.MinArgs(1),

		Hidden:    true,
		OwnerOnly: true,
		Command:   b.serverInfo,
	}))

	return
}
