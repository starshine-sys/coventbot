package types

import "github.com/diamondburned/arikawa/v3/discord"

// BotConfig ...
type BotConfig struct {
	Token        string
	DatabaseURL  string `yaml:"database_url"`
	SentryURL    string `yaml:"sentry_url"`
	DebugLogging bool   `yaml:"debug_logging"`

	AESKey string `yaml:"aes_key"`

	Prefixes []string
	Owners   []discord.UserID

	GlobalQuotes bool `yaml:"global_quotes"`

	SyncCommandsIn []discord.GuildID `yaml:"sync_commands_in"`

	GuildLogWebhook Webhook `yaml:"guild_log"`

	VerifyListen    string `yaml:"verify_listen"`
	VerifyBaseURL   string `yaml:"verify_base_url"`
	HCaptchaSitekey string `yaml:"hcaptcha_sitekey"`
	HCaptchaSecret  string `yaml:"hcaptcha_secret"`

	Branding struct {
		Name string

		Private  bool
		PublicID discord.UserID `yaml:"public_id"`

		SupportServer string `yaml:"support_server"`
	}

	DMs struct {
		Open    bool
		Webhook Webhook

		BlockedUsers []discord.UserID `yaml:"blocked_users"`
	} `yaml:"dms"`
}

// Webhook is a single webhook config
type Webhook struct {
	ID    discord.WebhookID
	Token string
}
