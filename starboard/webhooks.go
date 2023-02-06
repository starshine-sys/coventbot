// SPDX-License-Identifier: AGPL-3.0-only
package starboard

import (
	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
)

func (bot *Bot) Webhook(id discord.WebhookID, token string) *webhook.Client {
	bot.webhooksMu.Lock()
	defer bot.webhooksMu.Unlock()

	client, ok := bot.webhooks[id]
	if !ok {
		s, _ := bot.Router.StateFromGuildID(0)

		client := webhook.FromAPI(id, token, s.Client)
		bot.webhooks[id] = client
		return client
	}
	return client
}

func (bot *Bot) RemoveWebhook(id discord.WebhookID) {
	bot.webhooksMu.Lock()
	delete(bot.webhooks, id)
	bot.webhooksMu.Unlock()
}
