package reminders

import (
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

// Bot ...
type Bot struct {
	*bot.Bot

	activeReminders map[uint64]bool
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Reminders"

	b := &Bot{Bot: bot, activeReminders: make(map[uint64]bool)}

	go b.setReminderLoop()
	return
}

func (bot *Bot) setReminderLoop() {
	return
}
