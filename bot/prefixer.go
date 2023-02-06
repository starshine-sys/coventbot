// SPDX-License-Identifier: AGPL-3.0-only
package bot

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
)

// CheckPrefix checks the prefix
func (bot *Bot) CheckPrefix(m discord.Message) int {
	for _, mention := range []string{fmt.Sprintf("<@%v>", bot.Router.Bot.ID), fmt.Sprintf("<@!%v>", bot.Router.Bot.ID)} {
		if strings.HasPrefix(m.Content, mention) {
			return len(mention)
		}
	}

	if !m.GuildID.IsValid() {
		return bot.Router.DefaultPrefixer(m)
	}

	p, err := bot.DB.Prefixes(m.GuildID)
	if err != nil || len(p) == 0 {
		return bot.Router.DefaultPrefixer(m)
	}

	for _, p := range p {
		if strings.HasPrefix(strings.ToLower(m.Content), p) {
			return len(p)
		}
	}

	return -1
}
