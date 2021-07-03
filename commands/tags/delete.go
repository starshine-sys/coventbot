package tags

import (
	"context"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) deleteTag(ctx *bcr.Context) (err error) {
	t, err := bot.DB.GetTag(ctx.Message.GuildID, ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("No tag with that name found.")
		return
	}

	if t.CreatedBy != ctx.Author.ID && !bot.isModerator(ctx) {
		_, err = ctx.Send("You don't have permission to delete this tag.")
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

	yes, timeout := ctx.ConfirmButton(ctx.Author.ID, bcr.ConfirmData{
		Message: "**Are you sure you want to delete this tag?**",
		Embeds: []discord.Embed{{
			Author:      author,
			Title:       fmt.Sprintf("``%v``", bcr.EscapeBackticks(t.Name)),
			Description: t.Response,
			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("ID: %s | Created", t.ID),
			},
			Timestamp: discord.NewTimestamp(t.CreatedAt),
			Color:     ctx.Router.EmbedColor,
		}},
		YesPrompt: "Delete",
		YesStyle:  discord.DangerButton,

		Timeout: 5 * time.Minute,
	})
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
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("Tag deleted.")
	return
}
