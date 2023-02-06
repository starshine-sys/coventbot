// SPDX-License-Identifier: AGPL-3.0-only
package reactroles

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) simpleAdd(ctx *bcr.Context) (err error) {
	msg, err := ctx.ParseMessage(ctx.Args[0])
	if err != nil || msg.GuildID != ctx.Message.GuildID {
		_, err = ctx.Replyc(bcr.ColourRed, "Couldn't find that message.")
		return
	}

	dbMsg, err := bot.message(msg.ID)
	if err != nil || !dbMsg.Title.Valid {
		_, err = ctx.Replyc(bcr.ColourRed, "That message isn't an existing message with reaction roles.")
		return
	}

	entries, err := bot.getEntries(msg.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}
	if len(entries) >= 20 {
		_, err = ctx.Replyc(bcr.ColourRed, "That message already has the maximum number of reaction roles.")
		return
	}

	rls, n := ctx.GreedyRoleParser(ctx.Args[1:])
	if n == 0 {
		_, err = ctx.Replyc(bcr.ColourRed, "Couldn't parse any roles from your input.")
		return
	} else if n != -1 {
		_, err = ctx.Replyc(bcr.ColourOrange, "Warning: I could only parse %v out of %v roles from your input.", n, len(ctx.Args[1:]))
	}

	if len(entries)+len(rls) > 20 {
		_, err = ctx.Replyc(bcr.ColourRed, "Adding these roles to that message would result in it exceeding the maximum number of reaction roles.")
		return
	}

	simpleEmotes := simpleEmotes[len(entries):]

	e := msg.Embeds[0]

	for i, r := range rls {
		if dbMsg.Mention.Bool {
			e.Description += fmt.Sprintf("\n<:emoji:%v> %v", simpleEmotes[i], r.Mention())
		} else {
			e.Description += fmt.Sprintf("\n<:emoji:%v> %v", simpleEmotes[i], r.Name)
		}
	}

	_, err = ctx.State.EditEmbeds(msg.ChannelID, msg.ID, e)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "I couldn't edit the message.")
		return
	}

	for i, r := range rls {
		_, err = bot.DB.Pool.Exec(context.Background(), `insert into react_role_entries
		(message_id, emote, role_id) values ($1, $2, $3) on conflict (message_id, emote) do update
		set role_id = $3`, msg.ID, simpleEmotes[i], r.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		emoji := discord.APIEmoji("emoji:" + simpleEmotes[i])

		err = ctx.State.React(msg.ChannelID, msg.ID, emoji)
		if err != nil {
			ctx.Send("I couldn't react to the message.")
			return
		}
	}

	_, err = ctx.Sendf("Done! Added %v react roles.\nNote that my highest role must be above the roles you added to this message for the reactions to work.", len(rls))
	return
}
