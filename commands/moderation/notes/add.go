package notes

import (
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/db"
)

func (bot *Bot) addNote(ctx *bcr.Context) (err error) {
	perms := bot.globalPerms(ctx)
	if !perms.Has(discord.PermissionMoveMembers) && !perms.Has(discord.PermissionManageMessages) {
		_, err = ctx.Replyc(bcr.ColourRed, "You're not allowed to use this command.")
		return
	}

	u, err := ctx.ParseUser(ctx.Args[0])
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "User not found.")
		return
	}

	note := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))
	if note == ctx.RawArgs {
		note = strings.Join(ctx.Args[1:], " ")
	}

	if len(note) > 800 {
		_, err = ctx.Replyc(bcr.ColourRed, "Note too long (%v > 800 characters).", len(note))
	}

	n, err := bot.DB.NewNote(db.Note{
		ServerID:  ctx.Guild.ID,
		UserID:    u.ID,
		Note:      note,
		Moderator: ctx.Author.ID,
		Created:   time.Now().UTC(),
	})
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Replyc(bcr.ColourGreen, "✅ **Note taken.** (#%v)\n**Note:** %v", n.ID, n.Note)
	return
}
