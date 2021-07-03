package todos

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) channel(ctx *bcr.Context) (err error) {
	if len(ctx.Args) == 0 {
		var chID discord.ChannelID
		bot.DB.Pool.QueryRow(context.Background(), "select todo_channel from user_config where user_id = $1", ctx.Author.ID).Scan(&chID)

		if chID == 0 {
			_, err = ctx.Reply("You don't currently have a todo channel set.")
		} else {
			_, err = ctx.Reply("Your todo channel is currently set to %v.", chID.Mention())
		}
		return
	}

	if ctx.Guild.OwnerID != ctx.Author.ID {
		_, err = ctx.Replyc(bcr.ColourRed, "Your todo channel must be in a server you own.")
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

	_, err = bot.DB.Pool.Exec(context.Background(), `insert into user_config (user_id, todo_channel)
	values ($1, $2) on conflict (user_id) do
	update set todo_channel = $2`, ctx.Author.ID, chID)
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

func (bot *Bot) getChannel(u discord.UserID) (ch discord.ChannelID) {
	bot.DB.Pool.QueryRow(context.Background(), "select todo_channel from user_config where user_id = $1", u).Scan(&ch)
	return
}
