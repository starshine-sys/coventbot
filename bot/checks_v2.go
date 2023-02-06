// SPDX-License-Identifier: AGPL-3.0-only
package bot

import "github.com/starshine-sys/bcr/v2"

func (bot *Bot) NoDM(ctx *bcr.CommandContext) (err error) {
	if ctx.Guild == nil || ctx.Member == nil || ctx.Channel == nil {
		return bcr.NewCheckError[*bcr.CommandContext]("This command cannot be used in DMs.")
	}
	return nil
}
