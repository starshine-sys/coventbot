package config

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

		CustomPermissions: bot.ModRole,
		Command:           b.prefixAdd,
	})

	prefix.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a prefix.",
		Usage:   "<prefix>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.ModRole,
		Command:           b.prefixRemove,
	})

	list = append(list, prefix)

	wl := b.Router.AddCommand(&bcr.Command{
		Name:        "watchlist",
		Aliases:     []string{"watch-list", "wl"},
		Summary:     "Show the users currently on the watchlist.",
		Description: "The server watchlist notifies you when a member on it joins your server. Intended to be used for potential problem members who aren't worth banning.",

		CustomPermissions: bot.ModRole,
		Command:           b.watchlist,
	})

	wl.AddSubcommand(&bcr.Command{
		Name:        "channel",
		Aliases:     []string{"notifications", "notifs"},
		Summary:     "Set the notification channel",
		Description: "Set the channel where alerts will be sent when a user on the watchlist joins your server.",
		Usage:       "<new channel>",
		Args:        bcr.MinArgs(1),

		CustomPermissions: bot.ModRole,
		Command:           b.watchlistChannel,
	})

	wl.AddSubcommand(b.Router.AliasMust("show", nil, []string{"watchlist"}, nil))

	wl.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a user to the watch list.",
		Usage:   "<user>",

		Args:              bcr.MinArgs(1),
		CustomPermissions: bot.ModRole,
		Command:           b.watchlistAdd,
	})

	wl.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a user from the watch list.",
		Usage:   "<user>",

		Args:              bcr.MinArgs(1),
		CustomPermissions: bot.ModRole,
		Command:           b.watchlistRemove,
	})

	wl.AddSubcommand(&bcr.Command{
		Name:    "reason",
		Summary: "Set the reason for a user on the watchlist.",
		Usage:   "[user ID] [reason]",

		Args:              bcr.MinArgs(1),
		CustomPermissions: bot.ModRole,
		Command:           b.watchlistReason,
	})

	sb := b.Router.AddCommand(&bcr.Command{
		Name:    "starboard",
		Summary: "View or change this server's starboard settings.",

		GuildOnly:         true,
		CustomPermissions: bot.ModRole,
		Command:           b.starboardSettings,
	})

	sb.AddSubcommand(&bcr.Command{
		Name:    "channel",
		Summary: "Change this server's starboard channel.",
		Usage:   "<new channel|-clear>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.ModRole,
		Command:           b.starboardSetChannel,
	})

	sb.AddSubcommand(&bcr.Command{
		Name:    "emoji",
		Summary: "Change this server's starboard emoji.",
		Usage:   "<new emoji>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.ModRole,
		Command:           b.starboardSetEmoji,
	})

	sb.AddSubcommand(&bcr.Command{
		Name:    "limit",
		Summary: "Change this server's starboard limit.",
		Usage:   "<new limit>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.ModRole,
		Command:           b.starboardSetLimit,
	})

	sb.AddSubcommand(&bcr.Command{
		Name:    "stats",
		Aliases: []string{"statistics", "leaderboard", "lb"},
		Summary: "Show starboard statistics for this server.",

		GuildOnly: true,
		Command:   b.starboardStats,
	})

	bl := sb.AddSubcommand(&bcr.Command{
		Name:    "blacklist",
		Aliases: []string{"block", "bl"},
		Summary: "View this server's starboard blacklist.",

		CustomPermissions: bot.ModRole,
		Command:           b.blacklist,
	})

	bl.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a channel to the starboard blacklist.",
		Usage:   "<channel>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.ModRole,
		Command:           b.blacklistAdd,
	})

	bl.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a channel from the starboard blacklist.",
		Usage:   "<channel>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.ModRole,
		Command:           b.blacklistRemove,
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

		GuildOnly:         true,
		CustomPermissions: bot.ModRole,
		Command:           b.addTrigger,
	})

	triggers.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a trigger",
		Usage:   "<message> <emoji>",
		Args:    bcr.MinArgs(2),

		GuildOnly:         true,
		CustomPermissions: bot.ModRole,
		Command:           b.delTrigger,
	})

	helper := b.Router.AddCommand(&bcr.Command{
		Name:    "helper-roles",
		Aliases: []string{"helperroles"},
		Summary: "View this server's helper roles.",

		CustomPermissions: bot.ModRole,
		Command:           b.helperRoles,
	})

	helper.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a helper role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.ModRole,
		Command:           b.helperAddRole,
	})

	helper.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a helper role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.ModRole,
		Command:           b.helperRemoveRole,
	})

	mod := b.Router.AddCommand(&bcr.Command{
		Name:    "mod-roles",
		Aliases: []string{"modroles"},
		Summary: "View this server's mod roles.",

		CustomPermissions: bot.AdminRole,
		Command:           b.modRoles,
	})

	mod.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a mod role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.AdminRole,
		Command:           b.modAddRole,
	})

	mod.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a mod role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.AdminRole,
		Command:           b.modRemoveRole,
	})

	admin := b.Router.AddCommand(&bcr.Command{
		Name:    "admin-roles",
		Aliases: []string{"adminroles"},
		Summary: "View this server's admin roles.",

		CustomPermissions: bot.AdminRole,
		Command:           b.adminRoles,
	})

	admin.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a admin role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.AdminRole,
		Command:           b.adminAddRole,
	})

	admin.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a admin role.",
		Usage:   "<role>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.AdminRole,
		Command:           b.adminRemoveRole,
	})

	// add trigger handler
	b.Router.AddHandler(b.triggerReactionAdd)
	b.Router.AddHandler(b.triggerMessageDelete)

	// add join handler
	b.Router.AddHandler(b.watchlistMemberAdd)

	return s, append(list, wl, sb)
}
