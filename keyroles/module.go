// SPDX-License-Identifier: AGPL-3.0-only
package keyroles

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
	s = "Key roles"

	b := &Bot{
		Bot: bot,
	}

	kr := bot.Router.AddCommand(&bcr.Command{
		Name:    "keyrole",
		Aliases: []string{"key-role", "keyroles", "key-roles"},
		Summary: "Show this server's key roles.",

		Command: b.list,
	})

	kr.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a key role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		Command: b.add,
	})

	kr.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Aliases: []string{"rm", "del"},
		Summary: "Remove a key role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		Command: b.remove,
	})

	kr.AddSubcommand(&bcr.Command{
		Name:    "channel",
		Aliases: []string{"ch"},
		Summary: "Set the log channel for key roles.",
		Usage:   "<channel|-clear>",
		Args:    bcr.MinArgs(1),

		Command: b.channel,
	})

	b.Router.AddHandler(b.guildMemberUpdate)

	return
}
