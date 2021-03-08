package config

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/coventbot/bot"
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

	return s, append(list, wl)
}
