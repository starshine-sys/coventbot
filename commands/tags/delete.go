package tags

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) deleteTag(ctx *bcr.Context) (err error) {
	t, err := bot.DB.GetTag(ctx.Message.GuildID, ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("No tag with that name found.", nil)
		return
	}

	if t.CreatedBy != ctx.Author.ID && !bot.isModerator(ctx) {
		_, err = ctx.Send("You don't have permission to delete this tag.", nil)
		return
	}

	author := &discord.EmbedAuthor{
		Name: t.CreatedBy.String(),
	}
	u, err := ctx.State.User(t.CreatedBy)
	if err == nil {
		author = &discord.EmbedAuthor{
			Name: u.Username + "#" + u.Discriminator,
			Icon: u.AvatarURL(),
		}
	}

	m, err := ctx.Send("**Are you sure you want to delete this tag?** React with ✅ to confirm, ❌ to cancel.", &discord.Embed{
		Author:      author,
		Title:       fmt.Sprintf("``%v``", bcr.EscapeBackticks(t.Name)),
		Description: t.Response,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %s | Created", t.ID),
		},
		Timestamp: discord.NewTimestamp(t.CreatedAt),
		Color:     ctx.Router.EmbedColor,
	})

	yes, timeout := ctx.YesNoHandler(*m, ctx.Author.ID)
	if timeout {
		_, err = ctx.Sendf("Operation timed out.")
		return
	}
	if !yes {
		_, err = ctx.Sendf("Cancelled.")
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "delete from tags where id = $1 and server_id = $2", t.ID, ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Sendf("Error deleting tag: %v", err)
		return
	}

	_, err = ctx.Send("Tag deleted.", nil)
	return
}
