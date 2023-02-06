// SPDX-License-Identifier: AGPL-3.0-only
package db

import (
	"context"
	"errors"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
)

// Note is a note
type Note struct {
	ID uint64

	ServerID discord.GuildID
	UserID   discord.UserID

	Note      string
	Moderator discord.UserID
	Created   time.Time
}

// NewNote ...
func (db *DB) NewNote(n Note) (Note, error) {
	if n.ServerID == 0 || n.UserID == 0 || n.Note == "" || n.Moderator == 0 || n.Created.IsZero() {
		return n, errors.New("one or more required fields was empty")
	}

	err := db.Pool.QueryRow(context.Background(), "insert into notes (server_id, user_id, note, moderator, created) values ($1, $2, $3, $4, $5) returning id", n.ServerID, n.UserID, n.Note, n.Moderator, n.Created).Scan(&n.ID)
	return n, err
}

// ...
var (
	ErrNoteNotFound = errors.New("note not found")
)

// DelNote ...
func (db *DB) DelNote(guildID discord.GuildID, id uint64) (err error) {
	ct, err := db.Pool.Exec(context.Background(), "delete from notes where server_id = $1 and id = $2", guildID, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNoteNotFound
	}
	return nil
}

// UserNotes ...
func (db *DB) UserNotes(guildID discord.GuildID, userID discord.UserID) (notes []Note, err error) {
	err = pgxscan.Select(context.Background(), db.Pool, &notes, "select * from notes where server_id = $1 and user_id = $2 order by created desc", guildID, userID)
	return
}
