package cc

import (
	"fmt"
	"strconv"

	"github.com/diamondburned/arikawa/v3/discord"
	lua "github.com/yuin/gopher-lua"
)

// These functions are not used in Lua directly

func (s *State) _getChannelID(ls *lua.LState, idx int) *discord.ChannelID {
	v := ls.Get(idx)
	switch v.Type() {
	case lua.LTNil:
		return nil
	case lua.LTNumber:
		id := uint64(v.(lua.LNumber))
		return (*discord.ChannelID)(&id)
	case lua.LTString:
		i, err := strconv.ParseUint(string(v.(lua.LString)), 10, 64)
		if err != nil {
			return nil
		}

		return (*discord.ChannelID)(&i)
	default:
		ls.RaiseError("argument %d must be a string, number, or nil", idx)
		return nil
	}
}

func (s *State) _getMessageID(ls *lua.LState, idx int) *discord.MessageID {
	v := ls.Get(idx)
	switch v.Type() {
	case lua.LTNil:
		return nil
	case lua.LTNumber:
		id := uint64(v.(lua.LNumber))
		return (*discord.MessageID)(&id)
	case lua.LTString:
		i, err := strconv.ParseUint(string(v.(lua.LString)), 10, 64)
		if err != nil {
			return nil
		}

		return (*discord.MessageID)(&i)
	default:
		ls.RaiseError("argument %d must be a string, number, or nil", idx)
		return nil
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
