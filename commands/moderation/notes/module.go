package notes

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
	s = "Notes"

	b := &Bot{bot}

	b.Router.AddCommand(&bcr.Command{
		Name:    "notes",
		Summary: "List a user's notes.",
		Usage:   "<user>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.HelperRole,
		Command:           b.list,
	})

	b.Router.AddCommand(&bcr.Command{
		Name:    "setnote",
		Aliases: []string{"addnote"},
		Summary: "Add a note.",
		Usage:   "<user> <note>",
		Args:    bcr.MinArgs(2),

		CustomPermissions: bot.HelperRole,
		Command:           b.addNote,
	})

	b.Router.AddCommand(&bcr.Command{
		Name:    "delnote",
		Aliases: []string{"rmnote"},
		Summary: "Remove a note.",
		Usage:   "<note ID>",
		Args:    bcr.MinArgs(1),

		CustomPermissions: bot.ModRole,
		Command:           b.delNote,
	})

	b.Router.AddCommand(&bcr.Command{
		Name:    "bgc",
		Aliases: []string{"backgroundcheck"},
		Summary: "Show a background check for the given user.",
		Usage:   "[user]",

		CustomPermissions: bot.HelperRole,
		Command: func(ctx *bcr.Context) (err error) {
			if len(ctx.Args) == 0 {
				ctx.Args = []string{ctx.Author.ID.String()}
				ctx.RawArgs = ctx.Author.ID.String()
			}

			err = bot.Router.GetCommand("i").Command(ctx)
			if err != nil {
				return
			}
			err = bot.Router.GetCommand("notes").Command(ctx)
			if err != nil {
				return
			}
			return bot.Router.GetCommand("modlogs").Command(ctx)
		},
	})

	return
}
