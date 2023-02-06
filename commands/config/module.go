// SPDX-License-Identifier: AGPL-3.0-only
package config

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
	s = "Configuration"

	b := &Bot{Bot: bot}

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "user-cfg",
		Aliases: []string{"userconf", "user-config", "userconfig", "usercfg", "user-conf"},
		Summary: "Show or edit your user settings.",
		Usage:   "[show|<key> <value>]",

		Command: b.userCfg,
	}))

	prefix := b.Router.AddCommand(&bcr.Command{
		Name:    "prefix",
		Aliases: []string{"prefixes"},
		Summary: "Show the server's current prefixes.",

		Command: b.prefix,
	})

	prefix.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a prefix.",
		Usage:   "<prefix>",
		Args:    bcr.MinArgs(1),

		Command: b.prefixAdd,
	})

	prefix.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a prefix.",
		Usage:   "<prefix>",
		Args:    bcr.MinArgs(1),

		Command: b.prefixRemove,
	})

	list = append(list, prefix)

	wl := b.Router.AddCommand(&bcr.Command{
		Name:        "watchlist",
		Aliases:     []string{"watch-list", "wl"},
		Summary:     "Show the users currently on the watchlist.",
		Description: "The server watchlist notifies you when a member on it joins your server. Intended to be used for potential problem members who aren't worth banning.",

		Command: b.watchlist,
	})

	wl.AddSubcommand(&bcr.Command{
		Name:        "channel",
		Aliases:     []string{"notifications", "notifs"},
		Summary:     "Set the notification channel",
		Description: "Set the channel where alerts will be sent when a user on the watchlist joins your server.",
		Usage:       "<new channel>",
		Args:        bcr.MinArgs(1),

		Command: b.watchlistChannel,
	})

	wl.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a user to the watch list.",
		Usage:   "<user>",

		Args:    bcr.MinArgs(1),
		Command: b.watchlistAdd,
	})

	wl.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a user from the watch list.",
		Usage:   "<user>",

		Args:    bcr.MinArgs(1),
		Command: b.watchlistRemove,
	})

	wl.AddSubcommand(&bcr.Command{
		Name:    "reason",
		Summary: "Set the reason for a user on the watchlist.",
		Usage:   "[user ID] [reason]",

		Args:    bcr.MinArgs(1),
		Command: b.watchlistReason,
	})

	sb := b.Router.AddCommand(&bcr.Command{
		Name:    "starboard",
		Summary: "View or change this server's starboard settings.",

		GuildOnly: true,
		Command:   b.starboardSettings,
	})

	sb.AddSubcommand(&bcr.Command{
		Name:    "channel",
		Summary: "Change this server's starboard channel.",
		Usage:   "<new channel|-clear>",
		Args:    bcr.MinArgs(1),

		Command: b.starboardSetChannel,
	})

	sb.AddSubcommand(&bcr.Command{
		Name:    "emoji",
		Summary: "Change this server's starboard emoji.",
		Usage:   "<new emoji>",
		Args:    bcr.MinArgs(1),

		Command: b.starboardSetEmoji,
	})

	sb.AddSubcommand(&bcr.Command{
		Name:    "limit",
		Summary: "Change this server's starboard limit.",
		Usage:   "<new limit>",
		Args:    bcr.MinArgs(1),

		Command: b.starboardSetLimit,
	})

	sb.AddSubcommand(&bcr.Command{
		Name:    "stats",
		Aliases: []string{"statistics", "leaderboard", "lb"},
		Summary: "Show starboard statistics for this server.",

		GuildOnly: true,
		Command:   b.starboardStats,
	})

	sb.AddSubcommand(&bcr.Command{
		Name:    "name",
		Aliases: []string{"username"},
		Summary: "Set this server's starboard username.",
		Usage:   "<name>",
		Args:    bcr.MinArgs(1),

		Command: b.starboardSetUsername,
	})

	sb.AddSubcommand(&bcr.Command{
		Name:    "avatar",
		Aliases: []string{"pfp"},
		Summary: "Set this server's starboard avatar.",
		Usage:   "<link>",
		Args:    bcr.MinArgs(1),

		Command: b.starboardSetAvatar,
	})

	bl := sb.AddSubcommand(&bcr.Command{
		Name:    "blacklist",
		Aliases: []string{"block", "bl"},
		Summary: "View this server's starboard blacklist.",

		Command: b.blacklist,
	})

	bl.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a channel to the starboard blacklist.",
		Usage:   "<channel>",
		Args:    bcr.MinArgs(1),

		Command: b.blacklistAdd,
	})

	bl.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a channel from the starboard blacklist.",
		Usage:   "<channel>",
		Args:    bcr.MinArgs(1),

		Command: b.blacklistRemove,
	})

	triggers := b.Router.AddCommand(&bcr.Command{
		Name:    "triggers",
		Summary: "Add or remove triggers (reactions that trigger commands)",

		GuildOnly: true,
		Command:   func(ctx *bcr.Context) (err error) { return ctx.Help([]string{"triggers"}) },
	})

	triggers.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a trigger",
		Usage:   "<message> <emoji> <command>",
		Args:    bcr.MinArgs(3),

		GuildOnly: true,
		Command:   b.addTrigger,
	})

	triggers.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a trigger",
		Usage:   "<message> <emoji>",
		Args:    bcr.MinArgs(2),

		GuildOnly: true,
		Command:   b.delTrigger,
	})

	list = append(list, b.permCommands()...)

	// add trigger handler
	b.Router.AddHandler(b.triggerReactionAdd)
	b.Router.AddHandler(b.triggerMessageDelete)

	// add join handler
	b.Router.AddHandler(b.watchlistMemberAdd)

	return s, append(list, wl, sb)
}

