package tags

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
	s = "Tags"

	b := &Bot{Bot: bot}

	list = append(list, b.Router.AddCommand(&bcr.Command{
		Name:    "tags",
		Summary: "Show a list of tags in the current server, or the given server (in DMs).",
		Usage:   "[server ID]",

		Command: b.list,
	}))

	tag := b.Router.AddCommand(&bcr.Command{
		Name:        "tag",
		Summary:     "Display a tag.",
		Description: "Display the given tag. If the invoking message replied to a message, the response will reply to that message too.",
		Usage:       "<tag>",

		GuildOnly: true,
		Command:   b.tag,
	})

	tag.AddSubcommand(&bcr.Command{
		Name:    "add",
		Summary: "Add a tag.",
		Usage:   "<name>\n<response>",
		Args:    bcr.MinArgs(1),

		GuildOnly:   true,
		Permissions: discord.PermissionManageMessages,

		Command: b.addTag,
	})

	tag.AddSubcommand(b.Router.AliasMust("list", nil, []string{"tags"}, nil))

	return s, append(list, tag)
}
