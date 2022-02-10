package cc

import (
	"context"
	"time"

	"1f320.xyz/x/parameters"
	"github.com/starshine-sys/bcr"
	lua "github.com/yuin/gopher-lua"
)

type State struct {
	ctx    *bcr.Context
	params *parameters.Parameters

	ls *lua.LState
}

// NewState creates a new Lua state.
// It also adds functions.
func NewState(ctx *bcr.Context, params *parameters.Parameters) *State {
	s := &State{
		ctx:    ctx,
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

	return s
}

// Do executes the given Lua source with the given timeout.
// The Lua state is closed once this function returns and should not be reused.
// The timeout must be between 1ns and 5m inclusive.
func (s *State) Do(source string, timeout time.Duration) error {
	defer s.ls.Close()

	if timeout <= 0 || timeout > 5*time.Minute {
		timeout = time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s.ls.SetContext(ctx)

	return s.ls.DoString(source)
}

func (s *State) Close() {
	if !s.ls.IsClosed() {
		s.ls.Close()
	}
}
