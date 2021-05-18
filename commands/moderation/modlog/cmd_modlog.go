package modlog

import (
	"context"
	"fmt"
	"math"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

func (bot *ModLog) modlog(ctx *bcr.Context) (err error) {
	var entries []*Entry

	u, err := ctx.ParseUser(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("Couldn't find that user.", nil)
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
			mod, err = bot.State.User(entry.ModID)
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

	embeds := FieldPaginator("Mod logs", fmt.Sprintf("%v#%v - %v", u.Username, u.Discriminator, u.Mention()), bcr.ColourBlurple, fields, 5)

	_, err = ctx.PagedEmbed(embeds, false)
	return err
}

// FieldPaginator paginates embed fields, for use in ctx.PagedEmbed
func FieldPaginator(title, description string, colour discord.Color, fields []discord.EmbedField, perPage int) []discord.Embed {
	var (
		embeds []discord.Embed
		count  int

		pages = 1
		buf   = discord.Embed{
			Title:       title,
			Description: description,
			Color:       colour,
			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("Page 1/%v", math.Ceil(float64(len(fields))/float64(perPage))),
			},
		}
	)

	for _, field := range fields {
		if count >= perPage {
			embeds = append(embeds, buf)
			buf = discord.Embed{
				Title:       title,
				Description: description,
				Color:       colour,
				Footer: &discord.EmbedFooter{
					Text: fmt.Sprintf("Page %v/%v", pages+1, math.Ceil(float64(len(fields))/float64(perPage))),
				},
			}
			count = 0
			pages++
		}
		buf.Fields = append(buf.Fields, field)
		count++
	}

	embeds = append(embeds, buf)

	return embeds
}
