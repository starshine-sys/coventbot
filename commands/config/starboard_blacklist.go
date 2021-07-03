package config

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/db"
)

func (bot *Bot) blacklist(ctx *bcr.Context) (err error) {
	b, err := bot.DB.StarboardBlacklist(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	var x string
	// append all channel IDs (as mentions) to x
	for _, c := range b {
		x += fmt.Sprintf("<#%v>\n", c)
	}
	if len(b) == 0 {
		x = "No channels are blacklisted."
	}
	_, err = ctx.Send("", discord.Embed{
		Title:       "Starboard blacklist",
		Description: x,
		Color:       ctx.Router.EmbedColor,
	})
	return err
}

func (bot *Bot) blacklistRemove(ctx *bcr.Context) (err error) {
	if ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("You need to provide a channel.")
		return err
	}

	ch, err := ctx.ParseChannel(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("The channel you gave was not found.")
		return err
	}
	if ch.GuildID != ctx.Message.GuildID {
		_, err = ctx.Sendf("The given channel (%v) isn't in this server.", ch.Mention())
		return err
	}

	err = bot.DB.RemoveFromBlacklist(ctx.Message.GuildID, ch.ID)
	if err != nil {
		if err == db.ErrorNotBlacklisted {
			_, err = ctx.Send("That channel isn't blacklisted.")
			return err
		}

		return bot.Report(ctx, err)
	}

	_, err = ctx.Sendf("Removed %v from the blacklist.", ch.Mention())
	return
}

func (bot *Bot) blacklistAdd(ctx *bcr.Context) (err error) {
	if ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("You need to provide a channel.")
		return err
	}

	ch, err := ctx.ParseChannel(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("The channel you gave was not found.")
		return err
	}
	if ch.GuildID != ctx.Message.GuildID {
		_, err = ctx.Sendf("The given channel (%v) isn't in this server.", ch.Mention())
		return err
	}

	err = bot.DB.AddToBlacklist(ctx.Message.GuildID, ch.ID)
	if err != nil {
		if err == db.ErrorAlreadyBlacklisted {
			_, err = ctx.Send("That channel is already blacklisted.")
			return err
		}

		return bot.Report(ctx, err)
	}

	_, err = ctx.Sendf("Added %v to the blacklist.", ch.Mention())
	return
}
