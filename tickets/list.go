package tickets

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) list(ctx *bcr.Context) (err error) {
	categories := []struct {
		CategoryID discord.ChannelID
		ServerID   discord.GuildID
		LogChannel discord.ChannelID

		Name        string
		Mention     string
		Description string

		Count int

		PerUserLimit    int
		CanCreatorClose bool
	}{}

	err = pgxscan.Select(context.Background(), bot.DB.Pool, &categories, "select * from ticket_categories where server_id = $1", ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(categories) == 0 {
		_, err = ctx.Reply("This server has no ticket categories.")
		return
	}

	fields := []discord.EmbedField{}

	for _, c := range categories {
		cat, err := ctx.State.Channel(c.CategoryID)
		if err != nil {
			continue
		}

		fields = append(fields, discord.EmbedField{
			Name:  cat.Name,
			Value: fmt.Sprintf("Name in commands: ``%v``\nLog channel: %v\nCount: %v\nPer-user limit: %v\nCan creator close: %v", c.Name, c.LogChannel.Mention(), c.Count, c.PerUserLimit, c.CanCreatorClose),
		})

		mention, desc := c.Mention, c.Description
		if mention == "" {
			mention = "<none>"
		}
		if desc == "" {
			desc = "<none>"
		}

		fields = append(fields, []discord.EmbedField{
			{
				Name:  "Mention",
				Value: mention,
			},
			{
				Name:  "Description",
				Value: desc,
			},
		}...)
	}

	_, err = ctx.PagedEmbed(bcr.FieldPaginator("Ticket categories", "", bcr.ColourBlurple, fields, 3), false)
	return
}
