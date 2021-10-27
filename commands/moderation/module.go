package moderation

import (
	"sync"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
	"github.com/starshine-sys/tribble/commands/moderation/mirror"
	"github.com/starshine-sys/tribble/commands/moderation/modlog"
)

// Bot ...
type Bot struct {
	*bot.Bot

	ModLog *modlog.ModLog
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Moderation commands"

	b := &Bot{
		Bot:    bot,
		ModLog: modlog.New(bot),
	}

	// add handler for importing other bots' mod logs
	mirror.Init(bot)

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "members",
		Summary: "Show a filtered list of members.",
		Usage:   "--help",

		CustomPermissions: bot.HelperRole,
		Command:           b.members,
	}))

	roles := b.Router.AddCommand(&bcr.Command{
		Name:    "role",
		Summary: "Role info and management commands.",

		Command: func(ctx *bcr.Context) (err error) {
			return nil
		},
	})

	roles.AddSubcommand(&bcr.Command{
		Name:    "dump",
		Summary: "Show a list of *all* roles with permissions and basic information.",

		CustomPermissions: bot.HelperRole,
		Command:           b.roleDump,
	})

	roles.AddSubcommand(&bcr.Command{
		Name:    "create",
		Aliases: []string{"add"},
		Summary: "Create a role.",
		Usage:   "<name> [colour] [-h: hoist] [-m: mentionable]",

		Permissions: discord.PermissionManageRoles,
		Command:     b.roleCreate,
	})

	roles.AddSubcommand(
		b.Router.AliasMust("info", nil, []string{"roleinfo"}, nil),
	)

	list = append(list, roles)

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "massban",
		Summary: "Ban all the given members with an optional reason.",
		Usage:   "<users...> [reason]",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionBanMembers,
		Command:     b.massban,
	}))

	echo := b.Router.AddCommand(&bcr.Command{
		Name:        "echo",
		Aliases:     []string{"say"},
		Summary:     "Make the bot say something.",
		Description: "To echo something into a different channel, use the `echo to` subcommand.",

		CustomPermissions: bot.ModRole,
		Command:           b.echo,
	})

	echo.AddSubcommand(&bcr.Command{
		Name:    "to",
		Summary: "Echo something to the specified channel.",
		Usage:   "<channel>",

		CustomPermissions: bot.ModRole,
		Command:           b.echoTo,
	})

	embed := b.Router.AddCommand(&bcr.Command{
		Name:    "embed",
		Summary: "Send an embed to this channel.",
		Usage:   "<JSON>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.ModRole,
		Command:           b.embed,
	})

	embed.AddSubcommand(&bcr.Command{
		Name:    "to",
		Summary: "Send an embed to the given channel.",
		Usage:   "<channel> <JSON>",
		Args:    bcr.MinArgs(2),

		CustomPermissions: bot.ModRole,
		Command:           b.embedTo,
	})

	embed.AddSubcommand(&bcr.Command{
		Name:    "edit",
		Summary: "Edit the given message.",
		Usage:   "<message> <JSON>",
		Args:    bcr.MinArgs(2),

		CustomPermissions: bot.ModRole,
		Command:           b.editEmbed,
	})

	list = append(list, echo, embed)

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "makeinvite",
		Aliases: []string{"createinvite"},
		Summary: "Make an invite for the current channel, or the given channel.",
		Usage:   "[channel]",

		CustomPermissions: bot.ModRole,
		Permissions:       discord.PermissionCreateInstantInvite,
		Command:           b.makeInvite,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "transcript",
		Summary: "Make a transcript of the given channel.",
		Usage:   "<channel>",
		Args:    bcr.MinArgs(1),

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.StringP("out", "o", "", "Channel to output the transcript to")
			fs.UintP("limit", "l", 500, "Number of messages to make a transcript of (maximum 2000)")
			fs.BoolP("json", "j", false, "Output as a JSON file")
			fs.BoolP("html", "h", false, "Output as a HTML file")
			return fs
		},

		CustomPermissions: bot.ModRole,
		Permissions:       discord.PermissionManageChannels,
		Command:           b.transcript,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "lockdown",
		Aliases: []string{"lock"},
		Summary: "Toggles a channel being locked, hiding it from the `@everyone` role.",
		Usage:   "[channel]",

		CustomPermissions: bot.ModRole,
		Command:           b.lockdown,
	}))

	slowmode := b.Router.AddCommand(&bcr.Command{
		Name:    "slowmode",
		Aliases: []string{"slow", "sm"},
		Summary: "Configure slowmode settings.",

		Command: func(ctx *bcr.Context) (err error) { return },
	})

	slowmode.AddSubcommand(&bcr.Command{
		Name:    "discord",
		Summary: "Set the given channel's Discord slowmode.",
		Usage:   "<duration> [channel]",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.ModRole,
		Command:           b.discordSlowmode,
	})

	slowmode.AddSubcommand(&bcr.Command{
		Name:    "set",
		Summary: "Set the slowmode for the given channel.",
		Usage:   "<channel> <duration>",
		Args:    bcr.MinArgs(1),

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.BoolP("clear", "c", false, "Clear the channel's slowmode.")

			return fs
		},

		Command:           b.cmdSetSlowmode,
		CustomPermissions: bot.ModRole,
	})

	slowmode.AddSubcommand(&bcr.Command{
		Name:    "reset",
		Summary: "Reset the slowmode for the given user.",
		Usage:   "<user> [channel]",
		Args:    bcr.MinArgs(1),

		Command:           b.resetSlowmode,
		CustomPermissions: bot.HelperRole,
	})

	slowmode.AddSubcommand(&bcr.Command{
		Name:    "role",
		Summary: "Set the role that will ignore slowmode. (All bots automatically ignore slowmode)",
		Usage:   "[role|--clear]",
		Args:    bcr.MinArgs(1),

		Command:           b.slowmodeRole,
		CustomPermissions: bot.ModRole,
	})

	list = append(list, slowmode)

	bot.Router.AddHandler(b.slowmodeMessage)

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "channelban",
		Summary: "Ban a member from using a channel.",
		Usage:   "[channel] <member>",
		Args:    bcr.MinArgs(1),

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.BoolP("full", "f", false, "Also hide the channel from the user.")

			return fs
		},

		Command: b.channelban,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "unchannelban",
		Summary: "Ban a member from using a channel.",
		Usage:   "[channel] <member>",
		Args:    bcr.MinArgs(1),

		Command: b.unchannelban,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "warn",
		Summary: "Warn a member.",
		Usage:   "<member> <reason>",
		Args:    bcr.MinArgs(2),

		CustomPermissions: bot.HelperRole,
		Command:           b.warn,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "ban",
		Summary: "Ban a user.",
		Usage:   "<user> [reason]",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionBanMembers,
		Command:     b.ban,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "unban",
		Summary: "Unban a user.",
		Usage:   "<user> [reason]",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionBanMembers,
		Command:     b.unban,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "muterole",
		Summary: "Show or set this server's mute role.",
		Usage:   "[role]",
		Args:    bcr.MinArgs(1),

		CustomPermissions: b.ModRole,
		Command:           b.muteRole,
	}))

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "pauserole",
		Summary: "Show or set this server's pause role.",
		Usage:   "[role]",
		Args:    bcr.MinArgs(1),

		CustomPermissions: b.ModRole,
		Command:           b.pauseRole,
	}))

	muteme := bot.Router.AddCommand(&bcr.Command{
		Name:    "muteme",
		Summary: "Mute yourself for the specified duration.",
		Usage:   "<duration>",
		Args:    bcr.MinArgs(1),

		GuildOnly: true,
		Command:   b.muteme,
	})

	muteme.AddSubcommand(&bcr.Command{
		Name:    "message",
		Summary: "Set the message used for the `muteme` command.",
		Description: `Available templates:
- {mention}: replaced with the user's @mention
- {tag}: replaced with the user's tag (username#0000)
- {duration}: replaced with the duration
- {action}: replaced with the action type`,
		Usage: "[new message|-clear]",

		CustomPermissions: b.ModRole,
		Command:           b.cmdMutemeMessage,
	})

	list = append(list, bot.Router.AddCommand(&bcr.Command{
		Name:    "pauseme",
		Summary: "Pause yourself for the specified duration.",
		Usage:   "<duration>",
		Args:    bcr.MinArgs(1),

		GuildOnly: true,
		Command:   b.pauseme,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "purge",
		Summary: "Bulk delete the given number of messages in the current channel. Ignores pinned messages.",
		Usage:   "[number, default 100]",

		Permissions: discord.PermissionManageMessages,
		Command:     b.purge,
	}))

	bot.Router.AddHandler(b.channelbanOnJoin)
	bot.Router.AddHandler(b.muteRoleDelete)
	bot.Router.AddHandler(b.muteOnJoin)

	_, modLogList := modlog.InitCommands(bot)

	list = append(list, modLogList...)

	state, _ := bot.Router.StateFromGuildID(0)

	var o sync.Once
	state.AddHandler(func(_ *gateway.ReadyEvent) {
		o.Do(func() {
			go b.doPendingActions(state)
		})
	})

	return
}
