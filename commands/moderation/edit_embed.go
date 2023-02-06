// SPDX-License-Identifier: AGPL-3.0-only
package moderation

import (
	"encoding/json"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) editEmbed(ctx *bcr.Context) (err error) {
	msg, err := ctx.ParseMessage(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("I could not find that message.")
		return
	}

	args := []byte(
		strings.TrimSpace(
			strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]),
		),
	)

	var e discord.Embed
	err = json.Unmarshal(args, &e)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Edit(msg, "", true, e)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("Message edited!")
	return
}
