package chanmirror

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
	s = "Channel mirror"

	b := &Bot{bot}

	b.Router.AddCommand(&bcr.Command{
		Name:    "mirror",
		Summary: "Show a list of mirrored channels, or set a channel mirror.",
		Usage:   "[<source> <destination|--clear>]",

		Permissions: discord.PermissionManageWebhooks,
		Command:     b.set,
	})

	b.Router.AddHandler(b.messageCreate)
	b.Router.AddHandler(b.reactionAdd)

	return
}
