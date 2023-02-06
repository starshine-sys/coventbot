// SPDX-License-Identifier: AGPL-3.0-only
package modlog

import (
	"context"
	"strconv"
	"strings"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/starshine-sys/bcr"
)

func (bot *ModLog) reason(ctx *bcr.Context) (err error) {
	id := 0

	if equalFoldAny(ctx.Args[0], "latest", "last", "l") {
		err = bot.DB.Pool.QueryRow(context.Background(), "select id from mod_log where server_id = $1 order by id desc limit 1", ctx.Guild.ID).Scan(&id)
		if err != nil {
			if errors.Cause(err) == pgx.ErrNoRows {
				ctx.Replyc(bcr.ColourRed, "This server has no mod log entries.")
				return
			}
			return bot.Report(ctx, err)
		}
	} else {
		id, err = strconv.Atoi(ctx.Args[0])
		if err != nil {
			_, err = ctx.Replyc(bcr.ColourRed, "Couldn't parse ``%v`` as a number.", bcr.EscapeBackticks(ctx.Args[0]))
			return
		}
	}

	exists := false
	err = bot.DB.Pool.QueryRow(context.Background(), "select exists(select * from mod_log where id = $1 and server_id = $2)", id, ctx.Guild.ID).Scan(&exists)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if !exists {
		_, err = ctx.Replyc(bcr.ColourRed, "There's no mod log with the ID %v.", id)
		return
	}

	reason := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))
	if reason == ctx.RawArgs {
		reason = strings.Join(ctx.Args[1:], " ")
	}

	var oldEntry, entry Entry
	// get original
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &oldEntry, "select * from mod_log where id = $1 and server_id = $2", id, ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	entryReason := reason
	if oldEntry.ActionType == "channelban" || oldEntry.ActionType == "unchannelban" {
		ch, _, _ := strings.Cut(oldEntry.Reason, ": ")
		entryReason = ch + ": " + reason
	}

	// update and get the new one
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &entry, "update mod_log set reason = $1, mod_id = $2 where id = $3 and server_id = $4 returning *", entryReason, ctx.Author.ID, id, ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if !entry.ChannelID.IsValid() || !entry.MessageID.IsValid() {
		_, err = ctx.Replyc(bcr.ColourOrange, "I updated the reason, but couldn't update the log message.")
		return
	}

	e := bot.Embed(ctx.State, &entry)

	msg, err := ctx.State.Message(entry.ChannelID, entry.MessageID)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourOrange, "I updated the reason, but couldn't update the log message.")
		return
	}

	_, err = ctx.State.EditMessageComplex(msg.ChannelID, msg.ID, api.EditMessageData{
		Content: option.NewNullableString(""),
		Embeds:  &[]discord.Embed{e},
	})
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourOrange, "I updated the reason, but couldn't update the log message.")
		return
	}

	// ignore errors on reacting, the message might be deleted already
	ctx.State.React(ctx.Message.ChannelID, ctx.Message.ID, "✅")

	return
}

func equalFoldAny(s string, options ...string) bool {
	for _, o := range options {
		if strings.EqualFold(s, o) {
			return true
		}
	}
	return false
}
