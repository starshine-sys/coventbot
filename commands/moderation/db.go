// SPDX-License-Identifier: AGPL-3.0-only
package moderation

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/georgysavva/scany/pgxscan"
)

// MuteRoles ...
type MuteRoles struct {
	MuteRole  discord.RoleID
	PauseRole discord.RoleID
}

func (bot *Bot) muteRoles(guildID discord.GuildID) (r MuteRoles, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &r, "select mute_role, pause_role from servers where id = $1", guildID)
	return
}

func (bot *Bot) setMuteRole(guildID discord.GuildID, r discord.RoleID) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set mute_role = $1 where id = $2", r, guildID)
	return
}

func (bot *Bot) setPauseRole(guildID discord.GuildID, r discord.RoleID) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set pause_role = $1 where id = $2", r, guildID)
	return
}

func (bot *Bot) muteRoleDelete(ev *gateway.GuildRoleDeleteEvent) {
	r, err := bot.muteRoles(ev.GuildID)
	if err != nil {
		return
	}

	if ev.RoleID == r.MuteRole {
		bot.setMuteRole(ev.GuildID, 0)
	}
	if ev.RoleID == r.PauseRole {
		bot.setPauseRole(ev.GuildID, 0)
	}
}

func (bot *Bot) mutemeMessage(guildID discord.GuildID) (msg string, err error) {
	err = bot.DB.Pool.QueryRow(context.Background(), "select muteme_message from servers where id = $1", guildID).Scan(&msg)
	return
}

func (bot *Bot) setMutemeMessage(guildID discord.GuildID, msg string) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set muteme_message = $1 where id = $2", msg, guildID)
	return
}
