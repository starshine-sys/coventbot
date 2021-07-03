package tickets

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) add(ctx *bcr.Context) (err error) {
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

	u, err := ctx.ParseMember(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "Couldn't find that user.")
		return
	}

	var exists bool
	err = bot.DB.Pool.QueryRow(context.Background(), "select exists(select * from tickets where $1 = any(users) and channel_id = $2)", u.User.ID, ctx.Channel.ID).Scan(&exists)
	if err != nil {
		return bot.Report(ctx, err)
	}
	if exists {
		_, err = ctx.Replyc(bcr.ColourRed, "%v is already part of this ticket.", u.Mention())
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update tickets set users = array_append(users, $1) where channel_id = $2", u.User.ID, ctx.Channel.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	err = ctx.State.EditChannelPermission(ctx.Channel.ID, discord.Snowflake(u.User.ID), api.EditChannelPermissionData{
		Type:  discord.OverwriteMember,
		Allow: discord.PermissionViewChannel | discord.PermissionSendMessages | discord.PermissionReadMessageHistory | discord.PermissionAddReactions,
	})
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("", discord.Embed{
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
		_, err = ctx.Send("Couldn't find that user.")
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update tickets set users = array_remove(users, $1) where channel_id = $2", u.User.ID, ctx.Channel.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	err = ctx.State.EditChannelPermission(ctx.Channel.ID, discord.Snowflake(u.User.ID), api.EditChannelPermissionData{
		Type: discord.OverwriteMember,
		Deny: discord.PermissionViewChannel | discord.PermissionSendMessages | discord.PermissionReadMessageHistory | discord.PermissionAddReactions,
	})
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("", discord.Embed{
		Description: fmt.Sprintf("Removed %v from %v", u.Mention(), ctx.Channel.Mention()),
		Color:       ctx.Router.EmbedColor,
	})
	return
}

func (bot *Bot) guildMemberRemove(ev *gateway.GuildMemberRemoveEvent) {
	channels := []uint64{}

	err := bot.DB.Pool.QueryRow(context.Background(), "select array(select channel_id from tickets where $1 = any(users) and category_id = any(array(select category_id from ticket_categories where server_id = $2)))", ev.User.ID, ev.GuildID).Scan(&channels)
	if err != nil {
		bot.Sugar.Errorf("Error getting channels: %v", err)
	}

	s, _ := bot.Router.StateFromGuildID(ev.GuildID)

	for _, ch := range channels {
		_, err = s.Channel(discord.ChannelID(ch))
		if err != nil {
			continue
		}

		_, err = s.SendEmbeds(discord.ChannelID(ch), discord.Embed{
			Color:       bcr.ColourBlurple,
			Description: fmt.Sprintf("%v left.", ev.User.Mention()),
		})
		if err != nil {
			bot.Sugar.Errorf("Error sending message: %v", err)
		}
	}
}

func (bot *Bot) guildMemberAdd(ev *gateway.GuildMemberAddEvent) {
	channels := []uint64{}

	err := bot.DB.Pool.QueryRow(context.Background(), "select array(select channel_id from tickets where $1 = any(users) and category_id = any(array(select category_id from ticket_categories where server_id = $2)))", ev.User.ID, ev.GuildID).Scan(&channels)
	if err != nil {
		bot.Sugar.Errorf("Error getting channels: %v", err)
	}

	s, _ := bot.Router.StateFromGuildID(ev.GuildID)

	for _, ch := range channels {
		_, err = s.Channel(discord.ChannelID(ch))
		if err != nil {
			continue
		}

		err = s.EditChannelPermission(discord.ChannelID(ch), discord.Snowflake(ev.User.ID), api.EditChannelPermissionData{
			Type:  discord.OverwriteMember,
			Allow: discord.PermissionViewChannel | discord.PermissionSendMessages | discord.PermissionReadMessageHistory | discord.PermissionAddReactions,
		})
		if err != nil {
			bot.Sugar.Errorf("Error updating overwrites for %v: %v", ch, err)
			continue
		}

		_, err = s.SendEmbeds(discord.ChannelID(ch), discord.Embed{
			Description: fmt.Sprintf("Added %v after rejoin.", ev.Mention()),
			Color:       bcr.ColourBlurple,
		})
	}
}
