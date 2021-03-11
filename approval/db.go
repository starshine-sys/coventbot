package approval

import (
	"context"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
)

// ServerSettings ...
type ServerSettings struct {
	ID discord.GuildID

	ApproveRemoveRoles    []uint64
	ApproveAddRoles       []uint64
	ApproveWelcomeChannel discord.ChannelID
	ApproveWelcomeMessage string
}

func (bot *Bot) serverSettings(g discord.GuildID) (s ServerSettings, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &s, "select id, approve_remove_roles, approve_add_roles, approve_welcome_channel, approve_welcome_message from servers where id = $1", g)
	return
}

func (bot *Bot) setSettings(s ServerSettings) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set approve_remove_roles = $1, approve_add_roles = $2, approve_welcome_channel = $3, approve_welcome_message = $4 where id = $5", s.ApproveRemoveRoles, s.ApproveAddRoles, s.ApproveWelcomeChannel, s.ApproveWelcomeMessage, s.ID)
	return err
}
