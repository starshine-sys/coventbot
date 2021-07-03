package bot

import "github.com/diamondburned/arikawa/v3/gateway"

func (bot *Bot) guildMemberUpdate(ev *gateway.GuildMemberUpdateEvent) {
	s, _ := bot.Router.StateFromGuildID(ev.GuildID)

	m, err := s.Member(ev.GuildID, ev.User.ID)
	if err != nil {
		return
	}

	ev.Update(m)

	s.MemberSet(ev.GuildID, *m, true)
}
