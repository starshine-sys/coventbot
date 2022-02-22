package approval

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
	s = "Approval commands"

	b := &Bot{Bot: bot}

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "approve",
		Aliases: []string{"g"},
		Summary: "Approve the given member.",
		Usage:   "<member>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.ModRole,
		Command:           b.approve,
	}))

	conf := b.Router.AddCommand(&bcr.Command{
		Name:              "approval",
		Summary:           "Configure manual approval.",
		CustomPermissions: bot.ModRole,

		Command: func(ctx *bcr.Context) (err error) { return ctx.Help([]string{"approval"}) },
	})

	conf.AddSubcommand(&bcr.Command{
		Name:    "channel",
		Summary: "Configure the channel approval messages are sent in.",
		Usage:   "<new channel>",

		CustomPermissions: bot.ModRole,
		Command:           b.setChannel,
	})

	conf.AddSubcommand(&bcr.Command{
		Name:    "message",
		Summary: "Configure the message sent when a user is approved.",
		Usage:   "<new message>",

		CustomPermissions: bot.ModRole,
		Command:           b.setMessage,
	})

	conf.AddSubcommand(&bcr.Command{
		Name:    "add-roles",
		Summary: "Configure the roles added when a user is approved.",
		Usage:   "<roles|-clear>",

		CustomPermissions: bot.ModRole,
		Command:           b.setAddRoles,
	})

	conf.AddSubcommand(&bcr.Command{
		Name:    "remove-roles",
		Summary: "Configure the roles removed when a user is approved.",
		Usage:   "<roles|-clear>",

		CustomPermissions: bot.ModRole,
		Command:           b.setRemoveRoles,
	})

	return s, append(list, conf)
}
