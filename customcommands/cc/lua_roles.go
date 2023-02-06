// SPDX-License-Identifier: AGPL-3.0-only
package cc

import (
	"github.com/diamondburned/arikawa/v3/api"
	lua "github.com/yuin/gopher-lua"
)

func (s *State) setRoleFuncs() {
	s.ls.SetGlobal("add_role", s.ls.NewFunction(s.addRole))
	s.ls.SetGlobal("remove_role", s.ls.NewFunction(s.removeRole))
}

func (s *State) addRole(ls *lua.LState) int {
	if s.discordCalls >= 10 {
		ls.RaiseError("maximum number of Discord calls reached (10)")
		return 1
	}

	userID, ok := s._getUserID(ls, 1)
	if !ok {
		ls.ArgError(1, "must be a number or string")
		return 0
	}

	roleID, ok := s._getRoleID(ls, 2)
	if !ok {
		ls.ArgError(1, "must be a number or string")
		return 0
	}

	s.discordCalls++
	err := s.ctx.State.AddRole(s.ctx.Guild.ID, userID, roleID, api.AddRoleData{
		AuditLogReason: "`add_role` in custom command",
	})
	if err != nil {
		ls.RaiseError(err.Error())
		return 0
	}

	return 0
}

func (s *State) removeRole(ls *lua.LState) int {
	if s.discordCalls >= 10 {
		ls.RaiseError("maximum number of Discord calls reached (10)")
		return 1
	}

	userID, ok := s._getUserID(ls, 1)
	if !ok {
		ls.ArgError(1, "must be a number or string")
		return 0
	}

	roleID, ok := s._getRoleID(ls, 2)
	if !ok {
		ls.ArgError(1, "must be a number or string")
		return 0
	}

	s.discordCalls++
	err := s.ctx.State.RemoveRole(s.ctx.Guild.ID, userID, roleID, "`remove_role` in custom command")
	if err != nil {
		ls.RaiseError(err.Error())
		return 0
	}

	return 0
}
