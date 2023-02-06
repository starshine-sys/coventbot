// SPDX-License-Identifier: AGPL-3.0-only
package customcommands

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/customcommands/cc"
)

const maxScriptSize = 10 * 1024

func (bot *Bot) showOrAdd(ctx *bcr.Context) (err error) {
	allowCC := false
	for _, id := range bot.Config.AllowCCs {
		if id == ctx.Message.GuildID {
			allowCC = true
			break
		}
	}
	if !allowCC {
		return ctx.SendX("This guild does not have custom commands enabled.")
	}

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

	name := strings.ToLower(ctx.Args[0])

	if len(ctx.Message.Attachments) == 0 {
		cc, err := bot.DB.CustomCommand(ctx.Message.GuildID, name)
		if err != nil {
			return ctx.SendfX("No command with the name %v found.", bcr.AsCode(name))
		}

		return ctx.SendfX("`%v`: %v\n```lua\n%v\n```", cc.ID, bcr.AsCode(cc.Name), cc.Source)
	}

	if ctx.Message.Attachments[0].Size > maxScriptSize {
		return ctx.SendfX(":x: Maximum size of scripts is %v, this script is %v.", humanize.Bytes(ctx.Message.Attachments[0].Size), humanize.Bytes(maxScriptSize))
	}

	cctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(cctx, "GET", ctx.Message.Attachments[0].URL, nil)
	if err != nil {
		return bot.Report(ctx, err)
	}
	resp, err := bot.Client.Do(req)
	if err != nil {
		return bot.Report(ctx, err)
	}
	defer resp.Body.Close()

	src, err := io.ReadAll(resp.Body)
	if err != nil {
		return bot.Report(ctx, err)
	}

	s := cc.NewState(bot.Bot, ctx, nil)
	defer s.Close()

	err = s.Load(string(src))
	if err != nil {
		return ctx.SendfX("There was an error compiling the Lua code:\n```lua\n%v\n```", err)
	}

	cmd, err := bot.DB.SetCustomCommand(ctx.Message.GuildID, name, string(src))
	if err != nil {
		return bot.Report(ctx, err)
	}

	return ctx.SendfX("Saved command ``%v`` with ID `%v`!", bcr.EscapeBackticks(cmd.Name), cmd.ID)
}
