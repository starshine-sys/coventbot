package modlog

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

func (bot *ModLog) cmdImport(ctx *bcr.Context) (err error) {
	if len(ctx.Message.Attachments) == 0 && ctx.RawArgs == "" {
		_, err = ctx.Send("No file to import attached.")
		return
	}

	url := ctx.RawArgs
	if len(ctx.Message.Attachments) > 0 {
		url = ctx.Message.Attachments[0].URL
	}

	resp, err := http.Get(url)
	if err != nil {
		return bot.Report(ctx, err)
	}
	defer resp.Body.Close()

	var ex Export

	err = json.NewDecoder(resp.Body).Decode(&ex)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(ex.Entries) == 0 {
		_, err = ctx.Send("The export did not contain any entries.")
		return
	}

	if ctx.Message.GuildID != ex.ServerID {
		yes, timeout := ctx.ConfirmButton(ctx.Author.ID, bcr.ConfirmData{
			Message:   "The server ID in the export doesn't match this server's ID. Do you want to continue anyway?",
			YesPrompt: "Continue",
		})
		if !yes || timeout {
			_, err = ctx.Send("Import cancelled.")
			return err
		}
	}

	var currentEntries []Entry
	err = pgxscan.Select(context.Background(), bot.DB.Pool, &currentEntries, "select * from mod_log where server_id = $1", ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}
	contains := func(id int64, reason string) bool {
		for _, e := range currentEntries {
			if e.ID == id && e.Reason == reason {
				return true
			}
		}
		return false
	}

	tx, err := bot.DB.Pool.Begin(context.Background())
	if err != nil {
		return bot.Report(ctx, err)
	}

	// do the import
	var done int

	for _, e := range ex.Entries {
		if contains(e.ID, e.Reason) {
			continue
		}

		_, err = tx.Exec(context.Background(), insertSql, ctx.Message.GuildID, e.UserID, e.ModID, e.ActionType, e.Reason, e.Time)
		if err != nil {
			tx.Rollback(context.Background())

			return bot.Report(ctx, err)
		}
		done++
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Sendf("Success! Imported %v mod log entr%v.", done, func(b bool) string {
		if b {
			return "y"
		}
		return "ies"
	}(done == 1))
	return
}
