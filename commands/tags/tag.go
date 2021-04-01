package tags

import (
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) tag(ctx *bcr.Context) (err error) {
	t, err := bot.DB.GetTag(ctx.Message.GuildID, ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("No tag with that name found.", nil)
		return
	}

	m := ctx.NewMessage().Content(t.Response).BlockMentions()
	// reply to the same message the invocation replied to
	if ctx.Message.Reference != nil {
		m.Reference(ctx.Message.Reference.MessageID)
	}
	_, err = m.Send()
	return err
}
