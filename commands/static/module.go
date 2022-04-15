package static

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/spf13/pflag"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
	"github.com/starshine-sys/tribble/commands/static/info"
)

// Bot ...
type Bot struct {
	*bot.Bot
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Utility commands"
	b := &Bot{
		Bot: bot,
	}

	bot.Add(info.Init)

	bot.Interactions.Command("bubble").Check(bot.RequireNode("bubble")).Exec(b.bubbleSlash)
	bot.Interactions.Command("pride").Check(bot.RequireNode("pride")).Exec(b.prideSlash)
	bot.Interactions.Command("sampa").Check(bot.RequireNode("sampa")).Exec(b.sampaSlash)
	bot.Interactions.Command("linkto").Check(bot.RequireNode("linkto")).Exec(b.linktoSlash)

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "addemoji",
		Aliases: []string{"addemote", "steal"},
		Summary: "Add an emoji.",
		Description: `Adds an emoji. Source is optional if a file is attached.
Source can be either a link to an emote, an existing emote, or a link to a message (with the ` + "`-msg`" + ` flag).

If a message link is given as input, and the message has multiple emotes in it, a menu will pop up allowing you to choose the specific emote.`,
		Usage: "<source> [name]",

		Permissions: discord.PermissionManageEmojisAndStickers,
		Command:     b.addEmoji,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "exportemotes",
		Aliases: []string{"export-emotes"},
		Summary: "Export this server's emotes to a zip file.",

		Permissions: discord.PermissionManageEmojisAndStickers,
		Command:     b.exportEmotes,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "bubble",
		Summary: "Bubble wrap!",

		Command: b.bubble,
		Flags: func(fs *pflag.FlagSet) *pflag.FlagSet {
			fs.IntP("size", "s", 10, "Size of the bubble wrap")
			fs.BoolP("prepop", "p", false, "Whether to pre-pop some bubbles")
			return fs
		},
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "linkto",
		Aliases: []string{"moveto"},
		Summary: "Move a conversation to another channel.",
		Usage:   "<channel> [topic]",
		Args:    bcr.MinArgs(1),

		Command: b.linkto,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "enlarge",
		Aliases: []string{"e"},
		Summary: "Enlarge a custom emoji.",
		Usage:   "<emoji>",
		Args:    bcr.MinArgs(1),

		Command: b.enlarge,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "poll",
		Summary: "Make a poll using an embed.",
		Usage:   "<question> <option 1> <option 2> [options...]",
		Args:    bcr.MinArgs(3),

		GuildOnly: true,
		Command:   b.poll,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "quickpoll",
		Aliases: []string{"qp"},
		Summary: "Make a poll on the originating message.",
		Usage:   "[--options/-o num]",

		Command: b.quickpoll,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "meow",
		Summary: "Send a random meowmoji.",

		Command: b.meow,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "roll",
		Aliases: []string{"dice"},
		Summary: "Roll dice, defaults to 1d20.",
		Usage:   "[int?]d?num",

		Command: b.roll,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "embedsource",
		Aliases: []string{"embed-source"},
		Summary: "Show the source for a message's embed(s).",
		Usage:   "<message link>",
		Args:    bcr.MinArgs(1),

		Command: b.embedSource,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "hello",
		Aliases: []string{"hi", "hey", "heya"},
		Summary: "Say hi!",
		Hidden:  true,

		Command: b.hello,
	}))

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:        "sampa",
		Aliases:     []string{"xsampa", "x-sampa"},
		Summary:     "Convert X-SAMPA to IPA.",
		Description: "Convert [X-SAMPA](https://en.wikipedia.org/wiki/X-SAMPA) to IPA.\nThe converted message can be deleted by the user by reacting :x: or 🗑️.",
		Usage:       "<X-SAMPA>",
		Args:        bcr.MinArgs(1),

		Command: b.sampa,
	}))

	pride := b.Router.AddCommand(&bcr.Command{
		Name:    "pride",
		Summary: "Add a pride flag circle to your profile picture!",
		Usage:   "<flag>",

		Command: b.pride,
	})

	bot.Router.AddHandler(b.sampaReaction)

	// delete ?sampa messages (and potentially other responses) over a month old
	sf := discord.NewSnowflake(time.Now().UTC().Add(-720 * time.Hour))
	ct, err := bot.DB.Pool.Exec(context.Background(), "delete from command_responses where message_id < $1", sf)
	if err != nil {
		bot.Sugar.Errorf("Error cleaning command responses: %v", err)
	}
	if ct.RowsAffected() != 0 {
		bot.Sugar.Infof("Deleted %v command response(s)!", ct.RowsAffected())
	}

	return s, append(list, pride)
}
