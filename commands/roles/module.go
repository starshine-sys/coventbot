package roles

import (
	"github.com/diamondburned/arikawa/v2/discord"
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

	b := &Bot{Bot: bot}

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "roles",
		Summary: "Show a list of self-assignable role categories.",
		Usage:   "[category name]",

		GuildOnly: true,
		Command:   b.listCategories,
	}))

	// We're doing some :sparkles: stupid shit :sparkles: here to have the command "role" available early for the "moderation" category
	role := b.Router.GetCommand("role")
	role.Summary = "Give yourself a role."
	role.Usage = "<role name>"
	role.Args = bcr.MinArgs(1)
	role.GuildOnly = true

	cfg := role.AddSubcommand(&bcr.Command{
		Name:    "cfg",
		Aliases: []string{"config"},
		Summary: "Role configuration commands.",

		Permissions: discord.PermissionManageRoles,
		Command:     func(*bcr.Context) (err error) { return },
	})

	cat := cfg.AddSubcommand(&bcr.Command{
		Name:    "category",
		Aliases: []string{"cat"},
		Summary: "Role category management.",

		Permissions: discord.PermissionManageRoles,
		Command:     func(*bcr.Context) (err error) { return },
	})

	cat.AddSubcommand(&bcr.Command{
		Name:    "add",
		Aliases: []string{"create"},
		Summary: "Create a role category.",
		Usage:   "<name>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageRoles,
		Command:     b.addCategory,
	})

	_ = cat

	return
}
