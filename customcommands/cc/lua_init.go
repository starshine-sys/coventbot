package cc

import (
	"github.com/diamondburned/arikawa/v3/discord"
	lua "github.com/yuin/gopher-lua"
)

func (s *State) initGuild() {
	t := s.ls.NewTable()
	t.RawSetString("id", lua.LString(s.ctx.Guild.ID.String()))
	t.RawSetString("name", lua.LString(s.ctx.Guild.Name))
	t.RawSetString("icon", lua.LString(s.ctx.Guild.Icon))
	t.RawSetString("icon_url", lua.LString(s.ctx.Guild.IconURL()))

	s.ls.SetGlobal("guild", t)
	s.ls.SetGlobal("channel", s.channelToLua(*s.ctx.Channel))
}

func (s *State) initMessage() {
	t := s.messageToLua(s.ctx.Message)

	s.ls.SetGlobal("message", t)
}

func (s *State) messageToLua(m discord.Message) *lua.LTable {
	t := s.ls.NewTable()
	t.RawSetString("id", lua.LString(m.ID.String()))
	t.RawSetString("channel_id", lua.LString(m.ChannelID.String()))
	t.RawSetString("guild_id", lua.LString(m.GuildID.String()))

	t.RawSetString("pinned", lua.LBool(m.Pinned))
	t.RawSetString("mention_everyone", lua.LBool(m.MentionEveryone))

	t.RawSetString("author", s.userToLua(m.Author))
	t.RawSetString("webhook_id", lua.LString(m.WebhookID.String()))

	t.RawSetString("content", lua.LString(m.Content))

	return t
}

func (s *State) userToLua(u discord.User) *lua.LTable {
	t := s.ls.NewTable()

	t.RawSetString("id", lua.LString(u.ID.String()))
	t.RawSetString("username", lua.LString(u.Username))
	t.RawSetString("discriminator", lua.LString(u.Discriminator))
	t.RawSetString("tag", lua.LString(u.Tag()))
	t.RawSetString("avatar", lua.LString(u.Avatar))
	t.RawSetString("avatar_url", lua.LString(u.AvatarURL()))

	t.RawSetString("bot", lua.LBool(u.Bot))

	return t
}

func (s *State) channelToLua(ch discord.Channel) *lua.LTable {
	t := s.ls.NewTable()

	t.RawSetString("id", lua.LString(ch.ID.String()))
	t.RawSetString("guild_id", lua.LString(ch.GuildID.String()))
	t.RawSetString("parent_id", lua.LString(ch.ParentID.String()))

	t.RawSetString("type", lua.LNumber(ch.Type))
	t.RawSetString("nsfw", lua.LBool(ch.NSFW))
	t.RawSetString("position", lua.LNumber(ch.Position))

	t.RawSetString("name", lua.LString(ch.Name))
	t.RawSetString("topic", lua.LString(ch.Topic))

	return t
}
