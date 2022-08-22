package gatekeeper

import (
	"github.com/diamondburned/arikawa/v3/gateway"
)

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
