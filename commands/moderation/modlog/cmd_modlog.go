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

	err = pgxscan.Select(context.Background(), bot.DB.Pool, &entries, "select * from mod_log where user_id = $1 and server_id = $2 order by id desc", u.ID, ctx.Message.GuildID)
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
			Name: fmt.Sprintf("#%v | %v | %v", entry.ID, entry.ActionType, entry.Time.Format("2006-01-02")),
			Value: fmt.Sprintf(`Responsible moderator: %v#%v
Reason: %v`, mod.Username, mod.Discriminator, entry.Reason),
		})
	}

	embeds := bcr.FieldPaginator("Mod logs", fmt.Sprintf("%v#%v - %v", u.Username, u.Discriminator, u.Mention()), bcr.ColourBlurple, fields, 5)

	_, err = bot.PagedEmbed(ctx, embeds, 10*time.Minute)
	return err
}
