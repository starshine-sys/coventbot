// SPDX-License-Identifier: AGPL-3.0-only
package reminders

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"codeberg.org/eviedelta/detctime/durationparser"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/gateway"
	"gitlab.com/1f320/x/duration"
)

const dateReText = `(?i)hey %v,? remind (me|us) to ([\s\S]+) (in|at) (.+)`
const tomorrowReText = `(?i)hey %v,? remind (me|us) to ([\s\S]+) tomorrow`

var dateRegexp, tomorrowRegexp *regexp.Regexp

var o sync.Once

func (bot *Bot) setNameOnce() {
	dateRegexp = regexp.MustCompile(fmt.Sprintf(dateReText, strings.ToLower(bot.Router.Bot.Username)))
	tomorrowRegexp = regexp.MustCompile(fmt.Sprintf(tomorrowReText, strings.ToLower(bot.Router.Bot.Username)))
}

func (bot *Bot) messageCreate(m *gateway.MessageCreateEvent) {
	if m.Author.Bot {
		return
	}

	if dateRegexp == nil || tomorrowRegexp == nil {
		o.Do(bot.setNameOnce)
		return
	}

	s, _ := bot.Router.StateFromGuildID(m.GuildID)

	var text string
	var ts string
	if dateRegexp.MatchString(m.Content) {
		groups := dateRegexp.FindStringSubmatch(m.Content)
		if len(groups) < 5 {
			return
		}
		text = groups[2]
		ts = groups[4]
	} else if tomorrowRegexp.MatchString(m.Content) {
		groups := tomorrowRegexp.FindStringSubmatch(m.Content)
		if len(groups) < 3 {
			return
		}
		text = groups[2]
		ts = "tomorrow"
	} else {
		return
	}

	loc := bot.userTime(m.Author.ID)

	t, _, err := ParseTime(strings.Fields(ts), loc)
	if err != nil {
		var dur time.Duration

		fields := strings.Fields(ts)
		for i := len(fields) - 1; i > 0; i++ {
			dur, err = durationparser.Parse(strings.Join(fields[:i], " "))
			if err == nil {
				break
			}
		}

		if err != nil {
			dur, err = durationparser.Parse(fields[0])
			if err != nil {
				s.SendMessageComplex(m.ChannelID, api.SendMessageData{
					Content: fmt.Sprintf("Either you put the time in the wrong location (must be at the very end), or I didn't understand the format you used. Try again?"),
					AllowedMentions: &api.AllowedMentions{
						Parse: []api.AllowedMentionType{},
					},
				})
				return
			}
		}

		t = time.Now().In(loc).Add(dur)
	}

	var id uint64
	err = bot.DB.Pool.QueryRow(context.Background(), `insert into reminders
	(user_id, message_id, channel_id, server_id, reminder, expires)
	values
	($1, $2, $3, $4, $5, $6) returning id`, m.Author.ID, m.ID, m.ChannelID, m.GuildID, text, t.UTC()).Scan(&id)
	if err != nil {
		bot.Sugar.Errorf("Error storing reminder: %v", err)
		return
	}

	if len(text) > 128 {
		text = text[:128] + "..."
	}

	name := m.Author.Username
	if m.Member != nil && m.Member.Nick != "" {
		name = m.Member.Nick
	}

	msg, err := s.SendMessageComplex(m.ChannelID, api.SendMessageData{
		Content: fmt.Sprintf("Okay %v, I'll remind you about **%v** in %v. (<t:%v>, #%v)", name, text, duration.Format(time.Until(t.Add(time.Second))), t.Unix(), id),
		AllowedMentions: &api.AllowedMentions{
			Parse: []api.AllowedMentionType{},
		},
	})
	if err != nil {
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update reminders set message_id = $1 where id = $2", msg.ID, id)
	if err != nil {
		bot.Sugar.Errorf("Error storing reminder: %v", err)
	}
	return
}
