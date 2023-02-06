// SPDX-License-Identifier: AGPL-3.0-only
package cc

import (
	"context"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
	lua "github.com/yuin/gopher-lua"
	"gitlab.com/1f320/x/parameters"
)

// State is a Lua state.
// It is *not* thread safe and should be discarded once a command has been run.
type State struct {
	ctx         *bcr.Context
	bot         *bot.Bot
	params      *parameters.Parameters
	ccID        int64
	isScheduled bool

	ls *lua.LState

	discordCalls int
}

// NewState creates a new Lua state.
func NewState(bot *bot.Bot, ctx *bcr.Context, params *parameters.Parameters) *State {
	s := &State{
		ctx:    ctx,
		bot:    bot,
		params: params,
	}

	s.ls = lua.NewState(lua.Options{
		CallStackSize:       120,
		RegistrySize:        1024 * 20,
		RegistryMaxSize:     1024 * 80,
		RegistryGrowStep:    32,
		MinimizeStackMemory: true,
		SkipOpenLibs:        true,
	})

	// open libraries
	fn := lua.OpenBaseFiltered([]string{
		"collectgarbage",
		"dofile",
		"getfenv",
		"getmetatable",
		"load",
		"loadfile",
		"loadstring",
		"print",
		"require",
		"setfenv",
		"setmetatable",
	})

	s.ls.Push(s.ls.NewFunction(fn))
	s.ls.Push(lua.LString(""))
	s.ls.Call(1, 0)

	s.ls.OpenLibs()

	lua.OpenTable(s.ls)
	lua.OpenString(s.ls)
	lua.OpenMath(s.ls)

	// set globals
	s.initGuild()
	s.initMessage()

	// initialize all functions
	s.setMessageFuncs()
	s.setArgumentFuncs()
	s.setRequireFuncs()
	s.setScheduleFuncs()
	s.setRoleFuncs()

	return s
}

// Do executes the given Lua source with the given timeout.
// The Lua state is closed once this function returns and should not be reused.
// The timeout must be between 1ns and 5m inclusive.
func (s *State) Do(ctx context.Context, id int64, source string) error {
	defer s.ls.Close()

	s.ls.SetGlobal("cc_id", lua.LNumber(id))
	s.ccID = id

	s.ls.SetContext(ctx)

	return s.ls.DoString(source)
}

// Load loads the given string and compiles it.
// It is not stored; this is to be used for giving code a cursory check before saving it.
func (s *State) Load(source string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.ls.SetContext(ctx)

	_, err := s.ls.LoadString(source)
	return err
}

func (s *State) Close() {
	if !s.ls.IsClosed() {
		s.ls.Close()
	}
}

// If err is a user error (such as one raised from the require. functions), pretty-print that error and return true.
func (s *State) FilterErrors(err error) bool {
	if v, ok := err.(*lua.ApiError); ok {
		v, ok := v.Object.(lua.LString)
		if ok {
			// yes this is hacky, no i don't care
			if strings.Contains(string(v), "❌") {

				_, _ = s.ctx.State.SendMessageComplex(s.ctx.Message.ChannelID, api.SendMessageData{
					Content:         string(v),
					AllowedMentions: &api.AllowedMentions{},
				})
				return true
			}
		}
	}
	return false
}
