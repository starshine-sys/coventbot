package levels

import (
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
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
		Name:    "setlvl",
		Aliases: []string{"setlevel", "set-lvl", "set-level"},
		Summary: "Set the given user's XP to the minimum needed for the given level.",
		Usage:   "<user> <new level>",
		Args:    bcr.MinArgs(2),

		Permissions: discord.PermissionManageGuild,
		Command:     b.setlvl,
	})

	lvl.AddSubcommand(&bcr.Command{
		Name:    "colour",
		Aliases: []string{"color"},
		Summary: "Set the colour used in your level card.",
		Usage:   "[new colour|clear]",

		GuildOnly: true,
		Command:   b.colour,
	})

	bg := lvl.AddSubcommand(&bcr.Command{
		Name:        "background",
		Aliases:     []string{"bg"},
		Summary:     "Set the background used in your level card to the attached image.",
		Description: "The image will automatically be resized to be 1200 pixels wide, with the original aspect ratio preserved.\nOnly the top 400 pixels (when resized) will be shown as part of the background.",
		Usage:       "[clear]",

		GuildOnly: true,
		Command:   b.background,
	})

	bg.AddSubcommand(&bcr.Command{
		Name:    "server",
		Summary: "Set this server's default level background.",
		Usage:   "[clear]",

		Permissions: discord.PermissionManageGuild,
		Command:     b.serverBackground,
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

	bl := cfg.AddSubcommand(&bcr.Command{
		Name:    "blacklist",
		Aliases: []string{"bl"},
		Summary: "Configure this server's blacklists.",

		Permissions: discord.PermissionManageGuild,
		Command:     func(ctx *bcr.Context) error { return nil },
	})

	channels := bl.AddSubcommand(&bcr.Command{
		Name:    "channels",
		Aliases: []string{"channel", "ch"},
		Summary: "Show the channel blacklist.",

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklistChannels,
	})

	channels.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a channel to the blacklist.",
		Usage:   "<channel>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklistChannelAdd,
	})

	channels.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a channel from the blacklist.",
		Usage:   "<channel>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklistChannelRemove,
	})

	roles := bl.AddSubcommand(&bcr.Command{
		Name:    "roles",
		Aliases: []string{"role"},
		Summary: "Show the role blacklist.",

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklistRoles,
	})

	roles.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a role to the blacklist.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklistRoleAdd,
	})

	roles.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a role from the blacklist.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklistRoleRemove,
	})

	cfg.AddSubcommand(&bcr.Command{
		Name:    "block-categories",
		Aliases: []string{"blockcategories"},
		Summary: "Block the given category(s) from levels. Leave clear to unblock all categories.",
		Usage:   "[categories...]",

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklistCategories,
	})

	categories := bl.AddSubcommand(&bcr.Command{
		Name:    "categories",
		Aliases: []string{"category", "cat"},
		Summary: "Show the category blacklist.",

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklistCategories,
	})

	categories.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a category to the blacklist.",
		Usage:   "<category>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklistCategoryAdd,
	})

	categories.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a category from the blacklist.",
		Usage:   "<category>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklistCategoryRemove,
	})

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "leaderboard",
		Aliases: []string{"lb"},
		Summary: "Show this server's leaderboard.",

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.BoolP("full", "f", false, "Show the full leaderboard, including people who left the server.")

			return fs
		},

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
