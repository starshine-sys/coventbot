package cc

import (
	"fmt"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	lua "github.com/yuin/gopher-lua"
)

// These functions are not used in Lua directly

func (s *State) _getChannelID(ls *lua.LState, idx int) (discord.ChannelID, bool) {
	v := ls.Get(idx)
	switch v.Type() {
	case lua.LTNil:
		return discord.NullChannelID, false
	case lua.LTNumber:
		id := uint64(v.(lua.LNumber))
		return discord.ChannelID(id), true
	case lua.LTString:
		i, err := strconv.ParseUint(string(v.(lua.LString)), 10, 64)
		if err != nil {
			return discord.NullChannelID, false
		}

		return discord.ChannelID(i), true
	default:
		return discord.NullChannelID, false
	}
}

func (s *State) _getMessageID(ls *lua.LState, idx int) (discord.MessageID, bool) {
	v := ls.Get(idx)
	switch v.Type() {
	case lua.LTNil:
		return discord.NullMessageID, false
	case lua.LTNumber:
		id := uint64(v.(lua.LNumber))
		return discord.MessageID(id), true
	case lua.LTString:
		i, err := strconv.ParseUint(string(v.(lua.LString)), 10, 64)
		if err != nil {
			return discord.NullMessageID, false
		}

		return discord.MessageID(i), true
	default:
		return discord.NullMessageID, false
	}
}

func (s *State) _getUserID(ls *lua.LState, idx int) (discord.UserID, bool) {
	v := ls.Get(idx)
	switch v.Type() {
	case lua.LTNil:
		return discord.NullUserID, false
	case lua.LTNumber:
		id := uint64(v.(lua.LNumber))
		return discord.UserID(id), true
	case lua.LTString:
		i, err := strconv.ParseUint(string(v.(lua.LString)), 10, 64)
		if err != nil {
			return discord.NullUserID, false
		}

		return discord.UserID(i), true
	default:
		return discord.NullUserID, false
	}
}

func (s *State) _getRoleID(ls *lua.LState, idx int) (discord.RoleID, bool) {
	v := ls.Get(idx)
	switch v.Type() {
	case lua.LTNil:
		return discord.NullRoleID, false
	case lua.LTNumber:
		id := uint64(v.(lua.LNumber))
		return discord.RoleID(id), true
	case lua.LTString:
		i, err := strconv.ParseUint(string(v.(lua.LString)), 10, 64)
		if err != nil {
			return discord.NullRoleID, false
		}

		return discord.RoleID(i), true
	default:
		return discord.NullRoleID, false
	}
}

func (s *State) _getString(ls *lua.LState, idx int) string {
	v := ls.Get(idx)
	switch v.Type() {
	case lua.LTNil:
		return ""
	case lua.LTString:
		return string(v.(lua.LString))
	default:
		ls.RaiseError("argument %d must be a string or nil", idx)
		return ""
	}
}

// This error will be prepended with "❌"
func (s *State) _notStacktrace(ls *lua.LState, tmpl string, args ...interface{}) {
	ls.Push(lua.LString(fmt.Sprintf("❌ "+tmpl, args...)))
	ls.Panic(ls)
}
