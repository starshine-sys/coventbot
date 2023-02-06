// SPDX-License-Identifier: AGPL-3.0-only
package notes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) exportNotes(ctx *bcr.Context) (err error) {
	start := time.Now()

	var notes []ExportNote
	err = pgxscan.Select(
		context.Background(), bot.DB.Pool, &notes,
		"select user_id, note, moderator as moderator_id, created as timestamp from notes where server_id = $1", ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.SetIndent("", "  ")

	err = enc.Encode(notes)
	if err != nil {
		return bot.Report(ctx, err)
	}

	t := time.Since(start)

	_, err = ctx.State.SendMessageComplex(ctx.Message.ChannelID, api.SendMessageData{
		Content: fmt.Sprintf("Here you go! Exported %v notes in %v", len(notes), t.Round(time.Millisecond)),
		Files: []sendpart.File{{
			Name:   "export.json",
			Reader: b,
		}},
	})
	return err
}
