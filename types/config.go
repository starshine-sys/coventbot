package types

import "github.com/diamondburned/arikawa/v2/discord"

// BotConfig ...
type BotConfig struct {
	Token       string
	DatabaseURL string `yaml:"database_url"`

	Prefixes []string
	Owners   []discord.UserID

	GuildLogWebhook *Webhook `yaml:"guild_log"`
}

// Webhook is a single webhook config
type Webhook struct {
	ID    discord.WebhookID
	Token string
}
