// SPDX-License-Identifier: AGPL-3.0-only
package reactroles

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func (bot *Bot) reactionAdd(ev *gateway.MessageReactionAddEvent) {
	if !ev.GuildID.IsValid() || ev.UserID == bot.Router.Bot.ID {
		return
	}

	entries, err := bot.getEntries(ev.MessageID)
	if err != nil {
		bot.Sugar.Errorf("Error: %v", err)
		return
	}

	var toAdd []discord.RoleID
	for _, e := range entries {
		if ev.Emoji.ID.String() == e.Emote || ev.Emoji.Name == e.Emote {
			toAdd = append(toAdd, e.RoleID)
			break
		}
	}

	s, _ := bot.Router.StateFromGuildID(ev.GuildID)

	for _, r := range toAdd {
		err = s.AddRole(ev.GuildID, ev.UserID, r, api.AddRoleData{
			AuditLogReason: "Reaction role add role",
		})
		if err != nil {
			bot.Sugar.Errorf("Error adding role %v to user: %v", r, err)
		}
	}
}

func (bot *Bot) reactionRemove(ev *gateway.MessageReactionRemoveEvent) {
	if !ev.GuildID.IsValid() || ev.UserID == bot.Router.Bot.ID {
		return
	}

	entries, err := bot.getEntries(ev.MessageID)
	if err != nil {
		bot.Sugar.Errorf("Error: %v", err)
		return
	}

	var toRemove []discord.RoleID
	for _, e := range entries {
		if ev.Emoji.ID.String() == e.Emote || ev.Emoji.Name == e.Emote {
			toRemove = append(toRemove, e.RoleID)
			break
		}
	}

	s, _ := bot.Router.StateFromGuildID(ev.GuildID)

	for _, r := range toRemove {
		err = s.RemoveRole(ev.GuildID, ev.UserID, r, "Reaction role remove role")
		if err != nil {
			bot.Sugar.Errorf("Error removing role %v from user: %v", r, err)
		}
	}
}