func (bot *Bot) permCommands() (cmds []*bcr.Command) {
	root := bot.Router.AddCommand(&bcr.Command{
		Name:    "permissions",
		Aliases: []string{"perms"},
		Summary: "Configure permissions.",
		Command: func(ctx *bcr.Context) error { return ctx.Help([]string{"perms"}) },
	})

	node := root.AddSubcommand(&bcr.Command{
		Name:    "node",
		Aliases: []string{"nodes"},
		Summary: "List and edit permission nodes.",
		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.BoolP("edited", "e", false, "Only show edited nodes.")
			return fs
		},
		Command: bot.showNodes,
	})

	node.AddSubcommand(&bcr.Command{
		Name:    "set",
		Usage:   "<node> <level>",
		Summary: "Set a permission node.",
		Args:    bcr.MinArgs(2),
		Command: bot.setNode,
	})

	node.AddSubcommand(&bcr.Command{
		Name:    "reset",
		Usage:   "<node>",
		Summary: "Reset a permission node to the default level.",
		Args:    bcr.MinArgs(1),
		Command: bot.resetNode,
	})

	mod := root.AddSubcommand(&bcr.Command{
		Name:    "moderator",
		Aliases: []string{"mod", "moderatorroles"},
		Summary: "View or manage this server's moderator roles.",
		Command: bot.moderatorRoles,
	})

	mod.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a moderator role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		Command: bot.moderatorAddRole,
	})

	mod.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a moderator role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		Command: bot.moderatorRemoveRole,
	})

	manager := root.AddSubcommand(&bcr.Command{
		Name:    "manager",
		Aliases: []string{"managerroles"},
		Summary: "View or manage this server's manager roles.",
		Command: bot.managerRoles,
	})

	manager.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a manager role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),
		Command: bot.managerAddRole,
	})

	manager.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a manager role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),
		Command: bot.managerRemoveRole,
	})

	admin := root.AddSubcommand(&bcr.Command{
		Name:    "admin",
		Aliases: []string{"adminroles"},
		Summary: "View or manage this server's admin roles.",
		Command: bot.adminRoles,
	})

	admin.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add an admin role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		Command: bot.adminAddRole,
	})

	admin.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove an admin role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		Command: bot.adminRemoveRole,
	})

	return append(cmds, root)
}
