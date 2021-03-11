package gatekeeper

import (
	"context"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
)

func (bot *Bot) setPending(p PendingUser) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "insert into gatekeeper (server_id, user_id, key, pending) values ($1, $2, $3, true)", p.ServerID, p.UserID, p.Key)
	return err
}

func (bot *Bot) deletePending(g discord.GuildID, u discord.UserID) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "delete from gatekeeper where user_id = $1 and server_id = $2", u, g)
	return err
}

func (bot *Bot) isPending(g discord.GuildID, u discord.UserID) (b bool) {
	bot.DB.Pool.QueryRow(context.Background(), "select pending from gatekeeper where server_id = $1 and user_id = $2", g, u).Scan(&b)
	return
}

func (bot *Bot) pendingUser(g discord.GuildID, u discord.UserID) (p *PendingUser, err error) {
	p = &PendingUser{}
	err = pgxscan.Get(context.Background(), bot.DB.Pool, p, "select server_id, user_id, key, pending from gatekeeper where server_id = $1 and user_id = $2", g, u)
	if err != nil {
		return nil, err
	}
	return
}

func (bot *Bot) userByUUID(id uuid.UUID) (p *PendingUser, err error) {
	p = &PendingUser{}
	err = pgxscan.Get(context.Background(), bot.DB.Pool, p, "select server_id, user_id, key, pending from gatekeeper where key = $1", id)
	if err != nil {
		return nil, err
	}
	return
}

func (bot *Bot) completeCaptcha(g discord.GuildID, u discord.UserID) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "update gatekeeper set pending = false where server_id = $1 and user_id = $2", g, u)
	return
}

// ServerSettings ...
type ServerSettings struct {
	ID             discord.GuildID
	MemberRole     discord.RoleID
	WelcomeChannel discord.ChannelID
	WelcomeMessage string
}

func (bot *Bot) serverSettings(g discord.GuildID) (s ServerSettings, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &s, "select id, member_role, welcome_channel, welcome_message from servers where id = $1", g)
	return
}

func (bot *Bot) setSettings(s ServerSettings) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "update servers set member_role = $1, welcome_channel = $2, welcome_message = $3 where id = $4", s.MemberRole, s.WelcomeChannel, s.WelcomeMessage, s.ID)
	return err
}
