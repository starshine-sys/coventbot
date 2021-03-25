// Package names tracks usernames and nicknames
package names

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
	s = "Name logging"

	b := &Bot{bot}

	b.State.AddHandler(b.nicknameChange)
	b.State.AddHandler(b.usernameChange)

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "nicknames",
		Summary: "Show the nickname history for a user.",
		Usage:   "[user]",

		Permissions: discord.PermissionManageRoles,
		Command:     b.nicknames,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "usernames",
		Summary: "Show the username history for a user.",
		Usage:   "[user]",

		Permissions: discord.PermissionManageRoles,
		Command:     b.usernames,
	}))

	return
}
