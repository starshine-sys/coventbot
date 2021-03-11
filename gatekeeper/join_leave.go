package gatekeeper

import (
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/google/uuid"
)

func (bot *Bot) memberAdd(m *gateway.GuildMemberAddEvent) {
	settings, err := bot.serverSettings(m.GuildID)
	if err != nil {
		bot.Sugar.Errorf("Error getting server settings: %v", err)
		return
	}

	// if the member role isn't set, return
	if !settings.MemberRole.IsValid() {
		return
	}

	err = bot.setPending(PendingUser{
		UserID:   m.User.ID,
		ServerID: m.GuildID,
		Key:      uuid.New(),
		Pending:  true,
	})
	if err != nil {
		bot.Sugar.Errorf("Error setting user as pending: %v", err)
	}
}

func (bot *Bot) memberLeave(m *gateway.GuildMemberRemoveEvent) {
	settings, err := bot.serverSettings(m.GuildID)
	if err != nil {
		bot.Sugar.Errorf("Error getting server settings: %v", err)
		return
	}

	// if the member role isn't set, return
	if !settings.MemberRole.IsValid() {
		return
	}

	err = bot.deletePending(m.GuildID, m.User.ID)
	if err != nil {
		bot.Sugar.Errorf("Error deleting user entry from gatekeeper: %v", err)
	}
}
