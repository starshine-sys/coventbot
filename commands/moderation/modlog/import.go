package modlog

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jackc/pgx/v4"
	"github.com/starshine-sys/bcr"
)

func (bot *ModLog) cmdImport(ctx *bcr.Context) (err error) {
	if len(ctx.Message.Attachments) == 0 {
		_, err = ctx.Send("No file to import attached.")
		return
	}

	// check count
	var currentCount int
	err = bot.DB.Pool.QueryRow(context.Background(), "select count(*) from mod_log where server_id = $1", ctx.Message.GuildID).Scan(&currentCount)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if currentCount > 0 {
		yes, timeout := ctx.ConfirmButton(ctx.Author.ID, bcr.ConfirmData{
			Message:   "⚠️ There are existing mod logs for this server, which will be deleted if you proceed with this import. Are you sure you want to proceed?",
			YesPrompt: "Yes, clear data and proceed",
			YesStyle:  discord.DangerButtonStyle(),
		})
		if !yes || timeout {
			_, err = ctx.Send("Import cancelled.")
			return err
		}

		ct, err := bot.DB.Pool.Exec(context.Background(), "delete from mod_log where server_id = $1", ctx.Message.GuildID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		_, err = ctx.Sendf("Cleared mod logs, %v entries deleted.", ct.RowsAffected())
		if err != nil {
			return err
		}
	}

	resp, err := http.Get(ctx.Message.Attachments[0].URL)
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
		if !yes || !timeout {
			_, err = ctx.Send("Import cancelled.")
			return err
		}
	}

	// do the import
	count, err := bot.DB.Pool.CopyFrom(
		context.Background(),
		pgx.Identifier{"mod_log"},
		[]string{"id", "server_id", "user_id", "mod_id", "action_type", "reason", "time"},
		pgx.CopyFromSlice(len(ex.Entries), func(i int) ([]interface{}, error) {
			return []interface{}{ex.Entries[i].ID, ex.ServerID, ex.Entries[i].UserID, ex.Entries[i].ModID, ex.Entries[i].ActionType, ex.Entries[i].Reason, ex.Entries[i].Time}, nil
		}),
	)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Sendf("Success! Imported %v mod log entr%v.", count, func(b bool) string {
		if b {
			return "y"
		}
		return "ies"
	}(count == 1))
	return
}
