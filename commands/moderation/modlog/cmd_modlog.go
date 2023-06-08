// SPDX-License-Identifier: AGPL-3.0-only
package modlog

import (
	"context"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

func (bot *ModLog) modlog(ctx *bcr.Context) (err error) {
	var entries []*Entry

	u, err := ctx.ParseUser(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("Couldn't find that user.")
		return
	}

	err = pgxscan.Select(context.Background(), bot.DB.Pool, &entries, "select * from mod_log where user_id = $1 and server_id = $2 order by time desc", u.ID, ctx.Message.GuildID)
	if err != nil {
		bot.Sugar.Error(err)
		return bot.Report(ctx, err)
	}

	modCache := map[discord.UserID]*discord.User{}

	var fields []discord.EmbedField

	for _, entry := range entries {
		mod := modCache[entry.ModID]
		if mod == nil {
			mod, err = ctx.State.User(entry.ModID)
			if err != nil {
				return bot.Report(ctx, err)
			}
			modCache[entry.ModID] = mod
		}

		fields = append(fields, discord.EmbedField{
			Name: fmt.Sprintf("#%v | %v | <t:%v>", entry.ID, entry.ActionType, entry.Time.Unix()),
			Value: fmt.Sprintf(`Responsible moderator: %v
Reason: %v`, mod.Tag(), entry.Reason),
		})
	}

	embeds := bcr.FieldPaginator("Mod logs", fmt.Sprintf("%v - %v", u.Tag(), u.Mention()), bcr.ColourBlurple, fields, 5)

	_, err = bot.PagedEmbed(ctx, embeds, 10*time.Minute)
	return err
}
