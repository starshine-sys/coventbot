package cc

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	lua "github.com/yuin/gopher-lua"
)

// This file contains functions used for custom commands.

func (s *State) setMessageFuncs() {
	s.ls.SetGlobal("send_message", s.ls.NewFunction(s.sendMessage))
	s.ls.SetGlobal("react", s.ls.NewFunction(s.react))
}

func (s *State) sendMessage(ls *lua.LState) int {
	if s.discordCalls >= 10 {
		ls.RaiseError("maximum number of Discord calls reached (10)")
		return 1
	}

	chID := s.ctx.Message.ChannelID

	// first argument is channel ID
	id, ok := s._getChannelID(ls, 1)
	if ok {
		chID = id
	}

	ch, err := s.ctx.State.Channel(chID)
	if err != nil {
		ls.RaiseError("channel not found in state")
		return 1
	}

	if ch.GuildID != s.ctx.Message.GuildID {
		ls.RaiseError("channel is not in this guild")
		return 1
	}

	data := api.SendMessageData{
		AllowedMentions: allowedMentions(s.ctx),
	}

	v := ls.Get(2)
	switch v.Type() {
	case lua.LTString:
		data.Content = string(v.(lua.LString))
	case lua.LTTable:
		t := v.(*lua.LTable)

		e, ok := s.embed(t)
		if ok {
			data.Embeds = []discord.Embed{e}
		}

		v := t.RawGetString("content")
		s, ok := v.(lua.LString)
		if ok {
			data.Content = string(s)
		}

		v = t.RawGetString("allow_all_mentions")
		b, ok := v.(lua.LBool)
		if ok && bool(b) {
			data.AllowedMentions = nil
		}
	default:
		data.Content = v.String()
	}

	s.discordCalls++
	msg, err := s.ctx.State.SendMessageComplex(chID, data)
	if err != nil {
		ls.RaiseError("error sending message: %s", err.Error())
		return 1
	}

	ls.Push(lua.LString(msg.ID.String()))
	return 1
}

func (s *State) react(ls *lua.LState) int {
	if s.discordCalls >= 10 {
		ls.RaiseError("maximum number of Discord calls reached (10)")
		return 0
	}

	chID := s.ctx.Message.ChannelID
	mID := s.ctx.Message.ID

	// first argument is channel ID
	if v, ok := s._getChannelID(ls, 1); ok {
		chID = v
	}

	if v, ok := s._getMessageID(ls, 2); ok {
		mID = v
	}

	msg, err := s.ctx.State.Message(chID, mID)
	if err != nil {
		ls.RaiseError("message %d/%d not found", chID, mID)
		return 0
	}
	if msg.GuildID != s.ctx.Message.GuildID {
		ls.RaiseError("message %d/%d not in this guild", chID, mID)
		return 0
	}

	react := s._getString(ls, 3)

	s.discordCalls++
	err = s.ctx.State.React(msg.ChannelID, msg.ID, discord.APIEmoji(react))
	if err != nil {
		ls.RaiseError("error reacting to message: %s", err.Error())
	}
	return 0
}

func allowedMentions(ctx *bcr.Context) *api.AllowedMentions {
	mentions := &api.AllowedMentions{
		Parse: []api.AllowedMentionType{api.AllowUserMention},
	}

	if ctx.Guild == nil {
		return mentions
	}

	for _, r := range ctx.Guild.Roles {
		if r.Mentionable {
			mentions.Roles = append(mentions.Roles, r.ID)
		}
	}
	return mentions
}
