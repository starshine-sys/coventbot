package pklog

import (
	"context"
	"regexp"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

var botsToCheck = []discord.UserID{466378653216014359}

var (
	linkRegex   = regexp.MustCompile(`^https:\/\/discord.com\/channels\/\d+\/(\d+)\/\d+$`)
	footerRegex = regexp.MustCompile(`^System ID: (\w{5,6}) \| Member ID: (\w{5,6}) \| Sender: .+ \((\d+)\) \| Message ID: (\d+) \| Original Message ID: \d+$`)
)

func (bot *Bot) pkMessageCreate(m *gateway.MessageCreateEvent) {
	var shouldLog bool
	bot.DB.Pool.QueryRow(context.Background(), "select (pk_log_channel != 0) from servers where id = $1", m.GuildID).Scan(&shouldLog)
	if !shouldLog {
		return
	}

	// only handle PK message events
	var isPK bool
	for _, u := range botsToCheck {
		if m.Author.ID == u {
			isPK = true
			break
		}
	}
	if !isPK {
		return
	}

	// only handle events that are *probably* a log message
	if len(m.Embeds) == 0 || !linkRegex.MatchString(m.Content) {
		return
	}
	if m.Embeds[0].Footer == nil {
		return
	}
	if !footerRegex.MatchString(m.Embeds[0].Footer.Text) {
		return
	}

	groups := footerRegex.FindStringSubmatch(m.Embeds[0].Footer.Text)

	var (
		sysID     = groups[1]
		memberID  = groups[2]
		userID    discord.UserID
		msgID     discord.MessageID
		channelID discord.ChannelID
	)

	{
		sf, _ := discord.ParseSnowflake(groups[3])
		userID = discord.UserID(sf)
		sf, _ = discord.ParseSnowflake(groups[4])
		msgID = discord.MessageID(sf)
		sf, _ = discord.ParseSnowflake(linkRegex.FindStringSubmatch(m.Content)[1])
		channelID = discord.ChannelID(sf)
	}

	// get full message
	msg, err := bot.State.Message(channelID, msgID)
	if err != nil {
		bot.Sugar.Errorf("Error retrieving original message: %v", err)
		return
	}

	dbMsg := Message{
		MsgID:     msgID,
		UserID:    userID,
		ChannelID: channelID,
		ServerID:  m.GuildID,

		Username: msg.Author.Username,
		Member:   memberID,
		System:   sysID,

		Content: msg.Content,
	}

	err = bot.Insert(dbMsg)
	if err != nil {
		bot.Sugar.Errorf("Error inserting message %v: %v", msgID, err)
	}
}
