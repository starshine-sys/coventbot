package todos

import (
	"context"
	"errors"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
)

// Todo is a single todo
type Todo struct {
	ID     int64
	UserID discord.UserID

	Description string

	OrigMID       discord.MessageID
	OrigChannelID discord.ChannelID
	OrigServerID  discord.GuildID

	MID       discord.MessageID
	ChannelID discord.ChannelID
	ServerID  discord.GuildID

	Complete bool

	Created   time.Time
	Completed *time.Time
}

func (bot *Bot) newTodo(t Todo) (out *Todo, err error) {
	if t.UserID == 0 || t.Description == "" {
		return nil, errors.New("required fields empty")
	}

	err = bot.DB.Pool.QueryRow(context.Background(), `insert into todos
	(user_id, description, orig_mid, orig_channel_id, orig_server_id, mid, channel_id, server_id)
	values ($1, $2, $3, $4, $5, $6, $7, $8) returning id, created`, t.UserID, t.Description, t.OrigMID, t.OrigChannelID, t.OrigServerID, t.MID, t.ChannelID, t.ServerID).Scan(&t.ID, &t.Created)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (bot *Bot) getTodo(id int64, user discord.UserID) (*Todo, error) {
	t := &Todo{}

	err := pgxscan.Get(context.Background(), bot.DB.Pool, t, "select * from todos where id = $1 and user_id = $2", id, user)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (bot *Bot) getTodoMessage(id discord.MessageID, user discord.UserID) (*Todo, error) {
	t := &Todo{}

	err := pgxscan.Get(context.Background(), bot.DB.Pool, t, "select * from todos where mid = $1 and user_id = $2", id, user)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (bot *Bot) markComplete(id int64) (err error) {
	t := time.Now().UTC()

	_, err = bot.DB.Pool.Exec(context.Background(), "update todos set complete = true, completed = $1 where id = $2", &t, id)
	return
}
