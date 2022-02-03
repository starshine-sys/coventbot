package todos

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) channel(ctx *bcr.Context) (err error) {
	if len(ctx.Args) == 0 {
		id, err := bot.DB.UserIntGet(ctx.Author.ID, "todo_channel")
		if err != nil {
			return bot.Report(ctx, err)
		}

		chID := discord.ChannelID(id)

		if !chID.IsValid() {
			_, err = ctx.Reply("You don't currently have a todo channel set.")
		} else {
			_, err = ctx.Reply("Your todo channel is currently set to %v.", chID.Mention())
		}
		return err
	}

	if ctx.Guild.OwnerID != ctx.Author.ID {
		_, err = ctx.Replyc(bcr.ColourRed, "Your todo channel must be in a server you own.")
		return
	}

	var chID discord.ChannelID

	if ctx.RawArgs == "-clear" || ctx.RawArgs == "--clear" {
		chID = 0
	} else {
		ch, err := ctx.ParseChannel(ctx.RawArgs)
		if err != nil || ch.GuildID != ctx.Guild.ID || ch.Type != discord.GuildText {
			_, err = ctx.Replyc(bcr.ColourRed, "Channel not found!")
			return err
		}
		chID = ch.ID
	}

	err = bot.DB.UserIntSet(ctx.Author.ID, "todo_channel", int64(chID))
	if err != nil {
		return bot.Report(ctx, err)
	}

	if chID == 0 {
		_, err = ctx.Reply("Reset your todo channel!")
	} else {
		_, err = ctx.Reply("Set your todo channel to %v!", chID.Mention())
	}
	return
}
