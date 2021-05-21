package bot

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

// Send an embedded message
func (bot *Bot) Send(ctx *bcr.Context, format string, args ...interface{}) (m *discord.Message, err error) {
	return ctx.SendEmbed(bcr.SED{
		Message: fmt.Sprintf(format, args...),
	})
}
