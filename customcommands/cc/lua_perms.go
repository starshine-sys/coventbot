// SPDX-License-Identifier: AGPL-3.0-only
package cc

import (
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	lua "github.com/yuin/gopher-lua"
)

// This file contains functions used for custom commands.

func (s *State) setRequireFuncs() {
	t := s.ls.NewTable()
	t.RawSetString("role", s.ls.NewFunction(s.requireRole))
	t.RawSetString("not_role", s.ls.NewFunction(s.blockRole))

	s.ls.SetGlobal("require", t)
}

// role :: require any of the given roles
// - 2 arguments
// 1. error to raise (string)
// 2. role or list of roles (number, string, or table)
// - returns 0
func (s *State) requireRole(ls *lua.LState) int {
	reason := ls.CheckString(1)
	var roleIDs []discord.RoleID

	v := ls.Get(2)
	switch v.Type() {
	case lua.LTNumber:
		r, err := s.ctx.ParseRole(strconv.FormatUint(uint64(v.(lua.LNumber)), 10))
		if err != nil {
			ls.RaiseError("role %d not found", uint64(v.(lua.LNumber)))
			return 0
		}
		roleIDs = append(roleIDs, r.ID)
	case lua.LTString:
		r, err := s.ctx.ParseRole(string(v.(lua.LString)))
		if err != nil {
			ls.RaiseError("role %s not found", string(v.(lua.LString)))
			return 0
		}
		roleIDs = append(roleIDs, r.ID)
	case lua.LTTable:
		t := v.(*lua.LTable)
		t.ForEach(func(_, v lua.LValue) {
			switch v.Type() {
			case lua.LTNumber:
				r, err := s.ctx.ParseRole(strconv.FormatUint(uint64(v.(lua.LNumber)), 10))
				if err != nil {
					ls.RaiseError("role %d not found", uint64(v.(lua.LNumber)))
				}
				roleIDs = append(roleIDs, r.ID)
			case lua.LTString:
				r, err := s.ctx.ParseRole(string(v.(lua.LString)))
				if err != nil {
					ls.RaiseError("role %s not found", string(v.(lua.LString)))
				}
				roleIDs = append(roleIDs, r.ID)
			default:
				ls.RaiseError("argument must be string or number")
			}
		})
	default:
		ls.RaiseError("argument must be string, number, or table")
	}

	hasAnyRole := false
	for _, r := range roleIDs {
		for _, ur := range s.ctx.Member.RoleIDs {
			if r == ur {
				hasAnyRole = true
			}
		}
	}

	if !hasAnyRole {
		s._notStacktrace(ls, reason)
	}
	return 0
}

func (s *State) blockRole(ls *lua.LState) int {
	reason := ls.CheckString(1)
	var roleIDs []discord.RoleID

	v := ls.Get(2)
	switch v.Type() {
	case lua.LTNumber:
		r, err := s.ctx.ParseRole(strconv.FormatUint(uint64(v.(lua.LNumber)), 10))
		if err != nil {
			ls.RaiseError("role %d not found", uint64(v.(lua.LNumber)))
			return 0
		}
		roleIDs = append(roleIDs, r.ID)
	case lua.LTString:
		r, err := s.ctx.ParseRole(string(v.(lua.LString)))
		if err != nil {
			ls.RaiseError("role %s not found", string(v.(lua.LString)))
			return 0
		}
		roleIDs = append(roleIDs, r.ID)
	case lua.LTTable:
		t := v.(*lua.LTable)
		t.ForEach(func(_, v lua.LValue) {
			switch v.Type() {
			case lua.LTNumber:
				r, err := s.ctx.ParseRole(strconv.FormatUint(uint64(v.(lua.LNumber)), 10))
				if err != nil {
					ls.RaiseError("role %d not found", uint64(v.(lua.LNumber)))
				}
				roleIDs = append(roleIDs, r.ID)
			case lua.LTString:
				r, err := s.ctx.ParseRole(string(v.(lua.LString)))
				if err != nil {
					ls.RaiseError("role %s not found", string(v.(lua.LString)))
				}
				roleIDs = append(roleIDs, r.ID)
			default:
				ls.RaiseError("argument must be string or number")
			}
		})
	default:
		ls.RaiseError("argument must be string, number, or table")
	}

	for _, r := range roleIDs {
		for _, ur := range s.ctx.Member.RoleIDs {
			if r == ur {
				s._notStacktrace(ls, reason)
				return 0
			}
		}
	}
	return 0
}
