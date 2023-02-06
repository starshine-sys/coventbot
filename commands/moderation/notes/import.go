// SPDX-License-Identifier: AGPL-3.0-only
package notes

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

type ExportNote struct {
	UserID      discord.UserID `json:"user_id"`
	ModeratorID discord.UserID `json:"moderator_id"`
	Note        string         `json:"note"`
	Timestamp   time.Time      `json:"timestamp"`
}

func (bot *Bot) importNotes(ctx *bcr.Context) (err error) {
	if len(ctx.Message.Attachments) < 1 && len(ctx.Args) < 1 {
		return ctx.SendX("Please attach a JSON export, or provide a link to one.")
	}
	url := ctx.RawArgs
	if len(ctx.Message.Attachments) > 0 && strings.HasSuffix(ctx.Message.Attachments[0].URL, ".json") {
		url = ctx.Message.Attachments[0].URL
	}

	resp, err := http.Get(url)
	if err != nil {
		bot.Sugar.Errorf("downloading export file: %v", err)
		return bot.Report(ctx, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return ctx.SendfX("Encountered an error code while downloading the export file: %v\nPlease try again.", resp.StatusCode)
	}

	notes := make([]ExportNote, 0)
	err = json.NewDecoder(resp.Body).Decode(&notes)
	if err != nil {
		bot.Sugar.Errorf("decoding export file: %v", err)
		return ctx.SendfX("There was an error decoding the export file: %v", bcr.AsCode(err.Error()))
	}

	yes, timeout := ctx.ConfirmButton(ctx.Author.ID, bcr.ConfirmData{
		Message:   "Are you sure you want to import notes? This will overwrite **all** existing notes!",
		YesPrompt: "Yes, overwrite data",
		NoPrompt:  "Cancel",
	})
	if timeout {
		return ctx.SendX("Prompt timed out.")
	}
	if !yes {
		return ctx.SendX("Cancelled.")
	}

	err = bot.importNotesDB(ctx.Message.GuildID, notes)
	if err != nil {
		return bot.Report(ctx, err)
	}

	return ctx.SendfX("Success! Imported %v notes.", len(notes))
}

func (bot *Bot) importNotesDB(guildID discord.GuildID, notes []ExportNote) error {
	ctx := context.Background()

	tx, err := bot.DB.Pool.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "beginning transaction")
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "delete from notes where server_id = $1", guildID)
	if err != nil {
		return errors.Wrap(err, "deleting existing rows")
	}

	for _, note := range notes {
		_, err = tx.Exec(ctx, `insert into notes (server_id, user_id, note, moderator, created)
		values ($1, $2, $3, $4, $5)`, guildID, note.UserID, note.Note, note.ModeratorID, note.Timestamp)
		if err != nil {
			return errors.Wrap(err, "inserting note row")
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "committing transaction")
	}
	return nil
}
