package tags

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

// Bot ...
type Bot struct {
	*bot.Bot
}

func (Bot) String() string {
	return "Tag moderator"
}

// Check ...
func (bot *Bot) Check(ctx bcr.Contexter) (b bool, err error) {
	var id discord.RoleID
	err = bot.DB.Pool.QueryRow(context.Background(), "select tag_mod_role from servers where id = $1", ctx.GetGuild().ID).Scan(&id)
	if err != nil {
		return false, err
	}

	if !id.IsValid() {
		return true, nil
	}

	if ctx.GetMember() == nil {
		return false, nil
	}

	for _, r := range ctx.GetMember().RoleIDs {
		if r == id {
			return true, nil
		}
	}

	return false, nil
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

		GuildOnly:         true,
		CustomPermissions: b,
		Command:           b.addTag,
	})

	tag.AddSubcommand(&bcr.Command{
		Name:    "edit",
		Summary: "Edit an existing tag.",
		Usage:   "<name>\n<response>",
		Args:    bcr.MinArgs(1),

		GuildOnly:         true,
		CustomPermissions: b,
		Command:           b.editTag,
	})

	tag.AddSubcommand(&bcr.Command{
		Name:    "delete",
		Summary: "Delete a tag.",
		Usage:   "<name>",
		Args:    bcr.MinArgs(1),

		GuildOnly:         true,
		CustomPermissions: b,
		Command:           b.deleteTag,
	})

	tag.AddSubcommand(&bcr.Command{
		Name:    "info",
		Summary: "Show info on a tag.",
		Usage:   "<name>",
		Args:    bcr.MinArgs(1),

		GuildOnly: true,
		Command:   b.info,
	})

	tag.AddSubcommand(&bcr.Command{
		Name:    "role",
		Aliases: []string{"moderator"},
		Summary: "Restrict creating, editing, and deleting tags to a single role.",
		Description: `Restrict creating, editing, and deleting tags to a single role. If this is not set, anyone will be able to create, edit, or delete tags.
		
Tag ownership is bypassed completely with this setting. If there is no mod role set, only the creator of a tag, and anyone with the manage server permission, can edit or delete it.`,
		Usage: "[new role|-clear]",

		GuildOnly:         true,
		CustomPermissions: bot.ModRole,
		Command:           b.tagModerator,
	})

	tag.AddSubcommand(b.Router.AliasMust("list", nil, []string{"tags"}, nil))

	return s, append(list, tag)
}
