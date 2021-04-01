package tags

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) info(ctx *bcr.Context) (err error) {
	t, err := bot.DB.GetTag(ctx.Message.GuildID, ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("No tag with that name found.", nil)
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

	_, err = ctx.Send("", &discord.Embed{
		Author:      author,
		Title:       fmt.Sprintf("``%v``", bcr.EscapeBackticks(t.Name)),
		Description: t.Response,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %s | Created", t.ID),
		},
		Timestamp: discord.NewTimestamp(t.CreatedAt),
		Color:     ctx.Router.EmbedColor,
	})
	return err
}
