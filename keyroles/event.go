package keyroles

import (
	"context"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) guildMemberUpdate(ev *gateway.GuildMemberUpdateEvent) {
	s, _ := bot.Router.StateFromGuildID(ev.GuildID)

	old, err := bot.Member(ev.GuildID, ev.User.ID)
	if err != nil {
		return
	}

	var logChannel discord.ChannelID
	err = bot.DB.Pool.QueryRow(context.Background(), "select keyrole_channel from servers where id = $1", ev.GuildID).Scan(&logChannel)
	if err != nil {
		bot.Sugar.Errorf("Error getting log channel: %v", err)
		return
	}

	if !logChannel.IsValid() {
		return
	}

	var addedRoles, removedRoles []discord.RoleID
	for _, oldRole := range old.RoleIDs {
		if !roleIn(ev.RoleIDs, oldRole) {
			removedRoles = append(removedRoles, oldRole)
		}
	}
	for _, newRole := range ev.RoleIDs {
		if !roleIn(old.RoleIDs, newRole) {
			addedRoles = append(addedRoles, newRole)
		}
	}

	var keyRoles []uint64
	err = bot.DB.Pool.QueryRow(context.Background(), "select keyroles from servers where id = $1", ev.GuildID).Scan(&keyRoles)
	if err != nil {
		bot.Sugar.Errorf("Error getting key roles: %v", err)
		return
	}

	var addedKeyRoles, removedKeyRoles []discord.RoleID
	for _, r := range addedRoles {
		for _, k := range keyRoles {
			if r == discord.RoleID(k) {
				addedKeyRoles = append(addedKeyRoles, r)
			}
		}
	}

	for _, r := range removedRoles {
		for _, k := range keyRoles {
			if r == discord.RoleID(k) {
				removedKeyRoles = append(removedKeyRoles, r)
			}
		}
	}

	if len(addedKeyRoles) == 0 && len(removedKeyRoles) == 0 {
		return
	}

	e := discord.Embed{
		Title: "Key roles added or removed",
		Color: bcr.ColourBlurple,

		Author: &discord.EmbedAuthor{
			Name: ev.User.Username + "#" + ev.User.Discriminator,
			Icon: ev.User.AvatarURL(),
		},

		Footer: &discord.EmbedFooter{
			Text: "User ID: " + ev.User.ID.String(),
		},

		Timestamp: discord.NowTimestamp(),
	}

	if len(addedKeyRoles) > 0 {
		buf := ""

		for _, r := range addedKeyRoles {
			buf += r.Mention() + ", "
		}

		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Added role(s)",
			Value: buf,
		})
	}

	if len(removedKeyRoles) > 0 {
		buf := ""

		for _, r := range removedKeyRoles {
			buf += r.Mention() + ", "
		}

		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Removed role(s)",
			Value: buf,
		})
	}

	// sleep for a bit for the audit log
	time.Sleep(time.Second)

	logs, err := s.AuditLog(ev.GuildID, api.AuditLogData{
		ActionType: discord.MemberRoleUpdate,
		Limit:      100,
	})
	if err == nil {
		for _, l := range logs.Entries {
			if discord.UserID(l.TargetID) == ev.User.ID && l.ID.Time().After(time.Now().Add(-10*time.Second)) {
				mod, err := s.User(l.UserID)
				if err == nil {
					e.Fields = append(e.Fields, discord.EmbedField{
						Name:  "Actor",
						Value: fmt.Sprintf("%v#%v (%v, %v)", mod.Username, mod.Discriminator, mod.Mention(), mod.ID),
					})
				} else {
					e.Fields = append(e.Fields, discord.EmbedField{
						Name:  "Actor",
						Value: fmt.Sprintf("%v (%v)", l.UserID.Mention(), l.UserID),
					})
				}

				break
			}
		}
	} else {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Actor",
			Value: "*Unknown*",
		})
	}

	_, err = s.SendEmbeds(logChannel, e)
	if err != nil {
		bot.Sugar.Errorf("Error sending key role log message: %v", err)
	}
}

func roleIn(s []discord.RoleID, id discord.RoleID) (exists bool) {
	for _, r := range s {
		if id == r {
			return true
		}
	}
	return false
}
