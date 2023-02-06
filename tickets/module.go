// SPDX-License-Identifier: AGPL-3.0-only
package tickets

import (
	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

// Bot ...
type Bot struct {
	*bot.Bot
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Tickets"

	b := &Bot{bot}

	tickets := bot.Router.AddCommand(&bcr.Command{
		Name:    "tickets",
		Aliases: []string{"ticket"},
		Summary: "Ticket commands.",

		GuildOnly: true,
		Command:   func(ctx *bcr.Context) (err error) { return },
	})

	tickets.AddSubcommand(&bcr.Command{
		Name:    "new",
		Aliases: []string{"open"},
		Summary: "Open a ticket.",
		Usage:   "<category> [user]",
		Args:    bcr.MinArgs(1),

		GuildOnly: true,
		Command:   b.new,
	})

	tickets.AddSubcommand(&bcr.Command{
		Name:    "delete",
		Aliases: []string{"close"},
		Summary: "Close and delete a ticket.",

		GuildOnly: true,
		Command:   b.delete,
	})

	tickets.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a member to a ticket.",
		Usage:   "<member>",
		Args:    bcr.MinArgs(1),

		GuildOnly: true,
		Command:   b.add,
	})

	tickets.AddSubcommand(&bcr.Command{
		Name:    "rename",
		Summary: "Rename a ticket channel.",
		Usage:   "<new name>",
		Args:    bcr.MinArgs(1),

		GuildOnly: true,
		Command:   b.rename,
	})

	tickets.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a member from a ticket.",
		Usage:   "<member>",
		Args:    bcr.MinArgs(1),

		GuildOnly: true,
		Command:   b.remove,
	})

	tickets.AddSubcommand(&bcr.Command{
		Name:    "list",
		Summary: "Show a list of ticket categories.",

		Command: b.list,
	})

	cfg := tickets.AddSubcommand(&bcr.Command{
		Name:        "config",
		Aliases:     []string{"cfg"},
		Summary:     "Configure ticket categories.",
		Description: "Configure ticket categories. If you provide an existing ticket category, the existing configuration will be replaced.",
		Usage:       "<category> <name> <log channel>",
		Args:        bcr.MinArgs(3),

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.IntP("limit", "l", -1, "Per-user ticket limit (empty or -1 to disable)")
			fs.UintP("count", "c", 0, "Number to start numbering tickets at")
			fs.BoolP("creator-close", "C", false, "Whether or not the creator can close the ticket")

			return fs
		},

		Command: b.cfg,
	})

	cfg.AddSubcommand(&bcr.Command{
		Name:        "mention",
		Summary:     "Set the mention for a category.",
		Description: "Set the mention for a category. `{mention}` will be replaced with the user's mention, `{channel}` will be replaced with a link to the channel, `{here}` and `{everyone}` will be replaced with @here and @everyone, respectively.",
		Usage:       "<category> <mention|-clear>",
		Args:        bcr.MinArgs(2),

		Command: b.mention,
	})

	cfg.AddSubcommand(&bcr.Command{
		Name:    "description",
		Summary: "Set the description for a category.",
		Usage:   "<category> <description|-clear>",
		Args:    bcr.MinArgs(2),

		Command: b.description,
	})

	// add leave/join handlers
	b.Router.AddHandler(b.guildMemberRemove)
	b.Router.AddHandler(b.guildMemberAdd)

	return s, append(list, tickets)
}
