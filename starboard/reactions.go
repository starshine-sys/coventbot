package starboard

import (
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

// MessageReactionAdd ...
func (bot *Bot) MessageReactionAdd(ev *gateway.MessageReactionAddEvent) {
	if bot.mu[ev.MessageID] == nil {
		bot.mu[ev.MessageID] = &sync.Mutex{}
	}
	bot.mu[ev.MessageID].Lock()
	defer bot.mu[ev.MessageID].Unlock()

	settings, err := bot.DB.Starboard(ev.GuildID)
	if err != nil {
		bot.Sugar.Errorf("Error getting starboard settings: %v", err)
		return
	}

	// if the emoji isn't the starboard emoji, return
	if ev.Emoji.String() != settings.StarboardEmoji && ev.Emoji.Name != settings.StarboardEmoji {
		return
	}

	// if the channel is blacklisted, return
	if bot.DB.IsBlacklisted(ev.GuildID, ev.ChannelID) {
		return
	}

	bot.reactionInner(ev.UserID, ev.ChannelID, ev.MessageID, ev.Emoji, ev.GuildID)
}

// MessageReactionDelete ...
func (bot *Bot) MessageReactionDelete(ev *gateway.MessageReactionRemoveEvent) {
	if bot.mu[ev.MessageID] == nil {
		bot.mu[ev.MessageID] = &sync.Mutex{}
	}
	bot.mu[ev.MessageID].Lock()
	defer bot.mu[ev.MessageID].Unlock()

	// if the channel is blacklisted, return
	if bot.DB.IsBlacklisted(ev.GuildID, ev.ChannelID) {
		return
	}

	bot.reactionInner(ev.UserID, ev.ChannelID, ev.MessageID, ev.Emoji, ev.GuildID)
}

// MessageReactionRemoveEmoji ...
func (bot *Bot) MessageReactionRemoveEmoji(ev *gateway.MessageReactionRemoveEmojiEvent) {
	if bot.mu[ev.MessageID] == nil {
		bot.mu[ev.MessageID] = &sync.Mutex{}
	}
	bot.mu[ev.MessageID].Lock()
	defer bot.mu[ev.MessageID].Unlock()

	settings, err := bot.DB.Starboard(ev.GuildID)
	if err != nil {
		bot.Sugar.Errorf("Error getting starboard settings: %v", err)
		return
	}

	if ev.Emoji.String() != settings.StarboardEmoji && ev.Emoji.Name != settings.StarboardEmoji {
		return
	}

	// if the channel is blacklisted, return
	if bot.DB.IsBlacklisted(ev.GuildID, ev.ChannelID) {
		return
	}

	star, err := bot.DB.StarboardMessage(ev.MessageID)
	if err != nil {
		bot.Sugar.Errorf("Error getting database entry for message: %v", err)
		return
	}

	if star != nil {
		bot.deleteMessage(ev.ChannelID, ev.MessageID, settings, star)
	}
}

// MessageReactionRemoveAll ...
func (bot *Bot) MessageReactionRemoveAll(ev *gateway.MessageReactionRemoveAllEvent) {
	if bot.mu[ev.MessageID] == nil {
		bot.mu[ev.MessageID] = &sync.Mutex{}
	}
	bot.mu[ev.MessageID].Lock()
	defer bot.mu[ev.MessageID].Unlock()

	settings, err := bot.DB.Starboard(ev.GuildID)
	if err != nil {
		bot.Sugar.Errorf("Error getting starboard settings: %v", err)
		return
	}

	star, err := bot.DB.StarboardMessage(ev.MessageID)
	if err != nil {
		bot.Sugar.Errorf("Error getting database entry for message: %v", err)
		return
	}

	// if the channel is blacklisted, return
	if bot.DB.IsBlacklisted(ev.GuildID, ev.ChannelID) {
		return
	}

	if star != nil {
		bot.deleteMessage(ev.ChannelID, ev.MessageID, settings, star)
	}
}

func (bot *Bot) reactionInner(userID discord.UserID, channelID discord.ChannelID, messageID discord.MessageID, emoji discord.Emoji, guildID discord.GuildID) {
	// if it's a DM channel return
	if !guildID.IsValid() {
		return
	}

	settings, err := bot.DB.Starboard(guildID)
	if err != nil {
		bot.Sugar.Errorf("Error getting starboard settings: %v", err)
		return
	}

	// if no starboard channel is set, return
	if !settings.StarboardChannel.IsValid() {
		return
	}

	// make sure we're not starring our own starboard messages
	if channelID == settings.StarboardChannel {
		if bot.Router.Bot == nil {
			bot.Sugar.Info("Reaction event received, but bot user is not cached yet.")
			return
		}

		m, err := bot.State.Message(channelID, messageID)
		if err != nil {
			bot.Sugar.Errorf("Error getting message: %v", err)
			return
		}

		if m.Author.ID == bot.Router.Bot.ID {
			return
		}
	}

	m, err := bot.State.Message(channelID, messageID)
	if err != nil {
		bot.Sugar.Errorf("Error getting message: %v", err)
		return
	}

	// set count and break out of loop, hack to delete starboard message if all reactions were removed
	var count int
	for _, r := range m.Reactions {
		if r.Emoji.String() == settings.StarboardEmoji || r.Emoji.Name == settings.StarboardEmoji {
			count = r.Count
			break
		}
	}

	star, err := bot.DB.StarboardMessage(m.ID)
	if err != nil && errors.Cause(err) != pgx.ErrNoRows {
		bot.Sugar.Errorf("Error getting database entry for message: %v", err)
		return
	}

	if count < settings.StarboardLimit && (star != nil && star.MessageID != 0) {
		bot.deleteMessage(channelID, messageID, settings, star)
	}
	if count >= settings.StarboardLimit {
		bot.starboardMessage(*m, settings, star, count)
	}

	return
}
