package bot

import "github.com/diamondburned/arikawa/v2/gateway"

func (bot *Bot) guildMemberUpdate(ev *gateway.GuildMemberUpdateEvent) {
	m, err := bot.State.Member(ev.GuildID, ev.User.ID)
	if err != nil {
		return
	}

	ev.Update(m)

	bot.State.MemberSet(ev.GuildID, *m)
}
