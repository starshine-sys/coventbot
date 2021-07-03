package levels

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) blacklistRoles(ctx *bcr.Context) (err error) {
	sc, err := bot.getGuildConfig(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(sc.BlockedRoles) == 0 {
		_, err = ctx.Send("There are no blacklisted roles.")
		return
	}

	var rls []string
	for _, r := range sc.BlockedRoles {
		rls = append(rls, discord.RoleID(r).Mention()+"\n")
	}

	_, err = bot.PagedEmbed(ctx,
		bcr.StringPaginator("Blacklisted roles", bcr.ColourBlurple, rls, 20), 10*time.Minute,
	)
	return
}

func (bot *Bot) blacklistRoleAdd(ctx *bcr.Context) (err error) {
	r, err := ctx.ParseRole(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("Role not found.")
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update server_levels set blocked_roles = array_append(blocked_roles, $1) where id = $2", r.ID, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Added " + r.Mention() + " to the blacklist.")
	return
}

func (bot *Bot) blacklistRoleRemove(ctx *bcr.Context) (err error) {
	r, err := ctx.ParseRole(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("Role not found.")
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update server_levels set blocked_roles = array_remove(blocked_roles, $1) where id = $2", r.ID, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Removed " + r.Mention() + " from the blacklist.")
	return
}

func (bot *Bot) blacklistChannels(ctx *bcr.Context) (err error) {
	sc, err := bot.getGuildConfig(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(sc.BlockedChannels) == 0 {
		_, err = ctx.Send("There are no blacklisted channels.")
		return
	}

	var chs []string
	for _, ch := range sc.BlockedChannels {
		chs = append(chs, discord.ChannelID(ch).Mention()+"\n")
	}

	_, err = bot.PagedEmbed(ctx,
		bcr.StringPaginator("Blacklisted channels", bcr.ColourBlurple, chs, 20), 10*time.Minute,
	)
	return
}

func (bot *Bot) blacklistChannelAdd(ctx *bcr.Context) (err error) {
	r, err := ctx.ParseChannel(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("Channel not found.")
		return
	}
	if (r.Type != discord.GuildText && r.Type != discord.GuildNews) || r.GuildID != ctx.Channel.GuildID {
		_, err = ctx.Send("Channel not found.")
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update server_levels set blocked_channels = array_append(blocked_channels, $1) where id = $2", r.ID, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Added " + r.Mention() + " to the blacklist.")
	return
}

func (bot *Bot) blacklistChannelRemove(ctx *bcr.Context) (err error) {
	r, err := ctx.ParseChannel(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("Channel not found.")
		return
	}
	if (r.Type != discord.GuildText && r.Type != discord.GuildNews) || r.GuildID != ctx.Channel.GuildID {
		_, err = ctx.Send("Channel not found.")
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update server_levels set blocked_channels = array_remove(blocked_channels, $1) where id = $2", r.ID, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Removed " + r.Mention() + " from the blacklist.")
	return
}

func (bot *Bot) blacklistCategories(ctx *bcr.Context) (err error) {
	sc, err := bot.getGuildConfig(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(sc.BlockedCategories) == 0 {
		_, err = ctx.Send("There are no blacklisted categories.")
		return
	}

	var chs []string
	for _, ch := range sc.BlockedCategories {
		chs = append(chs, discord.ChannelID(ch).Mention()+"\n")
	}

	_, err = bot.PagedEmbed(ctx,
		bcr.StringPaginator("Blacklisted categories", bcr.ColourBlurple, chs, 20), 10*time.Minute,
	)
	return
}

func (bot *Bot) blacklistCategoryAdd(ctx *bcr.Context) (err error) {
	r, err := ctx.ParseChannel(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("Category not found.")
		return
	}
	if r.Type != discord.GuildCategory || r.GuildID != ctx.Channel.GuildID {
		_, err = ctx.Send("Category not found.")
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update server_levels set blocked_categories = array_append(blocked_categories, $1) where id = $2", r.ID, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Added " + r.Name + " to the blacklist.")
	return
}

func (bot *Bot) blacklistCategoryRemove(ctx *bcr.Context) (err error) {
	r, err := ctx.ParseChannel(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("Category not found.")
		return
	}
	if r.Type != discord.GuildCategory || r.GuildID != ctx.Channel.GuildID {
		_, err = ctx.Send("Category not found.")
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update server_levels set blocked_categories = array_remove(blocked_categories, $1) where id = $2", r.ID, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Removed " + r.Name + " from the blacklist.")
	return
}
