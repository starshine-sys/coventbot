package tickets

import (
	"context"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) rename(ctx *bcr.Context) (err error) {
	var isTicket bool
	err = bot.DB.Pool.QueryRow(context.Background(), "select exists(select * from tickets where channel_id = $1)", ctx.Channel.ID).Scan(&isTicket)
	if err != nil {
		return bot.Report(ctx, err)
	}
	if !isTicket {
		_, err = ctx.Replyc(bcr.ColourRed, "This isn't a ticket channel.")
		return
	}

	if perms, _ := ctx.State.Permissions(ctx.Channel.ID, ctx.Author.ID); !perms.Has(discord.PermissionManageMessages) {
		_, err = ctx.Replyc(bcr.ColourRed, "You're not allowed to use this command.")
		return
	}

	err = ctx.State.ModifyChannel(ctx.Channel.ID, api.ModifyChannelData{
		Name: ctx.RawArgs,
	})
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "I couldn't update the channel name.")
		return
	}

	// get new name, grab from API to not hit the (possibly out of date) cache
	ch, err := ctx.State.Session.Channel(ctx.Channel.ID)
	if err != nil {
		_, err = ctx.Reply("Renamed ticket to %v!", ctx.Channel.Mention())
		return
	}

	_, err = ctx.Reply("Renamed ticket to #%v!", ch.Name)
	return
}
