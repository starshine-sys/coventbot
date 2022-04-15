package roles

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
	s = "Role commands"

	b := &Bot{bot}

	roles := b.Router.AddCommand(&bcr.Command{
		Name:    "roles",
		Summary: "Show a list of role categories, or their roles.",
		Usage:   "[category]",

		GuildOnly: true,
		Command:   b.categories,
	})

	cfg := roles.AddSubcommand(&bcr.Command{
		Name:    "config",
		Aliases: []string{"cfg"},
		Summary: "Create or update a role category.",
		Usage:   "<name>",
		Args:    bcr.MinArgs(1),

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.StringP("require-role", "r", "", "Require role")
			fs.StringP("desc", "d", "", "Description (max 1000 characters)")
			fs.StringP("colour", "c", "", "Category colour")

			return fs
		},

		Command: b.config,
	})

	cfg.AddSubcommand(&bcr.Command{
		Name:    "delete",
		Summary: "Delete a role category.",
		Usage:   "<id>",
		Args:    bcr.MinArgs(1),

		Command: b.delete,
	})

	cfg.AddSubcommand(&bcr.Command{
		Name:    "set-roles",
		Aliases: []string{"setroles"},
		Summary: "Set the given category's roles.",
		Usage:   "<id> <roles...>",
		Args:    bcr.MinArgs(2),

		Command: b.setRoles,
	})

	role := bot.Router.GetCommand("role")
	if role != nil {
		role.Summary = "Assign a role to yourself."
		role.Usage = "<role>"
		role.Args = bcr.MinArgs(1)
		role.GuildOnly = true
		role.Command = b.addRole

		list = append(list, role)
	}

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "derole",
		Summary: "Remove a role from yourself.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		GuildOnly: true,
		Command:   b.removeRole,
	}))

	return s, append(list, roles)
}
