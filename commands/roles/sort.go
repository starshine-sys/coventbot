package roles

import (
	"errors"

	"github.com/diamondburned/arikawa/v2/discord"
)

func (bot *Bot) roles(guildID discord.GuildID, roleIDs []uint64) (roles []discord.Role, err error) {
	rls, err := bot.State.Roles(guildID)
	if err != nil {
		return nil, err
	}

	if len(rls) == 0 {
		return nil, errors.New("couldn't get roles")
	}

	for _, r := range rls {
		for _, id := range roleIDs {
			if r.ID == discord.RoleID(id) {
				roles = append(roles, r)
			}
		}
	}

	return
}

func sortByName(rls []discord.Role) func(i, j int) bool {
	return func(i, j int) bool {
		return rls[i].Name < rls[j].Name
	}
}

func sortByPosition(rls []discord.Role) func(i, j int) bool {
	return func(i, j int) bool {
		return rls[i].Position < rls[j].Position
	}
}
