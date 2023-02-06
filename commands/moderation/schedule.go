// SPDX-License-Identifier: AGPL-3.0-only
package moderation

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	botpkg "github.com/starshine-sys/tribble/bot"
	"github.com/starshine-sys/tribble/commands/moderation/modlog"
)

type changeRoles struct {
	UserID  discord.UserID  `json:"user_id"`
	GuildID discord.GuildID `json:"guild_id"`

	AddRoles       []discord.RoleID `json:"add_roles"`
	RemoveRoles    []discord.RoleID `json:"remove_roles"`
	AuditLogReason string           `json:"reason"`

	SendModLog   bool           `json:"send_mod_log"`
	ModeratorID  discord.UserID `json:"moderator_id"`
	ModLogReason string         `json:"mod_log_reason"`
	ModLogType   string         `json:"mod_log_type"`
}

func (dat *changeRoles) Execute(ctx context.Context, id int64, bot *botpkg.Bot) error {
	bot.Sugar.Debugf("executing change role event %v, u:%v/g:%v", dat.UserID, dat.UserID, dat.GuildID)

	s, _ := bot.Router.StateFromGuildID(dat.GuildID)

	rls, err := s.Roles(dat.GuildID)
	if err != nil {
		bot.Sugar.Errorf("error getting roles for guild %v: %v", dat.GuildID, err)
		return botpkg.Reschedule
	}

	m, err := s.Member(dat.GuildID, dat.UserID)
	if err != nil {
		bot.Sugar.Errorf("error getting member %v in guild %v: %v", dat.UserID, dat.GuildID, err)
		// they most likely left, so don't reschedule
		return err
	}

	var toSetTo []discord.RoleID
	for _, id := range dat.AddRoles {
		if roleIn(rls, id) {
			toSetTo = append(toSetTo, id)
		}
	}

	for _, id := range m.RoleIDs {
		if !roleIDIn(dat.RemoveRoles, id) && roleIn(rls, id) {
			toSetTo = append(toSetTo, id)
		}
	}

	err = s.ModifyMember(dat.GuildID, dat.UserID, api.ModifyMemberData{
		Roles:          &toSetTo,
		AuditLogReason: api.AuditLogReason(dat.AuditLogReason),
	})
	if err != nil {
		bot.Sugar.Errorf("error updating member roles for %v in guild %v: %v", dat.UserID, dat.GuildID, err)
		// we probably don't have perms, don't reschedule
		return err
	}

	if dat.SendModLog {
		err = modlog.New(bot).Log(s, modlog.ActionType(dat.ModLogType), dat.GuildID, dat.UserID, dat.ModeratorID, dat.ModLogReason)
		if err != nil {
			bot.Sugar.Errorf("error sending mod log message: %v", err)
		}
	}

	return nil
}

func roleIDIn(ids []discord.RoleID, toCheck discord.RoleID) bool {
	for _, id := range ids {
		if id == toCheck {
			return true
		}
	}
	return false
}

func roleIn(roles []discord.Role, toCheck discord.RoleID) bool {
	for _, role := range roles {
		if role.ID == toCheck {
			return true
		}
	}
	return false
}

func (dat *changeRoles) Offset() time.Duration { return time.Minute }
