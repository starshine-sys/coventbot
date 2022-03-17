package modlog

import (
	"context"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/georgysavva/scany/pgxscan"
	"gitlab.com/1f320/x/duration"
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
	ActionBan  ActionType = "ban"
	ActionKick ActionType = "kick"
	ActionWarn ActionType = "warn"
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

// Channelban logs a channel ban
func (bot *ModLog) Channelban(state *state.State, guildID discord.GuildID, channel discord.ChannelID, userID, modID discord.UserID, reason string) (err error) {
	if reason == "" {
		reason = "N/A"
	}

	entry, err := bot.InsertEntry(guildID, userID, modID, time.Now().UTC(), "channelban", fmt.Sprintf("%v: %v", channel.Mention(), reason))
	if err != nil {
		return err
	}

	ch := bot.logChannelFor(guildID)
	if !ch.IsValid() {
		return
	}

	if len(reason) > 1000 {
		reason = reason[:1000] + "..."
	}

	user, err := state.User(userID)
	if err != nil {
		return err
	}
	mod, err := state.User(modID)
	if err != nil {
		return err
	}

	text := fmt.Sprintf(`**Channel ban for %v | Case %v**
**User:** %v#%v (%v)
**Reason:** %v
**Responsible moderator:** %v#%v (%v)`, channel.Mention(), entry.ID, user.Username, user.Discriminator, entry.UserID, reason, mod.Username, mod.Discriminator, entry.ModID)

	msg, err := state.SendMessage(ch, text)
	if err != nil {
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update mod_log set channel_id = $1, message_id = $2 where id = $3 and server_id = $4", ch, msg.ID, entry.ID, guildID)
	return
}

// Unchannelban logs a channel unban
func (bot *ModLog) Unchannelban(state *state.State, guildID discord.GuildID, channel discord.ChannelID, userID, modID discord.UserID, reason string) (err error) {
	if reason == "" {
		reason = "N/A"
	}

	entry, err := bot.InsertEntry(guildID, userID, modID, time.Now().UTC(), "unchannelban", fmt.Sprintf("%v: %v", channel.Mention(), reason))
	if err != nil {
		return err
	}

	ch := bot.logChannelFor(guildID)
	if !ch.IsValid() {
		return
	}

	if len(reason) > 1000 {
		reason = reason[:1000] + "..."
	}

	user, err := state.User(userID)
	if err != nil {
		return err
	}
	mod, err := state.User(modID)
	if err != nil {
		return err
	}

	text := fmt.Sprintf(`**Channel unban for %v | Case %v**
**User:** %v#%v (%v)
**Reason:** %v
**Responsible moderator:** %v#%v (%v)`, channel.Mention(), entry.ID, user.Username, user.Discriminator, entry.UserID, reason, mod.Username, mod.Discriminator, entry.ModID)

	msg, err := state.SendMessage(ch, text)
	if err != nil {
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update mod_log set channel_id = $1, message_id = $2 where id = $3 and server_id = $4", ch, msg.ID, entry.ID, guildID)
	return
}

// Warn logs a warn
func (bot *ModLog) Warn(state *state.State, guildID discord.GuildID, userID, modID discord.UserID, reason string) (err error) {
	entry, err := bot.InsertEntry(guildID, userID, modID, time.Now().UTC(), "warn", reason)
	if err != nil {
		return err
	}

	ch := bot.logChannelFor(guildID)
	if !ch.IsValid() {
		return
	}

	if len(entry.Reason) > 1000 {
		entry.Reason = entry.Reason[:1000] + "..."
	}

	user, err := state.User(userID)
	if err != nil {
		return err
	}
	mod, err := state.User(modID)
	if err != nil {
		return err
	}

	text := fmt.Sprintf(`**Warn | Case %v**
**User:** %v#%v (%v)
**Reason:** %v
**Responsible moderator:** %v#%v (%v)`, entry.ID, user.Username, user.Discriminator, entry.UserID, entry.Reason, mod.Username, mod.Discriminator, entry.ModID)

	msg, err := state.SendMessage(ch, text)
	if err != nil {
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update mod_log set channel_id = $1, message_id = $2 where id = $3 and server_id = $4", ch, msg.ID, entry.ID, guildID)
	return
}

// Ban logs a ban
func (bot *ModLog) Ban(state *state.State, guildID discord.GuildID, userID, modID discord.UserID, reason string) (err error) {
	entry, err := bot.InsertEntry(guildID, userID, modID, time.Now().UTC(), "ban", reason)
	if err != nil {
		return err
	}

	ch := bot.logChannelFor(guildID)
	if !ch.IsValid() {
		return
	}

	if len(entry.Reason) > 1000 {
		entry.Reason = entry.Reason[:1000] + "..."
	}

	user, err := state.User(userID)
	if err != nil {
		return err
	}
	mod, err := state.User(modID)
	if err != nil {
		return err
	}

	text := fmt.Sprintf(`**Ban | Case %v**
**User:** %v#%v (%v)
**Reason:** %v
**Responsible moderator:** %v#%v (%v)`, entry.ID, user.Username, user.Discriminator, entry.UserID, entry.Reason, mod.Username, mod.Discriminator, entry.ModID)

	msg, err := state.SendMessage(ch, text)
	if err != nil {
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update mod_log set channel_id = $1, message_id = $2 where id = $3 and server_id = $4", ch, msg.ID, entry.ID, guildID)
	return
}

// Unban logs a unban
func (bot *ModLog) Unban(state *state.State, guildID discord.GuildID, userID, modID discord.UserID, reason string) (err error) {
	entry, err := bot.InsertEntry(guildID, userID, modID, time.Now().UTC(), "unban", reason)
	if err != nil {
		return err
	}

	ch := bot.logChannelFor(guildID)
	if !ch.IsValid() {
		return
	}

	if len(entry.Reason) > 1000 {
		entry.Reason = entry.Reason[:1000] + "..."
	}

	user, err := state.User(userID)
	if err != nil {
		return err
	}
	mod, err := state.User(modID)
	if err != nil {
		return err
	}

	text := fmt.Sprintf(`**Unban | Case %v**
**User:** %v#%v (%v)
**Reason:** %v
**Responsible moderator:** %v#%v (%v)`, entry.ID, user.Username, user.Discriminator, entry.UserID, entry.Reason, mod.Username, mod.Discriminator, entry.ModID)

	msg, err := state.SendMessage(ch, text)
	if err != nil {
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update mod_log set channel_id = $1, message_id = $2 where id = $3 and server_id = $4", ch, msg.ID, entry.ID, guildID)
	return
}

// Mute logs a mute
func (bot *ModLog) Mute(state *state.State, guildID discord.GuildID, userID, modID discord.UserID, dur time.Duration, reason string) (err error) {
	entry, err := bot.InsertEntry(guildID, userID, modID, time.Now().UTC(), "mute", reason)
	if err != nil {
		return err
	}

	ch := bot.logChannelFor(guildID)
	if !ch.IsValid() {
		return
	}

	if len(entry.Reason) > 1000 {
		entry.Reason = entry.Reason[:1000] + "..."
	}

	user, err := state.User(userID)
	if err != nil {
		return err
	}
	mod, err := state.User(modID)
	if err != nil {
		return err
	}

	time := "indefinite"
	if dur != 0 {
		time = duration.Format(dur)
	}

	text := fmt.Sprintf(`**Mute | Case %v**
**User:** %v#%v (%v)
**Reason:** %v
**Duration:** %v
**Responsible moderator:** %v#%v (%v)`, entry.ID, user.Username, user.Discriminator, entry.UserID, entry.Reason, time, mod.Username, mod.Discriminator, entry.ModID)

	msg, err := state.SendMessage(ch, text)
	if err != nil {
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update mod_log set channel_id = $1, message_id = $2 where id = $3 and server_id = $4", ch, msg.ID, entry.ID, guildID)
	return
}

// Pause logs a pause
func (bot *ModLog) Pause(state *state.State, guildID discord.GuildID, userID, modID discord.UserID, dur time.Duration, reason string) (err error) {
	entry, err := bot.InsertEntry(guildID, userID, modID, time.Now().UTC(), "pause", reason)
	if err != nil {
		return err
	}

	ch := bot.logChannelFor(guildID)
	if !ch.IsValid() {
		return
	}

	if len(entry.Reason) > 1000 {
		entry.Reason = entry.Reason[:1000] + "..."
	}

	user, err := state.User(userID)
	if err != nil {
		return err
	}
	mod, err := state.User(modID)
	if err != nil {
		return err
	}

	time := "indefinite"
	if dur != 0 {
		time = duration.Format(dur)
	}

	text := fmt.Sprintf(`**Pause | Case %v**
**User:** %v#%v (%v)
**Reason:** %v
**Duration:** %v
**Responsible moderator:** %v#%v (%v)`, entry.ID, user.Username, user.Discriminator, entry.UserID, entry.Reason, time, mod.Username, mod.Discriminator, entry.ModID)

	msg, err := state.SendMessage(ch, text)
	if err != nil {
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update mod_log set channel_id = $1, message_id = $2 where id = $3 and server_id = $4", ch, msg.ID, entry.ID, guildID)
	return
}

// Unmute logs a unmute
func (bot *ModLog) Unmute(state *state.State, guildID discord.GuildID, userID, modID discord.UserID, reason string) (err error) {
	entry, err := bot.InsertEntry(guildID, userID, modID, time.Now().UTC(), "unmute", reason)
	if err != nil {
		return err
	}

	ch := bot.logChannelFor(guildID)
	if !ch.IsValid() {
		return
	}

	if len(entry.Reason) > 1000 {
		entry.Reason = entry.Reason[:1000] + "..."
	}

	user, err := state.User(userID)
	if err != nil {
		return err
	}
	mod, err := state.User(modID)
	if err != nil {
		return err
	}

	text := fmt.Sprintf(`**Unmute | Case %v**
**User:** %v#%v (%v)
**Reason:** %v
**Responsible moderator:** %v#%v (%v)`, entry.ID, user.Username, user.Discriminator, entry.UserID, entry.Reason, mod.Username, mod.Discriminator, entry.ModID)

	msg, err := state.SendMessage(ch, text)
	if err != nil {
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update mod_log set channel_id = $1, message_id = $2 where id = $3 and server_id = $4", ch, msg.ID, entry.ID, guildID)
	return
}

// Unpause logs a unpause
func (bot *ModLog) Unpause(state *state.State, guildID discord.GuildID, userID, modID discord.UserID, reason string) (err error) {
	entry, err := bot.InsertEntry(guildID, userID, modID, time.Now().UTC(), "unpause", reason)
	if err != nil {
		return err
	}

	ch := bot.logChannelFor(guildID)
	if !ch.IsValid() {
		return
	}

	if len(entry.Reason) > 1000 {
		entry.Reason = entry.Reason[:1000] + "..."
	}

	user, err := state.User(userID)
	if err != nil {
		return err
	}
	mod, err := state.User(modID)
	if err != nil {
		return err
	}

	text := fmt.Sprintf(`**Unpause | Case %v**
**User:** %v#%v (%v)
**Reason:** %v
**Responsible moderator:** %v#%v (%v)`, entry.ID, user.Username, user.Discriminator, entry.UserID, entry.Reason, mod.Username, mod.Discriminator, entry.ModID)

	msg, err := state.SendMessage(ch, text)
	if err != nil {
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update mod_log set channel_id = $1, message_id = $2 where id = $3 and server_id = $4", ch, msg.ID, entry.ID, guildID)
	return
}
