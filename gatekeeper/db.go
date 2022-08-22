package gatekeeper

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
)

func (bot *Bot) setPending(guildID discord.GuildID, userID discord.UserID) (p PendingUser, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &p, `insert into gatekeeper
	(guild_id, user_id, key, pending) values ($1, $2, $3, true)
	on conflict (guild_id, user_id) do update
	set key = $3
	returning *`, guildID, userID, uuid.New())
	return p, err
}

func (bot *Bot) deletePending(g discord.GuildID, u discord.UserID) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "delete from gatekeeper where user_id = $1 and guild_id = $2", u, g)
	return err
}

func (bot *Bot) userByUUID(id uuid.UUID) (p *PendingUser, err error) {
	p = &PendingUser{}
	err = pgxscan.Get(context.Background(), bot.DB.Pool, p, "select guild_id, user_id, key, pending from gatekeeper where key = $1", id)
	if err != nil {
		return nil, err
	}
	return
}

func (bot *Bot) completeCaptcha(g discord.GuildID, u discord.UserID) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "update gatekeeper set pending = false where guild_id = $1 and user_id = $2", g, u)
	return
}

// ServerSettings ...
type ServerSettings struct {
	ID             discord.GuildID
	MemberRole     discord.RoleID
	WelcomeChannel discord.ChannelID
	WelcomeMessage string
	GatekeeperLog  discord.ChannelID
}

func (bot *Bot) serverSettings(g discord.GuildID) (s ServerSettings, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &s, "select id, member_role, welcome_channel, welcome_message, gatekeeper_log from servers where id = $1", g)
	return
}

func (bot *Bot) setSettings(s ServerSettings) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set member_role = $1, welcome_channel = $2, welcome_message = $3, gatekeeper_log = $4 where id = $5", s.MemberRole, s.WelcomeChannel, s.WelcomeMessage, s.GatekeeperLog, s.ID)
	return err
}
