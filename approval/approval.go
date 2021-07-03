package approval

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

// Bot ...
type Bot struct {
	*bot.Bot
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Approval commands"

	b := &Bot{Bot: bot}

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "approve",
		Aliases: []string{"g"},
		Summary: "Approve the given member.",
		Usage:   "<member>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageRoles,
		Command:     b.approve,
	}))

	conf := b.Router.AddCommand(&bcr.Command{
		Name:    "approval",
		Summary: "Configure manual approval.",

		Command: func(ctx *bcr.Context) (err error) { return },
	})

	conf.AddSubcommand(&bcr.Command{
		Name:    "channel",
		Summary: "Configure the channel approval messages are sent in.",
		Usage:   "<new channel>",

		Permissions: discord.PermissionManageGuild,
		Command:     b.setChannel,
	})

	conf.AddSubcommand(&bcr.Command{
		Name:    "message",
		Summary: "Configure the message sent when a user is approved.",
		Usage:   "<new message>",

		Permissions: discord.PermissionManageGuild,
		Command:     b.setMessage,
	})

	conf.AddSubcommand(&bcr.Command{
		Name:    "add-roles",
		Summary: "Configure the roles added when a user is approved.",
		Usage:   "<roles|-clear>",

		Permissions: discord.PermissionManageGuild,
		Command:     b.setAddRoles,
	})

	conf.AddSubcommand(&bcr.Command{
		Name:    "remove-roles",
		Summary: "Configure the roles removed when a user is approved.",
		Usage:   "<roles|-clear>",

		Permissions: discord.PermissionManageGuild,
		Command:     b.setRemoveRoles,
	})

	return s, append(list, conf)
}
