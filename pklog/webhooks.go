package pklog

import (
	"errors"

	"github.com/diamondburned/arikawa/v3/discord"
)

// ErrNotExists ...
var ErrNotExists = errors.New("webhooks not found in cache")

// Webhooks ...
type Webhooks struct {
	MessageWebhookID    discord.WebhookID
	MessageWebhookToken string
}

// SetWebhooks ...
func (bot *Bot) SetWebhooks(id discord.GuildID, w *Webhooks) {
	bot.WebhookCache.Set(id.String(), w)
}

// GetWebhooks ...
func (bot *Bot) GetWebhooks(id discord.GuildID) (*Webhooks, error) {
	v, err := bot.WebhookCache.Get(id.String())
	if err != nil {
		return nil, ErrNotExists
	}
	if _, ok := v.(*Webhooks); !ok {
		return nil, errors.New("could not convert interface to Webhooks")
	}

	return v.(*Webhooks), nil
}

// ResetCache ...
func (bot *Bot) ResetCache(id discord.GuildID) {
	bot.WebhookCache.Remove(id.String())
}
