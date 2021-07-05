// Package names tracks usernames and nicknames
package names

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
	s = "Name logging"

	b := &Bot{bot}

	b.Router.AddHandler(b.nicknameChange)
	b.Router.AddHandler(b.usernameChange)

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "nicknames",
		Summary: "Show the nickname history for a user.",
		Usage:   "[user]",

		CustomPermissions: bot.ModRole,
		Command:           b.nicknames,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "usernames",
		Summary: "Show the username history for a user.",
		Usage:   "[user]",

		CustomPermissions: bot.ModRole,
		Command:           b.usernames,
	}))

	return
}
