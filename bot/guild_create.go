package bot

import "github.com/diamondburned/arikawa/v2/gateway"

// GuildCreate logs the bot joining a server, and creates a database entry if one doesn't exist
func (bot *Bot) GuildCreate(g *gateway.GuildCreateEvent) {
	// create the server if it doesn't exist
	exists, err := bot.DB.CreateServerIfNotExists(g.ID)
	// if the server exists, don't log the join
	if exists {
		return
	}
	if err != nil {
		bot.Sugar.Errorf("Error creating database entry for server: %v", err)
		return
	}

	bot.Sugar.Infof("Joined server %v (%v).", g.Name, g.ID)
	return
}
