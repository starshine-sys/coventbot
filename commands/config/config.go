package config

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
	s = "Configuration"

	b := &Bot{Bot: bot}

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:        "prefix",
		Aliases:     []string{"prefixes"},
		Summary:     "Show the server's current prefixes, or set new prefixes (maximum of 10, space separated).",
		Description: "Use `-clear` to clear custom prefixes.",
		Usage:       "[new prefixes...]",

		Permissions: discord.PermissionManageGuild,

		Command: b.prefix,
	}))

	wl := b.Router.AddCommand(&bcr.Command{
		Name:        "watchlist",
		Aliases:     []string{"watch-list", "wl"},
		Summary:     "Show the users currently on the watchlist.",
		Description: "The server watchlist notifies you when a member on it joins your server. Intended to be used for potential problem members who aren't worth banning.",

		Permissions: discord.PermissionKickMembers,

		Command: b.watchlist,
	})

	wl.AddSubcommand(&bcr.Command{
		Name:        "channel",
		Aliases:     []string{"notifications", "notifs"},
		Summary:     "Set the notification channel",
		Description: "Set the channel where alerts will be sent when a user on the watchlist joins your server.",
		Usage:       "<new channel>",
		Args:        bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,

		Command: b.watchlistChannel,
	})

	wl.AddSubcommand(b.Router.AliasMust("show", nil, []string{"watchlist"}, nil))

	wl.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a user to the watch list.",
		Usage:   "<user>",

		Args:        bcr.MinArgs(1),
		Permissions: discord.PermissionKickMembers,
		Command:     b.watchlistAdd,
	})

	wl.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a user from the watch list.",
		Usage:   "<user>",

		Args:        bcr.MinArgs(1),
		Permissions: discord.PermissionKickMembers,
		Command:     b.watchlistRemove,
	})

	wl.AddSubcommand(&bcr.Command{
		Name:    "reason",
		Summary: "Set the reason for a user on the watchlist.",
		Usage:   "[user ID] [reason]",

		Args:        bcr.MinArgs(1),
		Permissions: discord.PermissionKickMembers,
		Command:     b.watchlistReason,
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

		Permissions: discord.PermissionManageGuild,
		Command:     b.starboardSetChannel,
	})

	sb.AddSubcommand(&bcr.Command{
		Name:    "emoji",
		Summary: "Change this server's starboard emoji.",
		Usage:   "<new emoji>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.starboardSetEmoji,
	})

	sb.AddSubcommand(&bcr.Command{
		Name:    "limit",
		Summary: "Change this server's starboard limit.",
		Usage:   "<new limit>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.starboardSetLimit,
	})

	bl := sb.AddSubcommand(&bcr.Command{
		Name:    "blacklist",
		Aliases: []string{"block", "bl"},
		Summary: "View this server's starboard blacklist.",

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklist,
	})

	bl.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a channel to the starboard blacklist.",
		Usage:   "<channel>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklistAdd,
	})

	bl.AddSubcommand(&bcr.Command{
		Name:    "remove",
		Summary: "Remove a channel from the starboard blacklist.",
		Usage:   "<channel>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.blacklistRemove,
	})

	triggers := b.Router.AddCommand(&bcr.Command{
		Name:    "triggers",
		Summary: "Add or remove triggers (reactions that trigger commands)",

		GuildOnly: true,
		Command:   func(ctx *bcr.Context) (err error) { return },
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

	// add trigger handler
	b.State.AddHandler(b.triggerReactionAdd)
	b.State.AddHandler(b.triggerMessageDelete)

	// add join handler
	b.State.AddHandler(b.watchlistMemberAdd)

	return s, append(list, wl, sb)
}
