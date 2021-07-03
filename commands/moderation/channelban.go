package moderation

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

var channelRe = regexp.MustCompile(`^<#\d{15,}>$`)

func (bot *Bot) getBool(sql string, params ...interface{}) (b bool, err error) {
	err = bot.DB.Pool.QueryRow(context.Background(), sql, params...).Scan(&b)
	return b, err
}

func (bot *Bot) channelban(ctx *bcr.Context) (err error) {
	go func() { ctx.State.Typing(ctx.Channel.ID) }()

	channel := ctx.Channel

	if channelRe.MatchString(ctx.Args[0]) {
		channel, err = ctx.ParseChannel(ctx.Args[0])
		if err != nil || channel.GuildID != ctx.Message.GuildID {
			_, err = ctx.Send("Channel not found.")
			return
		}

		if len(ctx.Args) > 1 {
			ctx.Args = ctx.Args[1:]
		} else {
			ctx.Args = []string{}
		}
	}

	if len(ctx.Args) < 1 {
		_, err = ctx.Send("You must give a member to channel ban.")
		return
	}

	member, err := ctx.ParseMember(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("Member not found.")
		return err
	}

	// get permissions for the invoking user
	if perms, _ := ctx.State.Permissions(channel.ID, ctx.Author.ID); !perms.Has(discord.PermissionManageMessages) {
		_, err = ctx.Send("You're not a mod, you can't do that!")
		return
	}

	// don't channelban ourselves
	if member.User.ID == ctx.Bot.ID {
		_, err = ctx.Send("No.")
		return
	}

	banned, err := bot.getBool("select exists (select * from channel_bans where server_id = $1 and channel_id = $2 and user_id = $3)", ctx.Message.GuildID, channel.ID, member.User.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}
	if banned {
		_, err = ctx.Send("User is already channel banned.")
		return
	}

	var allow, deny discord.Permissions
	for _, p := range channel.Permissions {
		if p.ID == discord.Snowflake(member.User.ID) {
			allow = p.Allow
			deny = p.Deny
		}
	}

	full, _ := ctx.Flags.GetBool("full")

	_, err = bot.DB.Pool.Exec(context.Background(), "insert into channel_bans (server_id, channel_id, user_id, full_ban) values ($1, $2, $3, $4)", ctx.Message.GuildID, channel.ID, member.User.ID, full)
	if err != nil {
		return bot.Report(ctx, err)
	}

	deny = deny | discord.PermissionSendMessages | discord.PermissionAddReactions

	if full {
		deny = deny | discord.PermissionViewChannel
	}

	err = ctx.State.EditChannelPermission(channel.ID, discord.Snowflake(member.User.ID), api.EditChannelPermissionData{
		Type:  discord.OverwriteMember,
		Allow: allow,
		Deny:  deny,
	})
	if err != nil {
		_, err = ctx.Sendf("I was unable to change the permissions in %v.", channel.Mention())
		return
	}

	// get reason
	var reason string
	if len(ctx.Args) > 1 {
		reason = strings.Join(ctx.Args[1:], " ")
	}

	err = bot.ModLog.Channelban(ctx, ctx.Message.GuildID, channel.ID, member.User.ID, ctx.Author.ID, reason)
	if err != nil {
		bot.Sugar.Errorf("Error logging channelban: %v", err)
	}

	_, err = ctx.Send("", discord.Embed{
		Color:       bcr.ColourBlurple,
		Description: fmt.Sprintf("Banned %v from %v", member.Mention(), channel.Mention()),
	})
	return
}

