// SPDX-License-Identifier: AGPL-3.0-only
package quotes

import (
	"context"
	"strings"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) toggle(ctx *bcr.Context) (err error) {
	var enable bool

	switch strings.ToLower(ctx.RawArgs) {
	case "true", "enable", "on":
		enable = true
	case "false", "disable", "off":
		enable = false
	}

	current := bot.quotesEnabled(ctx.Guild.ID)
	if current == enable {
		if enable {
			_, err = ctx.Send("Quotes are already enabled on this server.")
		} else {
			_, err = ctx.Send("Quotes are already disabled on this server.")
		}
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set quotes_enabled = $1 where id = $2", enable, ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if enable {
		_, err = ctx.Send("Enabled quotes for this server!")
	} else {
		_, err = ctx.Send("Disabled quotes for this server!")
	}
	return
}

func (bot *Bot) toggleSuppressMessages(ctx *bcr.Context) (err error) {
	var disable bool

	switch strings.ToLower(ctx.RawArgs) {
	case "true", "enable", "on":
		disable = false
	case "false", "disable", "off":
		disable = true
	}

	current := bot.suppressMessages(ctx.Guild.ID)
	if current == disable {
		if disable {
			_, err = ctx.Send("Quote messages are already disabled on this server.")
		} else {
			_, err = ctx.Send("Quote messages are already enabled on this server.")
		}
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set quote_suppress_messages = $1 where id = $2", disable, ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if disable {
		_, err = ctx.Send("Disabled quote messages for this server!")
	} else {
		_, err = ctx.Send("Enabled quote messages for this server!")
	}
	return
}
