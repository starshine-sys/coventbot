// SPDX-License-Identifier: AGPL-3.0-only
package bot

import "github.com/diamondburned/arikawa/v3/gateway"

func (bot *Bot) guildMemberUpdate(ev *gateway.GuildMemberUpdateEvent) {
	s, _ := bot.Router.StateFromGuildID(ev.GuildID)

	m, err := s.Member(ev.GuildID, ev.User.ID)
	if err != nil {
		return
	}

	ev.UpdateMember(m)

	s.MemberSet(ev.GuildID, m, true)
}
