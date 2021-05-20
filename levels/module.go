package levels

import (
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

// Bot ...
type Bot struct {
	*bot.Bot
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Levels"

	b := &Bot{bot}

	bot.State.AddHandler(b.messageCreate)

	lvl := bot.Router.AddCommand(&bcr.Command{
		Name:    "level",
		Aliases: []string{"lvl", "rank"},
		Summary: "Show your or another user's level.",
		Usage:   "[user]",

		GuildOnly: true,
		Command:   b.level,
	})

	lvl.AddSubcommand(&bcr.Command{
		Name:    "setxp",
		Aliases: []string{"setexp", "set-xp", "set-exp"},
		Summary: "Set the given user's XP.",
		Usage:   "<user> <new XP>",
		Args:    bcr.MinArgs(2),

		Permissions: discord.PermissionManageGuild,
		Command:     b.setXP,
	})

	lvl.AddSubcommand(&bcr.Command{
		Name:    "colour",
		Aliases: []string{"color"},
		Summary: "Set the colour used in your level embed.",
		Usage:   "[new colour|clear]",

		GuildOnly: true,
		Command:   b.colour,
	})

	cfg := lvl.AddSubcommand(&bcr.Command{
		Name:    "config",
		Aliases: []string{"cfg"},
		Summary: "Configure levels.",
		Usage:   "[key <new value>]",

		Permissions: discord.PermissionManageGuild,
		Command:     b.config,
	})

	cfg.AddSubcommand(&bcr.Command{
		Name:    "add-reward",
		Aliases: []string{"addreward"},
		Summary: "Add a level reward.",
		Usage:   "<lvl> <role>",
		Args:    bcr.MinArgs(2),

		Permissions: discord.PermissionManageRoles,
		Command:     b.cmdAddReward,
	})

	cfg.AddSubcommand(&bcr.Command{
		Name:    "del-reward",
		Aliases: []string{"delreward"},
		Summary: "Remove a level reward.",
		Usage:   "<lvl>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageRoles,
		Command:     b.cmdDelReward,
	})

	cfg.AddSubcommand(&bcr.Command{
		Name:    "block-channels",
		Aliases: []string{"blockchannels"},
		Summary: "Block the given channel(s) from levels. Leave clear to unblock all channels.",
		Usage:   "[channels...]",

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklistChannels,
	})

	cfg.AddSubcommand(&bcr.Command{
		Name:    "block-roles",
		Aliases: []string{"blockroles"},
		Summary: "Block the given role(s) from levels. Leave clear to unblock all roles.",
		Usage:   "[roles...]",

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklistRoles,
	})

	cfg.AddSubcommand(&bcr.Command{
		Name:    "block-categories",
		Aliases: []string{"blockcategories"},
		Summary: "Block the given category(s) from levels. Leave clear to unblock all categories.",
		Usage:   "[categories...]",

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklistCategories,
	})

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "leaderboard",
		Summary: "Show this server's leaderboard.",

		GuildOnly: true,
		Command:   b.leaderboard,
	}))

	nolevels := bot.Router.AddCommand(&bcr.Command{
		Name:    "nolevels",
		Aliases: []string{"nolevel"},
		Summary: "Manage the user blacklist for levels.",

		Permissions: discord.PermissionManageMessages,
		Command:     b.nolevelsList,
	})

	nolevels.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Blacklist a user.",
		Usage:   "<user> [time]",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageMessages,
		Command:     b.nolevelsAdd,
	})

	nolevels.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Unblacklist a user.",
		Usage:   "<user>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageMessages,
		Command:     b.nolevelsRemove,
	})

	var o sync.Once
	bot.State.AddHandler(func(_ *gateway.ReadyEvent) {
		o.Do(func() {
			go b.nolevelLoop()
		})
	})

	return s, append(list, lvl, nolevels)
}
