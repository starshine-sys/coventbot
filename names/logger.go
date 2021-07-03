package names

import (
	"context"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/jackc/pgx/v4"
)

func (bot *Bot) nicknameChange(m *gateway.GuildMemberUpdateEvent) {
	nick := m.Nick
	if nick == "" {
		nick = m.User.Username
	}

	var oldNick string
	err := bot.DB.Pool.QueryRow(context.Background(), "select name from nicknames where server_id = $1 and user_id = $2 order by time desc limit 1", m.GuildID, m.User.ID).Scan(&oldNick)
	if err != nil && err != pgx.ErrNoRows {
		bot.Sugar.Errorf("Error getting old nickname: %v", err)
		return
	}

	// if the old nickname is the same as the new nickname, return
	if oldNick == nick {
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "insert into nicknames (server_id, user_id, name) values ($1, $2, $3)", m.GuildID, m.User.ID, nick)
	if err != nil {
		bot.Sugar.Errorf("Error adding to nickname log: %v", err)
		return
	}
}

func (bot *Bot) usernameChange(m *gateway.GuildMemberUpdateEvent) {
	user := m.User.Username + "#" + m.User.Discriminator

	var optedOut bool
	bot.DB.Pool.QueryRow(context.Background(), "select usernames_opt_out from user_config where user_id = $1", m.User.ID).Scan(&optedOut)
	if optedOut {
		return
	}

	var oldUser string
	err := bot.DB.Pool.QueryRow(context.Background(), "select name from usernames where user_id = $1 order by time desc limit 1", m.User.ID).Scan(&oldUser)
	if err != nil && err != pgx.ErrNoRows {
		bot.Sugar.Errorf("Error getting old nickname: %v", err)
		return
	}

	// if the old username is the same as the new username, return
	if oldUser == user {
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "insert into usernames (user_id, name) values ($1, $2)", m.User.ID, user)
	if err != nil {
		bot.Sugar.Errorf("Error adding to username log: %v", err)
		return
	}
}
