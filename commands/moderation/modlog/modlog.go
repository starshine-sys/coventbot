package modlog

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

// ModLog can be created in two ways: either by bcr for commands, or other
type ModLog struct {
	*bot.Bot
}

// InitCommands ...
func InitCommands(bot *bot.Bot) (s string, list []*bcr.Command) {
	b := &ModLog{Bot: bot}

	s = "Moderation logging"

	cfg := bot.Router.AddCommand(&bcr.Command{
		Name:    "modlog",
		Aliases: []string{"modlogs", "mod-log", "mod-logs"},
		Summary: "Get the moderation log for the specified user.",
		Usage:   "<user>",
		Args:    bcr.MinArgs(1),

		GuildPermissions: discord.PermissionManageMessages,
		Command:          b.modlog,
	})

	cfg.AddSubcommand(&bcr.Command{
		Name:    "channel",
		Summary: "Set the moderation log channel.",
		Usage:   "<channel|-clear>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.setchannel,
	})

	cfg.AddSubcommand(&bcr.Command{
		Name:    "export",
		Summary: "Export a user's (or the entire server's) mod log.",
		Usage:   "[user]",

		Permissions: discord.PermissionManageGuild,
		Command:     b.export,
	})

	cfg.AddSubcommand(&bcr.Command{
		Name:    "import",
		Summary: "Import a mod log. *This will overwrite any existing entries.*",

		Permissions: discord.PermissionManageGuild,
		Command:     b.cmdImport,
	})

	return s, append(list, cfg)
}

// New creates a new ModLog
func New(bot *bot.Bot) *ModLog {
	return &ModLog{Bot: bot}
}
