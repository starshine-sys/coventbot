package quotes

import (
	"sync"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

// Bot ...
type Bot struct {
	*bot.Bot

	AESKey [32]byte

	mu map[discord.MessageID]*sync.Mutex
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Quotes"

	b := &Bot{
		Bot: bot,
		mu:  make(map[discord.MessageID]*sync.Mutex),
	}

	copy(b.AESKey[:], bot.Config.AESKey)

	b.Router.AddHandler(b.reactionAdd)

	cmd := b.Router.AddCommand(&bcr.Command{
		Name:    "quote",
		Summary: "Show a random quote, or a quote from a given user, or a quote by ID.",
		Usage:   "[quote ID|user ID]",

		GuildOnly: true,
		Command:   b.quote,
	})

	cmd.AddSubcommand(&bcr.Command{
		Name:    "delete",
		Aliases: []string{"yeet"},
		Summary: "Delete the quote with the given ID.",
		Usage:   "<quote ID>",
		Args:    bcr.MinArgs(1),

		GuildPermissions: discord.PermissionManageMessages,
		Command:          b.cmdQuoteDelete,
	})

	cmd.AddSubcommand(&bcr.Command{
		Name:    "leaderboard",
		Aliases: []string{"lb"},
		Summary: "Show this server's quote leaderboard, for users or channels.",
		Usage:   "<\"channel\" or \"user\">",
		Args:    bcr.MinArgs(1),

		GuildOnly: true,
		Command:   b.leaderboard,
	})

	quotes := b.Router.AddCommand(&bcr.Command{
		Name:    "quotes",
		Summary: "Show a list of quotes.",

		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.StringP("user", "u", "", "Filter by a user.")
			fs.StringP("channel", "c", "", "Filter by a channel.")

			fs.BoolP("sort-by-message", "m", false, "Sort by message ID.")
			fs.BoolP("reversed", "r", false, "Reverse sorting.")

			return fs
		},

		GuildOnly: true,
		Command:   b.list,
	})

	cmd.AddSubcommand(b.Router.AliasMust("list", nil, []string{"quotes"}, nil))
	quotes.AddSubcommand(b.Router.AliasMust("leaderboard", []string{"lb"}, []string{"quote", "leaderboard"}, nil))

	quotes.AddSubcommand(&bcr.Command{
		Name:    "toggle",
		Summary: "Enable or disable quotes for this server.",
		Usage:   "<on|off>",
		Args:    bcr.MinArgs(1),

		GuildPermissions: discord.PermissionManageGuild,
		Command:          b.toggle,
	})

	quotes.AddSubcommand(&bcr.Command{
		Name:    "messages",
		Summary: "Enable or disable quote messages for this server.",
		Usage:   "<on|off>",
		Args:    bcr.MinArgs(1),

		GuildPermissions: discord.PermissionManageGuild,
		Command:          b.toggleSuppressMessages,
	})

	return s, append(list, cmd, quotes)
}
