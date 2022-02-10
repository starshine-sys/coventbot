package customcommands

import (
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

type Bot struct {
	*bot.Bot
}

// Init ...
func Init(b *bot.Bot) (s string, list []*bcr.Command) {
	s = "Custom commands"

	bot := &Bot{b}

	bot.Router.AddCommand(&bcr.Command{
		Name:    "cc",
		Summary: "Show or create a custom command",
		Usage:   "[name]",
		Command: bot.showOrAdd,
	})

	return
}
