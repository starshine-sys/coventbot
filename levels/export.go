// SPDX-License-Identifier: AGPL-3.0-only
package levels

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

func (bot *Bot) exportLevels(ctx *bcr.Context) (err error) {
	start := time.Now()

	var lvls []importedLevel
	err = pgxscan.Select(
		context.Background(), bot.DB.Pool, &lvls,
		"select user_id, xp, colour from levels where server_id = $1", ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.SetIndent("", "  ")

	err = enc.Encode(lvls)
	if err != nil {
		return bot.Report(ctx, err)
	}

	t := time.Since(start)

	_, err = ctx.State.SendMessageComplex(ctx.Message.ChannelID, api.SendMessageData{
		Content: fmt.Sprintf("Here you go! Exported %v members in %v", len(lvls), t.Round(time.Millisecond)),
		Files: []sendpart.File{{
			Name:   "export.json",
			Reader: b,
		}},
	})
	return err
}
