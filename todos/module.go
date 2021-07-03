package todos

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
	s = "Todo commands"

	b := &Bot{bot}

	todo := b.Router.AddCommand(&bcr.Command{
		Name:    "todo",
		Summary: "Set a personal todo.",
		Usage:   "<text>",
		Args:    bcr.MinArgs(1),

		Command: b.todo,
	})

	todo.AddSubcommand(&bcr.Command{
		Name:      "channel",
		Summary:   "Set your personal todo channel.",
		Usage:     "[channel|-clear]",
		GuildOnly: true,

		Command: b.channel,
	})

	todo.AddSubcommand(&bcr.Command{
		Name:    "list",
		Summary: "List your todos.",

		Command: b.list,
	})

	todo.AddSubcommand(&bcr.Command{
		Name:    "delete",
		Aliases: []string{"del", "remove", "rm"},
		Summary: "Remove a todo",
		Usage:   "<id>",
		Args:    bcr.MinArgs(1),

		Command: b.delete,
	})

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "complete",
		Summary: "Complete a todo.",
		Usage:   "<id>",
		Args:    bcr.MinArgs(1),

		Command: b.cmdComplete,
	}))

	b.Router.AddHandler(b.reactionAdd)

	return s, append(list, todo)
}
