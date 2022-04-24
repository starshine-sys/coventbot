package modlog

import (
	"context"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

// Entry ...
type Entry struct {
	ID       int64           `json:"id"`
	ServerID discord.GuildID `json:"-"`

	UserID discord.UserID `json:"user_id"`
	ModID  discord.UserID `json:"mod_id"`

	ActionType ActionType `json:"action_type"`
	Reason     string     `json:"reason,omitempty"`

	Time time.Time `json:"timestamp"`

	ChannelID discord.ChannelID `json:"-"`
	MessageID discord.MessageID `json:"-"`
}

func (bot *ModLog) logChannelFor(guildID discord.GuildID) (logChannel discord.ChannelID) {
	bot.DB.Pool.QueryRow(context.Background(), "select mod_log_channel from servers where id = $1", guildID).Scan(&logChannel)
	return
}

// ActionType ...
type ActionType string

// Constants for action types
const (
	ActionBan          ActionType = "ban"
	ActionUnban        ActionType = "unban"
	ActionKick         ActionType = "kick"
	ActionWarn         ActionType = "warn"
	ActionChannelban   ActionType = "channelban"
	ActionUnchannelban ActionType = "unchannelban"
)

// InsertEntry inserts a mod log entry
func (bot *ModLog) InsertEntry(guildID discord.GuildID, user, mod discord.UserID, timestamp time.Time, actionType ActionType, reason string) (log *Entry, err error) {
	if reason == "" {
		reason = "N/A"
	}

	log = &Entry{}

	err = pgxscan.Get(context.Background(), bot.DB.Pool, log, `insert into mod_log
	(id, server_id, user_id, mod_id, action_type, reason, time)
	values
	(
			coalesce(
					(select id + 1 from mod_log where server_id = $1 order by id desc limit 1), 1
			),
			$1, $2, $3, $4, $5, $6
	)
	returning *`, guildID, user, mod, actionType, reason, timestamp)
	if err != nil {
		return nil, err
	}
	return
}

func (bot *ModLog) Log(
	s *state.State,
	actionType ActionType,
	guildID discord.GuildID,
	userID, modID discord.UserID,
	reason string,
) (err error) {
	if reason == "" {
		reason = "N/A"
	}

	entry, err := bot.InsertEntry(guildID, userID, modID, time.Now().UTC(), actionType, reason)
	if err != nil {
		return err
	}

	ch := bot.logChannelFor(guildID)
	if !ch.IsValid() {
		return
	}

	e := bot.Embed(s, entry)
	msg, err := s.SendEmbeds(ch, e)
	if err != nil {
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update mod_log set channel_id = $1, message_id = $2 where id = $3 and server_id = $4", ch, msg.ID, entry.ID, guildID)
	return
}

func (bot *ModLog) Embed(s *state.State, entry *Entry) (embed discord.Embed) {
	reason := entry.Reason
	if len(reason) > 1000 {
		reason = reason[:1000] + "..."
	}

	var u, mod discord.User
	if user, err := bot.Member(entry.ServerID, entry.UserID); err == nil {
		u = user.User
	} else if user, err := s.User(entry.UserID); err == nil {
		u = *user
	} else {
		u = discord.User{Username: "unknown", Discriminator: "0000", ID: entry.UserID}
	}

	if user, err := bot.Member(entry.ServerID, entry.UserID); err == nil {
		mod = user.User
	} else if user, err := s.User(entry.UserID); err == nil {
		mod = *user
	} else {
		mod = discord.User{Username: "unknown", Discriminator: "0000", ID: entry.UserID}
	}

	e := discord.Embed{
		Title:       fmt.Sprintf("%s | case %d", entry.ActionType, entry.ID),
		Description: fmt.Sprintf("**User:** %v <@!%v>\n**Reason:** %v\n**Moderator:** %v (%v)", u.Tag(), entry.UserID, entry.Reason, mod.Tag(), entry.ModID),
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v", entry.UserID),
		},
		Timestamp: discord.Timestamp(entry.Time),
	}

	switch entry.ActionType {
	case ActionBan, ActionChannelban, ActionKick:
		e.Color = bcr.ColourRed
	case ActionUnban, ActionUnchannelban:
		e.Color = bcr.ColourGreen
	case ActionWarn:
		e.Color = bcr.ColourOrange
	default:
		e.Color = bcr.ColourBlurple
	}

	return e
}
