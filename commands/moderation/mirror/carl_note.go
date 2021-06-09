package mirror

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/starshine-sys/tribble/db"
)

var carlNoteRegexp = regexp.MustCompile(`.*\n\*\*Note:\*\* ([\s\S]+)`)

func (bot *Bot) processCarlNote(m *gateway.MessageCreateEvent) {
	groups := carlNoteRegexp.FindStringSubmatch(m.Content)
	if len(groups) < 2 {
		bot.Sugar.Infof("Didn't match regex")
		return
	}

	reason := groups[1]

	msgs, err := bot.State.MessagesAround(m.ChannelID, m.ID, 10)
	if err != nil {
		bot.Sugar.Errorf("Error getting messages: %v", err)
		return
	}
	var orig discord.Message
	var found bool
	for _, m := range msgs {
		if strings.Contains(m.Content, reason) && !m.Author.Bot {
			orig = m
			found = true
			break
		}
	}

	if !found {
		bot.Sugar.Infof("No original message found")
		return
	}

	var userID discord.UserID
	for _, arg := range strings.Fields(orig.Content) {
		u, err := bot.ParseUser(arg)
		if err == nil {
			userID = u.ID
			break
		}
	}

	if userID == 0 {
		bot.Sugar.Infof("No user ID found")
		return
	}

	note := db.Note{
		ServerID:  m.GuildID,
		UserID:    userID,
		Note:      reason,
		Moderator: orig.Author.ID,
		Created:   time.Now().UTC(),
	}

	note, err = bot.DB.NewNote(note)
	if err != nil {
		bot.Sugar.Errorf("Error importing Carl-bot note: %v", err)
		return
	}
}

var (
	userMentionRegex = regexp.MustCompile("<@!?(\\d+)>")
	idRegex          = regexp.MustCompile("^\\d+$")
)

// ParseUser finds a user by mention or ID
func (bot *Bot) ParseUser(s string) (u *discord.User, err error) {
	if idRegex.MatchString(s) {
		sf, err := discord.ParseSnowflake(s)
		if err != nil {
			return nil, err
		}
		return bot.State.User(discord.UserID(sf))
	}

	if userMentionRegex.MatchString(s) {
		matches := userMentionRegex.FindStringSubmatch(s)
		if len(matches) < 2 {
			return nil, errors.New("user not found")
		}
		sf, err := discord.ParseSnowflake(matches[1])
		if err != nil {
			return nil, err
		}
		return bot.State.User(discord.UserID(sf))
	}

	return nil, errors.New("user not found")
}
