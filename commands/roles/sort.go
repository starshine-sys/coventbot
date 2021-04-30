package roles

import (
	"sort"

	"github.com/diamondburned/arikawa/v2/discord"
)

// SortByPos sorts a list of roles by position.
func SortByPos(roles []discord.Role) {
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Position < roles[j].Position
	})
}

// SortByName sorts a list of roles by name.
func SortByName(roles []discord.Role) {
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Name < roles[j].Name
	})
}

// FilterRoles filters a slice of roles, returning only those whose IDs match the given list.
func FilterRoles(ids []discord.RoleID, roles []discord.Role) (filtered []discord.Role) {
	for _, r := range roles {
		for _, i := range ids {
			if r.ID == i {
				filtered = append(filtered, r)
			}
		}
	}

	return
}

// HasRole returns true if the given member has the given role ID.
func HasRole(id discord.RoleID, member *discord.Member) (hasRole bool) {
	for _, r := range member.RoleIDs {
		if id == r {
			return true
		}
	}

	return false
}
