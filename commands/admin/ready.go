package admin

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

func (bot *Bot) ready(_ *gateway.ReadyEvent) {
	bot.updateStatus()
	return
}

func (bot *Bot) updateStatus() {
	s := bot.Settings()

	usd := gateway.UpdateStatusData{
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

	err := bot.State.Gateway.UpdateStatus(usd)
	if err != nil {
		bot.Sugar.Errorf("Error setting status: %v", err)
	}
}
