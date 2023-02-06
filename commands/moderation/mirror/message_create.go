// SPDX-License-Identifier: AGPL-3.0-only
package mirror

import (
	"strings"

	"github.com/diamondburned/arikawa/v3/gateway"
)

const (
	yagID  = 204255221017214977
	carlID = 235148962103951360
)

func (bot *Bot) messageCreate(m *gateway.MessageCreateEvent) {
	if !m.GuildID.IsValid() {
		return
	}

	if m.Author.ID == yagID {
		bot.processYAG(m)
		return
	}

	if m.Author.ID == carlID {
		if len(m.Embeds) > 0 {
			bot.processCarlLog(m)
			return
		}

		if strings.Contains(m.Content, "Note taken.") {
			bot.processCarlNote(m)
			return
		}
	}
}
