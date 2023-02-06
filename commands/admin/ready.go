// SPDX-License-Identifier: AGPL-3.0-only
package admin

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

func (bot *Bot) updateStatus(state *state.State) {
	s := bot.Settings()

	usd := gateway.UpdatePresenceCommand{
		Status:     s.Status,
		Activities: []discord.Activity{},
	}

	if s.ActivityType != "" && s.Activity != "" {
		t := discord.GameActivity

		switch s.ActivityType {
		case "listening", "listening to":
			t = discord.ListeningActivity
		case "watching":
			t = discord.WatchingActivity
		case "playing":
			t = discord.GameActivity
		}

		usd.Activities = []discord.Activity{{
			Name: s.Activity,
			Type: t,
		}}
	}

	err := state.Gateway().Send(context.Background(), &usd)
	if err != nil {
		bot.Sugar.Errorf("Error setting status: %v", err)
	}
}
