package moderation

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
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

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "members",
		Summary: "Show a filtered list of members.",
		Usage:   "--help",

		Permissions: discord.PermissionManageMessages,
		Command:     b.members,
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

		Permissions: discord.PermissionManageRoles,
		Command:     b.roleDump,
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

		Permissions: discord.PermissionManageMessages,

		Command: b.echo,
	})

	echo.AddSubcommand(&bcr.Command{
		Name:    "to",
		Summary: "Echo something to the specified channel.",
		Usage:   "<channel>",

		Permissions: discord.PermissionManageMessages,

		Command: b.echoTo,
	})

	embed := b.Router.AddCommand(&bcr.Command{
		Name:    "embed",
		Summary: "Send an embed to this channel.",
		Usage:   "<JSON>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageMessages,

		Command: b.embed,
	})

	embed.AddSubcommand(&bcr.Command{
		Name:    "to",
		Summary: "Send an embed to the given channel.",
		Usage:   "<channel> <JSON>",
		Args:    bcr.MinArgs(2),

		Permissions: discord.PermissionManageMessages,

		Command: b.embedTo,
	})

	embed.AddSubcommand(&bcr.Command{
		Name:    "edit",
		Summary: "Edit the given message.",
		Usage:   "<message> <JSON>",
		Args:    bcr.MinArgs(2),

		Permissions: discord.PermissionManageMessages,

		Command: b.editEmbed,
	})

	list = append(list, echo, embed)

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "makeinvite",
		Aliases: []string{"createinvite"},
		Summary: "Make an invite for the current channel, or the given channel.",
		Usage:   "[channel]",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionCreateInstantInvite,
		Command:     b.makeInvite,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "lockdown",
		Aliases: []string{"lock"},
		Summary: "Toggles a channel being locked, hiding it from the `@everyone` role.",
		Usage:   "[channel]",

		Permissions: discord.PermissionManageRoles,
		Command:     b.lockdown,
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

		Command: b.discordSlowmode,
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

		Command: b.cmdSetSlowmode,

		Permissions: discord.PermissionManageGuild,
	})

	slowmode.AddSubcommand(&bcr.Command{
		Name:    "reset",
		Summary: "Reset the slowmode for the given user.",
		Usage:   "<user> [channel]",
		Args:    bcr.MinArgs(1),

		Command: b.resetSlowmode,

		Permissions: discord.PermissionManageMessages,
	})

	slowmode.AddSubcommand(&bcr.Command{
		Name:    "role",
		Summary: "Set the role that will ignore slowmode. (All bots automatically ignore slowmode)",
		Usage:   "[role|--clear]",
		Args:    bcr.MinArgs(1),

		Command: b.slowmodeRole,

		Permissions: discord.PermissionManageGuild,
	})

	list = append(list, slowmode)

	bot.State.AddHandler(b.slowmodeMessage)

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

		Permissions: discord.PermissionManageMessages,
		Command:     b.warn,
	}))

	bot.State.AddHandler(b.channelbanOnJoin)

	_, modLogList := modlog.InitCommands(bot)

	list = append(list, modLogList...)
	return
}
