package highlights

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
)

// Highlight ...
type Highlight struct {
	UserID   discord.UserID
	ServerID discord.GuildID

	Highlights []string
	Blocked    []uint64
}

// HighlightConfig ...
type HighlightConfig struct {
	ServerID discord.GuildID

	HlEnabled bool
	Blocked   []uint64
}

func (bot *Bot) hlConfig(guildID discord.GuildID) (conf HighlightConfig, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &conf, "insert into highlight_config (server_id) values ($1) on conflict (server_id) do update set server_id = $1 returning *", guildID)
	return conf, err
}

func (bot *Bot) setConfig(conf HighlightConfig) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "insert into highlight_config (server_id, hl_enabled, blocked) values ($1, $2, $3) on conflict (server_id) do update set hl_enabled = $2, blocked = $3", conf.ServerID, conf.HlEnabled, conf.Blocked)
	return
}

func (bot *Bot) guildHighlights(guildID discord.GuildID) (hls []Highlight, err error) {
	err = pgxscan.Select(context.Background(), bot.DB.Pool, &hls, "select * from highlights where server_id = $1", guildID)
	return
}

func (bot *Bot) userHighlights(guildID discord.GuildID, userID discord.UserID) (hl Highlight, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &hl, "insert into highlights (user_id, server_id) values ($1, $2) on conflict (user_id, server_id) do update set server_id = $2 returning *", userID, guildID)
	return hl, err
}

func (bot *Bot) setUserHighlights(hl Highlight) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), `insert into highlights
	(user_id, server_id, highlights, blocked) values ($1, $2, $3, $4) on conflict (user_id, server_id) do update set highlights = $3, blocked = $4`, hl.UserID, hl.ServerID, hl.Highlights, hl.Blocked)
	return
}
