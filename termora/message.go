// SPDX-License-Identifier: AGPL-3.0-only
package termora

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/diamondburned/arikawa/v3/gateway"
)

const embedColour = 0xd14171

var linkRegexp = regexp.MustCompile(`\[\[(.*?)\]\]`)

func (bot *Bot) messageCreate(m *gateway.MessageCreateEvent) {
	if !m.GuildID.IsValid() || m.Author.Bot ||
		m.GuildID != bot.Config.Termora.Guild || m.ChannelID == bot.Config.Termora.TermChannel {
		return
	}

	if !linkRegexp.MatchString(m.Content) {
		return
	}

	ts, err := bot.Terms()
	if err != nil {
		bot.Sugar.Errorf("Error fetching terms: %v", err)
	}

	var output string
	matches := linkRegexp.FindAllStringSubmatch(m.Content, -1)

	for _, i := range matches {
		if len(i) < 2 {
			continue
		}

		t, found := findTerm(ts, i[1])
		if found {
			output += fmt.Sprintf("%v: <https://termora.org/term/%v>\n", t.Name, t.ID)
		}
	}

	if output == "" {
		return
	}

	s, _ := bot.Router.StateFromGuildID(m.GuildID)

	_, err = s.SendMessage(m.ChannelID, output)
	if err != nil {
		bot.Sugar.Errorf("Error sending term link message: %v", err)
	}
}

func findTerm(ts []Term, input string) (t Term, found bool) {
	for _, t := range ts {
		if strings.EqualFold(t.Name, input) {
			return t, true
		}

		for _, alias := range t.Aliases {
			if strings.EqualFold(alias, input) {
				return t, true
			}
		}
	}

	return t, false
}