func (bot *Bot) unchannelban(ctx *bcr.Context) (err error) {
	go func() { ctx.State.Typing(ctx.Channel.ID) }()

	channel := ctx.Channel

	if channelRe.MatchString(ctx.Args[0]) {
		channel, err = ctx.ParseChannel(ctx.Args[0])
		if err != nil || channel.GuildID != ctx.Message.GuildID {
			_, err = ctx.Send("Channel not found.")
			return
		}

		if len(ctx.Args) > 1 {
			ctx.Args = ctx.Args[1:]
		} else {
			ctx.Args = []string{}
		}
	}

	if len(ctx.Args) < 1 {
		_, err = ctx.Send("You must give a member to unban.")
		return
	}

	member, err := ctx.ParseMember(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("Member not found.")
		return err
	}

	// get permissions for the invoking user
	if perms, _ := ctx.State.Permissions(channel.ID, ctx.Author.ID); !perms.Has(discord.PermissionManageMessages) {
		_, err = ctx.Send("You're not a mod, you can't do that!")
		return
	}

	banned, err := bot.getBool("select exists (select * from channel_bans where server_id = $1 and channel_id = $2 and user_id = $3)", ctx.Message.GuildID, channel.ID, member.User.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}
	if !banned {
		_, err = ctx.Send("User is not banned from that channel.")
		return
	}

	var allow, deny discord.Permissions
	for _, p := range channel.Permissions {
		if p.ID == discord.Snowflake(member.User.ID) {
			allow = p.Allow
			deny = p.Deny
		}
	}
	if deny != 0 && deny.Has(discord.PermissionViewChannel) {
		deny = deny ^ discord.PermissionViewChannel
	}
	if deny != 0 && deny.Has(discord.PermissionSendMessages) {
		deny = deny ^ discord.PermissionSendMessages
	}
	if deny != 0 && deny.Has(discord.PermissionAddReactions) {
		deny = deny ^ discord.PermissionAddReactions
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "delete from channel_bans where server_id = $1 and channel_id = $2 and user_id = $3", ctx.Message.GuildID, channel.ID, member.User.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	err = ctx.State.EditChannelPermission(channel.ID, discord.Snowflake(member.User.ID), api.EditChannelPermissionData{
		Type:  discord.OverwriteMember,
		Allow: allow,
		Deny:  deny,
	})
	if err != nil {
		_, err = ctx.Sendf("I was unable to change the permissions in %v.", channel.Mention())
		return
	}

	// get reason
	var reason string
	if len(ctx.Args) > 1 {
		reason = strings.Join(ctx.Args[1:], " ")
	}

	err = bot.ModLog.Unchannelban(ctx, ctx.Message.GuildID, channel.ID, member.User.ID, ctx.Author.ID, reason)
	if err != nil {
		bot.Sugar.Errorf("Error logging unchannelban: %v", err)
	}

	_, err = ctx.Send("", discord.Embed{
		Color:       bcr.ColourBlurple,
		Description: fmt.Sprintf("Unbanned %v from %v", member.Mention(), channel.Mention()),
	})
	return
}

func (bot *Bot) channelbanOnJoin(m *gateway.GuildMemberAddEvent) {
	channels := []struct {
		ChannelID discord.ChannelID
		FullBan   bool
	}{}

	err := pgxscan.Select(context.Background(), bot.DB.Pool, &channels, "select channel_id, full_ban from channel_bans where server_id = $1 and user_id = $2", m.GuildID, m.User.ID)
	if err != nil {
		bot.Sugar.Errorf("Error getting channel bans for %v: %v", m.User.ID, err)
		return
	}

	s, _ := bot.Router.StateFromGuildID(m.GuildID)

	for _, ch := range channels {
		deny := discord.PermissionSendMessages | discord.PermissionAddReactions
		if ch.FullBan {
			deny |= discord.PermissionViewChannel
		}

		err = s.EditChannelPermission(ch.ChannelID, discord.Snowflake(m.User.ID), api.EditChannelPermissionData{
			Type: discord.OverwriteMember,
			Deny: deny,
		})
		if err != nil {
			bot.Sugar.Errorf("Error setting channel ban for %v: %v", ch.ChannelID, err)
		}
	}

	return
}
