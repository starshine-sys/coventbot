package starboard

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
	s = "Starboard"

	b := &Bot{Bot: bot}

	sb := b.Router.AddCommand(&bcr.Command{
		Name:    "starboard",
		Summary: "View or change this server's starboard settings.",

		GuildOnly: true,
		Command:   b.settings,
	})

	sb.AddSubcommand(&bcr.Command{
		Name:    "channel",
		Summary: "Change this server's starboard channel.",
		Usage:   "<new channel|-clear>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.setChannel,
	})

	sb.AddSubcommand(&bcr.Command{
		Name:    "emoji",
		Summary: "Change this server's starboard emoji.",
		Usage:   "<new emoji>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.setEmoji,
	})

	sb.AddSubcommand(&bcr.Command{
		Name:    "limit",
		Summary: "Change this server's starboard limit.",
		Usage:   "<new limit>",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.setLimit,
	})

	b.State.AddHandler(b.MessageReactionAdd)
	b.State.AddHandler(b.MessageReactionDelete)
	b.State.AddHandler(b.MessageReactionRemoveEmoji)
	return s, append(list, sb)
}
