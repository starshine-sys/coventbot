package levels

import (
	"context"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) blacklistRoles(ctx *bcr.Context) (err error) {
	rls, _ := ctx.GreedyRoleParser(ctx.Args)
	roleIDs := []uint64{}
	for _, r := range rls {
		roleIDs = append(roleIDs, uint64(r.ID))
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update server_levels set blocked_roles = $1 where id = $2", roleIDs, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("Blocked roles updated!", nil)
	return
}

func (bot *Bot) blacklistChannels(ctx *bcr.Context) (err error) {
	chs, _ := ctx.GreedyChannelParser(ctx.Args)
	channelIDs := []uint64{}
	for _, ch := range chs {
		if (ch.Type != discord.GuildText && ch.Type != discord.GuildNews) || ch.GuildID != ctx.Channel.GuildID {
			_, err = ctx.Sendf("Channel %v is not a valid channel to block.", ch.Mention())
			return
		}
		channelIDs = append(channelIDs, uint64(ch.ID))
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update server_levels set blocked_channels = $1 where id = $2", channelIDs, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("Blocked channels updated!", nil)
	return
}

func (bot *Bot) blacklistCategories(ctx *bcr.Context) (err error) {
	chs, _ := ctx.GreedyChannelParser(ctx.Args)
	channelIDs := []uint64{}
	for _, ch := range chs {
		if ch.Type != discord.GuildCategory || ch.GuildID != ctx.Channel.GuildID {
			_, err = ctx.Sendf("%v is not a valid category to block.", ch.Mention())
			return
		}
		channelIDs = append(channelIDs, uint64(ch.ID))
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update server_levels set blocked_categories = $1 where id = $2", channelIDs, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("Blocked categories updated!", nil)
	return
}
