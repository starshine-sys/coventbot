package config

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/db"
)

func (bot *Bot) watchlist(ctx *bcr.Context) (err error) {
	b, err := bot.DB.Watchlist(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
	}

	var x string
	// append all user IDs (as mentions) to x
	for _, c := range b {
		x += fmt.Sprintf("<@%v> (%v)\n", c, c)
	}
	if len(b) == 0 {
		x = "No users are on the watchlist."
	}
	_, err = ctx.Send("", &discord.Embed{
		Title:       "Watchlist",
		Description: x,
		Color:       ctx.Router.EmbedColor,
	})
	return err
}

func (bot *Bot) watchlistRemove(ctx *bcr.Context) (err error) {
	if ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("You need to provide a channel.", nil)
		return err
	}

	u, err := ctx.ParseUser(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("The user you gave was not found.", nil)
		return err
	}

	err = bot.DB.RemoveFromWatchlist(ctx.Message.GuildID, u.ID)
	if err != nil {
		if err == db.ErrorNotBlacklisted {
			_, err = ctx.Send("That user isn't on the watchlist.", nil)
			return err
		}

		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	_, err = ctx.Sendf("Removed %v / **%v#%v** from the watchlist.", u.Mention(), u.Username, u.Discriminator)
	return
}

func (bot *Bot) watchlistAdd(ctx *bcr.Context) (err error) {
	if ctx.CheckMinArgs(1); err != nil {
		_, err = ctx.Send("You need to provide a user.", nil)
		return err
	}

	u, err := ctx.ParseUser(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("The user you gave was not found.", nil)
		return err
	}

	err = bot.DB.AddToWatchlist(ctx.Message.GuildID, u.ID)
	if err != nil {
		if err == db.ErrorAlreadyBlacklisted {
			_, err = ctx.Send("That user is already on the watchlist.", nil)
			return err
		}

		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	_, err = ctx.Sendf("Added %v / **%v#%v** to the watchlist.", u.Mention(), u.Username, u.Discriminator)
	return
}
