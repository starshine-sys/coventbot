package tickets

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) add(ctx *bcr.Context) (err error) {
	var isTicket bool
	err = bot.DB.Pool.QueryRow(context.Background(), "select exists(select * from tickets where channel_id = $1)", ctx.Channel.ID).Scan(&isTicket)
	if err != nil {
		return bot.Report(ctx, err)
	}
	if !isTicket {
		_, err = ctx.Sendf("This isn't a ticket channel.")
		return
	}

	if perms, _ := ctx.State.Permissions(ctx.Channel.ID, ctx.Author.ID); !perms.Has(discord.PermissionManageMessages) {
		_, err = ctx.Sendf("You're not allowed to use this command.")
		return
	}

	u, err := ctx.ParseMember(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("Couldn't find that user.", nil)
		return
	}

	err = bot.State.EditChannelPermission(ctx.Channel.ID, discord.Snowflake(u.User.ID), api.EditChannelPermissionData{
		Type:  discord.OverwriteMember,
		Allow: discord.PermissionViewChannel | discord.PermissionSendMessages | discord.PermissionReadMessageHistory | discord.PermissionAddReactions,
	})
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("", &discord.Embed{
		Description: fmt.Sprintf("Added %v to %v", u.Mention(), ctx.Channel.Mention()),
		Color:       ctx.Router.EmbedColor,
	})
	return
}

func (bot *Bot) remove(ctx *bcr.Context) (err error) {
	var isTicket bool
	err = bot.DB.Pool.QueryRow(context.Background(), "select exists(select * from tickets where channel_id = $1)", ctx.Channel.ID).Scan(&isTicket)
	if err != nil {
		return bot.Report(ctx, err)
	}
	if !isTicket {
		_, err = ctx.Sendf("This isn't a ticket channel.")
		return
	}

	if perms, _ := ctx.State.Permissions(ctx.Channel.ID, ctx.Author.ID); !perms.Has(discord.PermissionManageMessages) {
		_, err = ctx.Sendf("You're not allowed to use this command.")
		return
	}

	u, err := ctx.ParseMember(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("Couldn't find that user.", nil)
		return
	}

	err = bot.State.EditChannelPermission(ctx.Channel.ID, discord.Snowflake(u.User.ID), api.EditChannelPermissionData{
		Type: discord.OverwriteMember,
		Deny: discord.PermissionViewChannel | discord.PermissionSendMessages | discord.PermissionReadMessageHistory | discord.PermissionAddReactions,
	})
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("", &discord.Embed{
		Description: fmt.Sprintf("Removed %v from %v", u.Mention(), ctx.Channel.Mention()),
		Color:       ctx.Router.EmbedColor,
	})
	return
}
