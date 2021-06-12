package bot

import (
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

func (bot *Bot) requestGuildMembers(g *gateway.GuildCreateEvent) {
	bot.State.Gateway.RequestGuildMembers(gateway.RequestGuildMembersData{
		GuildID: []discord.GuildID{g.ID},
		Limit:   0,
	})
}

func (bot *Bot) guildMemberChunk(g *gateway.GuildMembersChunkEvent) {
	bot.membersMu.Lock()
	defer bot.membersMu.Unlock()
	for _, m := range g.Members {
		bot.members[memberKey{
			GuildID: g.GuildID,
			UserID:  m.User.ID,
		}] = member{m, g.GuildID}
	}
}

func (bot *Bot) memberUpdateEvent(ev *gateway.GuildMemberUpdateEvent) {
	// wait a bit so stuff can grab the old member object
	time.Sleep(time.Second)

	bot.membersMu.Lock()
	defer bot.membersMu.Unlock()
	v, ok := bot.members[memberKey{
		GuildID: ev.GuildID,
		UserID:  ev.User.ID,
	}]
	if !ok {
		return
	}

	ev.Update(&v.Member)

	bot.members[memberKey{
		GuildID: ev.GuildID,
		UserID:  ev.User.ID,
	}] = member{v.Member, ev.GuildID}
}

func (bot *Bot) memberAddEvent(ev *gateway.GuildMemberAddEvent) {
	bot.membersMu.Lock()
	defer bot.membersMu.Unlock()
	bot.members[memberKey{
		GuildID: ev.GuildID,
		UserID:  ev.User.ID,
	}] = member{ev.Member, ev.GuildID}
}

func (bot *Bot) memberRemoveEvent(ev *gateway.GuildMemberRemoveEvent) {
	time.Sleep(time.Second)

	bot.membersMu.Lock()
	defer bot.membersMu.Unlock()

	delete(bot.members, memberKey{ev.GuildID, ev.User.ID})
}

// Member gets a member from the cache, or tries fetching it from Discord
func (bot *Bot) Member(guildID discord.GuildID, userID discord.UserID) (m discord.Member, err error) {
	bot.membersMu.RLock()
	if m, ok := bot.members[memberKey{guildID, userID}]; ok {
		bot.membersMu.RUnlock()
		return m.Member, nil
	}
	bot.membersMu.RUnlock()

	gm, err := bot.State.Session.Member(guildID, userID)
	if err != nil {
		return
	}

	bot.membersMu.Lock()
	defer bot.membersMu.Unlock()
	bot.members[memberKey{
		GuildID: guildID,
		UserID:  userID,
	}] = member{*gm, guildID}
	return *gm, nil
}

// Members gets all *cached* members for a guild
func (bot *Bot) Members(guildID discord.GuildID) (members []discord.Member) {
	bot.membersMu.RLock()
	defer bot.membersMu.RUnlock()

	for _, m := range bot.members {
		if m.GuildID == guildID {
			members = append(members, m.Member)
		}
	}

	return
}

// CacheLen gets the size of the member cache
func (bot *Bot) CacheLen() int64 {
	bot.membersMu.RLock()
	defer bot.membersMu.RUnlock()
	return int64(len(bot.members))
}
