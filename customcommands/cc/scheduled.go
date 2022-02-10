package cc

import (
	"context"
	"strconv"
	"time"

	"1f320.xyz/x/parameters"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
	"github.com/starshine-sys/tribble/db"
	lua "github.com/yuin/gopher-lua"
)

type ScheduledCC struct {
	ID         int64           `json:"id"`
	Guild      discord.Guild   `json:"guild"`
	Channel    discord.Channel `json:"channel"`
	Message    discord.Message `json:"message"`
	Member     discord.Member  `json:"member"`
	Parameters string          `json:"parameters"`
}

var _ bot.Event = (*ScheduledCC)(nil)

func (dat *ScheduledCC) Execute(cctx context.Context, id int64, bot *bot.Bot) error {
	cmd, err := bot.DB.CustomCommandID(dat.Guild.ID, dat.ID)
	if err != nil {
		if err != db.ErrCommandNotFound {
			bot.Sugar.Errorf("error getting command %v: %v", dat.ID, err)
		}
		return err
	}

	ds, sid := bot.Router.StateFromGuildID(dat.Guild.ID)
	ctx := &bcr.Context{
		State:   ds,
		ShardID: sid,
		Bot:     bot.Router.Bot,
		Message: dat.Message,
		Channel: &dat.Channel,
		Guild:   &dat.Guild,
		Member:  &dat.Member,
		Router:  bot.Router,
	}

	s := NewState(bot, ctx, parameters.NewParameters(dat.Parameters, false))
	s.ls.SetGlobal("is_scheduled", lua.LBool(true))
	s.ls.SetGlobal("scheduled_args", lua.LString(dat.Parameters))
	s.isScheduled = true

	return s.Do(cctx, cmd.ID, cmd.Source)
}

func (dat *ScheduledCC) Offset() time.Duration { return time.Minute }

func (s *State) setScheduleFuncs() {
	s.ls.SetGlobal("schedule_cc", s.ls.NewFunction(s.scheduleCC))
}

func (s *State) scheduleCC(ls *lua.LState) int {
	if s.isScheduled {
		ls.RaiseError("cannot schedule a cc in a scheduled cc")
		return 1
	}

	id := s.ccID

	v := ls.Get(1)
	switch v.Type() {
	case lua.LTNil:
	case lua.LTNumber:
		id = int64(v.(lua.LNumber))
	case lua.LTString:
		i, err := strconv.ParseInt(string(v.(lua.LString)), 10, 64)
		if err != nil {
			ls.RaiseError("argument 1 should be nil, number, or integer string")
			return 1
		}
		id = i
	default:
		ls.RaiseError("argument 1 should be nil, number, or integer string")
		return 1
	}

	_, err := s.bot.DB.CustomCommandID(s.ctx.Guild.ID, id)
	if err != nil {
		ls.RaiseError("command id %d does not exist on this server", id)
		return 1
	}

	dur, err := time.ParseDuration(ls.CheckString(2))
	if err != nil {
		ls.RaiseError("%q cannot be parsed as a duration", ls.CheckString(2))
		return 1
	}

	if dur < time.Minute || dur > 7*24*time.Hour {
		ls.RaiseError("delay must be between 1 minute and 7 days")
		return 1
	}

	params := s.ctx.RawArgs
	v = ls.Get(3)
	switch v.Type() {
	case lua.LTNil:
	case lua.LTString:
		params = string(v.(lua.LString))
	default:
		ls.TypeError(3, lua.LTString)
		return 1
	}

	sid, err := s.bot.Scheduler.Add(time.Now().UTC().Add(dur), &ScheduledCC{
		ID:         id,
		Guild:      *s.ctx.Guild,
		Channel:    *s.ctx.Channel,
		Message:    s.ctx.Message,
		Member:     *s.ctx.Member,
		Parameters: params,
	})
	if err != nil {
		s.bot.Sugar.Errorf("error scheduling cc")
		ls.RaiseError("could not schedule cc")
		return 1
	}

	ls.Push(lua.LNumber(sid))
	return 1
}
