package keyroles

import (
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"

	"github.com/starshine-sys/tribble/bot"
)

// Bot ...
type Bot struct {
	*bot.Bot

	members   map[key][]discord.RoleID
	membersMu sync.Mutex
}

type key struct {
	GuildID discord.GuildID
	UserID  discord.UserID
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Key roles"

	b := &Bot{
		Bot:     bot,
		members: map[key][]discord.RoleID{},
	}

	kr := bot.Router.AddCommand(&bcr.Command{
		Name:    "keyrole",
		Aliases: []string{"key-role", "keyroles", "key-roles"},
		Summary: "Show this server's key roles.",

		Permissions: discord.PermissionManageGuild,
		Command:     b.list,
	})

	kr.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a key role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.add,
	})

	kr.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Aliases: []string{"rm", "del"},
		Summary: "Remove a key role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.remove,
	})

	kr.AddSubcommand(&bcr.Command{
		Name:    "channel",
		Aliases: []string{"ch"},
		Summary: "Set the log channel for key roles.",
		Usage:   "<channel|-clear>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.channel,
	})

	b.State.AddHandler(b.guildMemberUpdate)
	b.State.AddHandler(b.requestGuildMembers)
	b.State.AddHandler(b.guildMemberChunk)

	return
}
