package types

import "github.com/diamondburned/arikawa/v2/discord"

// BotConfig ...
type BotConfig struct {
	Token       string
	DatabaseURL string `yaml:"database_url"`

	Prefixes []string
	Owners   []discord.UserID
}
