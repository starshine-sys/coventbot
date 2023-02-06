// SPDX-License-Identifier: AGPL-3.0-only
package modlog

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

// Export is an export of a server's mod logs
type Export struct {
	ServerID  discord.GuildID `json:"server_id"`
	Timestamp time.Time       `json:"timestamp"`

	Entries []Entry `json:"entries"`
}

func (bot *ModLog) export(ctx *bcr.Context) (err error) {
	sql := "select * from mod_log where server_id = $1"
	if len(ctx.Args) > 0 {
		sql += " and user_id = $2"
	}
	sql += " order by id"

	ex := Export{
		ServerID:  ctx.Message.GuildID,
		Timestamp: time.Now().UTC(),
	}

	if len(ctx.Args) > 0 {
		u, err := ctx.ParseUser(ctx.RawArgs)
		if err != nil {
			_, err = ctx.Send("Couldn't find that user.")
			return err
		}
		err = pgxscan.Select(context.Background(), bot.DB.Pool, &ex.Entries, sql, ctx.Message.GuildID, u.ID)
	} else {
		err = pgxscan.Select(context.Background(), bot.DB.Pool, &ex.Entries, sql, ctx.Message.GuildID)
	}
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(ex.Entries) == 0 {
		_, err = ctx.Send("There are no mod logs to export.")
		return
	}

	b, err := json.MarshalIndent(&ex, "", "    ")
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.NewMessage().AddFile("export.json", bytes.NewReader(b)).Send()
	return
}
