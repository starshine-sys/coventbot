// SPDX-License-Identifier: AGPL-3.0-only
package moderation

import (
	"fmt"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize/english"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) purge(ctx *bcr.Context) (err error) {
	botPerms, _ := ctx.State.Permissions(ctx.Message.ChannelID, ctx.Router.Bot.ID)
	if !botPerms.Has(discord.PermissionManageMessages) {
		_, err = ctx.Replyc(bcr.ColourRed, "%v does not have permission to manage messages in this channel.", ctx.Router.Bot.Username)
		return
	}

	var i uint64 = 100
	if len(ctx.Args) > 0 {
		i, err = strconv.ParseUint(ctx.RawArgs, 10, 10)
		if err != nil {
			_, err = ctx.Replyc(bcr.ColourRed, "Couldn't parse your input as a number.")
			return
		}
	}
	if i > 100 || i == 0 {
		_, err = ctx.Replyc(bcr.ColourRed, "Number %v is out of maximum range: should be between 1 and 100.", i)
		return
	}

	msgs, err := ctx.State.Session.Messages(ctx.Message.ChannelID, uint(i))
	if err != nil {
		return bot.Report(ctx, err)
	}
	var deleteIDs []discord.MessageID

	for _, msg := range msgs {
		if !msg.Pinned {
			deleteIDs = append(deleteIDs, msg.ID)
		}
	}

	var msg *discord.Message
	err = ctx.State.DeleteMessages(ctx.Message.ChannelID, deleteIDs,
		api.AuditLogReason(fmt.Sprintf("Bulk message delete by %v (%v)", ctx.Author.Tag(), ctx.Author.ID)))
	if err != nil {
		msg, err = ctx.Replyc(bcr.ColourRed, "Couldn't delete messages. Sorry :(")
	} else {
		msg, err = ctx.Replyc(bcr.ColourGreen, "Deleted %v!", english.Plural(len(deleteIDs), "message", "messages"))
	}
	if err != nil {
		return
	}

	time.AfterFunc(10*time.Second, func() {
		ctx.State.DeleteMessage(ctx.Message.ChannelID, msg.ID, "")
	})
	return
}
