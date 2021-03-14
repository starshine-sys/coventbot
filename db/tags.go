package db

import (
	"context"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/snowflake/v2"
)

// Tag is a single tag
type Tag struct {
	ID       snowflake.SmallID
	ServerID discord.GuildID

	Name     string
	Response string

	CreatedBy discord.UserID
	CreatedAt time.Time
}

// AddTag adds a tag with the given name and response to the database
func (db *DB) AddTag(ctx *bcr.Context, tag Tag) (t Tag, err error) {
	err = db.Pool.QueryRow(context.Background(), "insert into tags (id, server_id, name, response, created_by) values ($1, $2, $3, $4, $5) returning id, created_at", sfGen.Get(), ctx.Message.GuildID, tag.Name, tag.Response, ctx.Author.ID).Scan(&tag.ID, &tag.CreatedAt)
	return tag, err
}

// GetTag gets a tag by name
func (db *DB) GetTag(guildID discord.GuildID, s string) (t *Tag, err error) {
	t = &Tag{}

	err = pgxscan.Get(context.Background(), db.Pool, t, "select id, server_id, name, response, created_by, created_at from tags where lower(name) = lower($1) and server_id = $2", s, guildID)
	return t, err
}

// Tags returns all tags for the given server
func (db *DB) Tags(guildID discord.GuildID) (t []*Tag, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &t, "select id, server_id, name, response, created_by, created_at from tags where server_id = $1 order by name, id", guildID)
	return
}
