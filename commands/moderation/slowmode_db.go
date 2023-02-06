// SPDX-License-Identifier: AGPL-3.0-only
package moderation

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
)

func (bot *Bot) setSlowmode(guildID discord.GuildID, channelID discord.ChannelID, slowmode time.Duration) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), `insert into slowmode (server_id, channel_id, slowmode) values ($1, $2, $3)
on conflict (channel_id) do update
set slowmode = $3`, guildID, channelID, slowmode)
	return
}

func (bot *Bot) clearSlowmode(channelID discord.ChannelID) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), `delete from slowmode where channel_id = $1`, channelID)
	return
}

func (bot *Bot) setUserSlowmode(guildID discord.GuildID, channelID discord.ChannelID, userID discord.UserID, expires time.Time) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), `insert into user_slowmode (server_id, channel_id, user_id, expiry)
values ($1, $2, $3, $4)
on conflict (user_id, channel_id) do update
set expiry = $4`, guildID, channelID, userID, expires)
	return err
}

func (bot *Bot) resetUserChannel(channelID discord.ChannelID, userID discord.UserID) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), `delete from user_slowmode where channel_id = $1 and user_id = $2`, channelID, userID)
	return err
}

func (bot *Bot) resetUserGuild(guildID discord.GuildID, userID discord.UserID) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), `delete from user_slowmode where server_id = $1 and user_id = $2`, guildID, userID)
	return err
}

func (bot *Bot) userSlowmode(channelID discord.ChannelID, userID discord.UserID) (slowmode bool) {
	bot.DB.Pool.QueryRow(context.Background(), `select exists (select * from user_slowmode where channel_id = $1 and user_id = $2 and expiry >= $3)`, channelID, userID, time.Now().UTC()).Scan(&slowmode)
	return
}

func (bot *Bot) hasSlowmode(channelID discord.ChannelID) (hasSlowmode bool, duration time.Duration) {
	bot.DB.Pool.QueryRow(context.Background(), "select exists (select * from slowmode where channel_id = $1)", channelID).Scan(&hasSlowmode)

	if !hasSlowmode {
		return false, 0
	}

	bot.DB.Pool.QueryRow(context.Background(), "select slowmode from slowmode where channel_id = $1", channelID).Scan(&duration)
	return true, duration
}

func (bot *Bot) slowmodeIgnore(guildID discord.GuildID, roles []discord.RoleID) (ignore bool, err error) {
	r := []uint64{}
	for _, role := range roles {
		r = append(r, uint64(role))
	}

	err = bot.DB.Pool.QueryRow(context.Background(), "select slowmode_ignore_role = any($1) from servers where id = $2", r, guildID).Scan(&ignore)
	return ignore, err
}
