package admin

import (
	"context"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/georgysavva/scany/pgxscan"
)

// Settings ...
type Settings struct {
	ID     int
	Status gateway.Status

	ActivityType string
	Activity     string
}

// Settings returns the bot settings
func (bot *Bot) Settings() (s Settings) {
	err := pgxscan.Get(context.Background(), bot.DB.Pool, &s, "select * from bot_settings")
	if err != nil {
		bot.Sugar.Errorf("Error: %v", err)
	}
	return
}

// SetSettings sets the bot's settings
func (bot *Bot) SetSettings(s Settings) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), `update bot_settings set
status = $1,
activity_type = $2,
activity = $3`, s.Status, s.ActivityType, s.Activity)
	return
}
