package customcommands

import (
	"fmt"
	"strings"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) showOrAdd(ctx *bcr.Context) (err error) {
	if len(ctx.Args) == 0 {
		ccs, err := bot.DB.AllCustomCommands(ctx.Message.GuildID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		s := "Custom commands:\n```\n"
		for _, cc := range ccs {
			s += fmt.Sprintf("%03d. %v\n", cc.ID, cc.Name)
		}
		s += "\n```"

		return ctx.SendX(s)
	}

	if len(ctx.Message.Attachments) == 0 {
		cc, err := bot.DB.CustomCommand(ctx.Message.GuildID, strings.ToLower(ctx.Args[0]))
		if err != nil {
			return ctx.SendfX("No command with the name %v found.", bcr.AsCode(ctx.Args[0]))
		}

		return ctx.SendfX("`%v`: %v\n```lua\n%v\n```", cc.ID, bcr.AsCode(cc.Name), cc.Source)
	}

	return nil
}
