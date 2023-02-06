// SPDX-License-Identifier: AGPL-3.0-only
package moderation

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) resetSlowmode(ctx *bcr.Context) (err error) {
	u, err := ctx.ParseMember(ctx.Args[0])
	if err != nil {
		_, err = ctx.Sendf("Couldn't find that user.")
		return
	}

	if len(ctx.Args) < 2 {
		yes, timeout := ctx.ConfirmButton(ctx.Author.ID, bcr.ConfirmData{
			Embeds: []discord.Embed{{
				Description: fmt.Sprintf("Are you sure that you want to reset %v's slowmode for **the entire server**?", u.Mention()),
				Color:       bcr.ColourBlurple,
			}},
			YesPrompt: "Yes",
			YesStyle:  discord.DangerButtonStyle(),
		})
		if !yes || timeout {
			_, err = ctx.Send(":x: Cancelled.")
			return err
		}

		err = bot.resetUserGuild(ctx.Message.GuildID, u.User.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		_, err = ctx.NewMessage().Content(fmt.Sprintf("Reset slowmode for %v!", u.Mention())).BlockMentions().Send()
		return err
	}

	ch, err := ctx.ParseChannel(ctx.Args[1])
	if err != nil {
		_, err = ctx.Send("Couldn't find that channel.")
		return
	}
	if ch.GuildID != ctx.Channel.GuildID {
		_, err = ctx.Send("That channel isn't in this server.")
		return
	}

	err = bot.resetUserChannel(ch.ID, u.User.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.NewMessage().Content(fmt.Sprintf("Reset slowmode in %v for %v!", ch.Mention(), u.Mention())).BlockMentions().Send()
	return
}
