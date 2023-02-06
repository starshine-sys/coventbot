// SPDX-License-Identifier: AGPL-3.0-only
package cc

import (
	"fmt"
	"strconv"
	"time"

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

// embed transforms a lua table into an embed
func (s *State) embed(t *lua.LTable) (e discord.Embed, ok bool) {
	t.ForEach(func(key, val lua.LValue) {
		switch key.String() {
		case "title":
			ok = true
			e.Title = val.String()
		case "description":
			ok = true
			e.Description = val.String()
		case "url":
			ok = true
			e.URL = val.String()
		case "timestamp":
			t, err := toTimestamp(val)
			if err != nil {
				return
			}
			ok = true
			e.Timestamp = discord.NewTimestamp(t)
		case "color":
			v, innerOk := val.(lua.LNumber)
			if !innerOk {
				return
			}
			ok = true
			e.Color = discord.Color(v)
		case "footer":
			f, innerOk := s.embedFooter(val)
			if !innerOk {
				return
			}
			ok = true
			e.Footer = &f
		case "author":
			f, innerOk := s.embedAuthor(val)
			if !innerOk {
				return
			}
			ok = true
			e.Author = &f
		case "image":
			e.Image = &discord.EmbedImage{
				URL: val.String(),
			}
		case "thumbnail":
			e.Thumbnail = &discord.EmbedThumbnail{
				URL: val.String(),
			}
		case "fields":
			fields, innerOk := s.embedFields(val)
			if !innerOk {
				return
			}
			ok = true
			e.Fields = fields
		}
	})

	return e, ok
}

func (s *State) embedFooter(v lua.LValue) (f discord.EmbedFooter, ok bool) {
	t, ok := v.(*lua.LTable)
	if !ok {
		return f, false
	}

	switch t.Len() {
	default:
		fallthrough
	case 2:
		f.Icon = t.RawGetInt(2).String()
		fallthrough
	case 1:
		f.Text = t.RawGetInt(1).String()
	case 0:
		return f, false
	}
	return f, true
}

func (s *State) embedAuthor(v lua.LValue) (e discord.EmbedAuthor, ok bool) {
	t, ok := v.(*lua.LTable)
	if !ok {
		return e, false
	}

	switch t.Len() {
	default:
		fallthrough
	case 2:
		e.Icon = t.RawGetInt(2).String()
		fallthrough
	case 1:
		e.Name = t.RawGetInt(1).String()
	case 0:
		return e, false
	}
	return e, true
}

func (s *State) embedFields(v lua.LValue) (fields []discord.EmbedField, ok bool) {
	t, ok := v.(*lua.LTable)
	if !ok {
		return nil, false
	}

	t.ForEach(func(_, val lua.LValue) {
		f, innerOk := s.embedField(val)
		if innerOk {
			ok = true
			fields = append(fields, f)
		}
	})

	return fields, ok
}

func (s *State) embedField(v lua.LValue) (f discord.EmbedField, ok bool) {
	t, ok := v.(*lua.LTable)
	if !ok {
		return f, false
	}

	switch t.Len() {
	default:
		fallthrough
	case 3:
		v := t.RawGetInt(3)
		b, ok := v.(lua.LBool)
		if ok {
			f.Inline = bool(b)
		}
		fallthrough
	case 2:
		f.Name = t.RawGetInt(1).String()
		f.Value = t.RawGetInt(2).String()
	case 1, 0:
		return f, false
	}
	return f, true
}

// toTimestamp turns the given Lua value into a timestamp.
// Numbers are interpreted as Unix time.
// Strings are:
// - parsed as RFC3339
// - parsed as RFC1123
// - parsed as a Unix time number
func toTimestamp(v lua.LValue) (time.Time, error) {
	float, ok := v.(lua.LNumber)
	if ok {
		return time.Unix(int64(float), 0), nil
	}

	t, err := time.Parse(time.RFC3339, v.String())
	if err == nil {
		return t, nil
	}

	t, err = time.Parse(time.RFC1123, v.String())
	if err == nil {
		return t, nil
	}

	i, err := strconv.ParseInt(v.String(), 10, 64)
	if err != nil {
		return t, err
	}
	return time.Unix(i, 0), nil
}
