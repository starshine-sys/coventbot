package levels

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) colour(ctx *bcr.Context) (err error) {
	sc, err := bot.getGuildConfig(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}
	if !sc.LevelsEnabled {
		_, err = ctx.Send("Levels are disabled on this server.")
		return
	}

	if len(ctx.Args) == 0 {
		uc, err := bot.getUser(ctx.Message.GuildID, ctx.Author.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		if uc.Colour == 0 {
			_, err = ctx.Sendf("You don't currently have a colour set. Set one with `%vlvl colour <hex code>`.", ctx.Prefix)
			return err
		}
		url := fmt.Sprintf("https://fakeimg.pl/256x256/%06X/?text=%%20", uc.Colour)
		_, err = ctx.Send("", discord.Embed{
			Thumbnail: &discord.EmbedThumbnail{
				URL: url,
			},

			Description: fmt.Sprintf("Your level colour is currently set to **%v**. To clear it, type `%vlvl colour clear`.", uc.Colour.String(), ctx.Prefix),

			Color: uc.Colour,
		})
		return err
	}

	if ctx.RawArgs == "clear" {
		_, err = bot.DB.Pool.Exec(context.Background(), "update levels set colour = $1 where server_id = $2 and user_id = $3", 0, ctx.Message.GuildID, ctx.Author.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		_, err = ctx.Reply("Level colour cleared!")
		return
	}

	ctx.RawArgs = strings.TrimPrefix(ctx.RawArgs, "#")

	clr, err := strconv.ParseUint(ctx.RawArgs, 16, 0)
	if err != nil {
		_, err = ctx.Send("You didn't give a valid hex code colour.")
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update levels set colour = $1 where server_id = $2 and user_id = $3", clr, ctx.Message.GuildID, ctx.Author.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	url := fmt.Sprintf("https://fakeimg.pl/256x256/%06X/?text=%%20", clr)
	_, err = ctx.Send("", discord.Embed{
		Thumbnail: &discord.EmbedThumbnail{
			URL: url,
		},

		Color:       discord.Color(clr),
		Description: fmt.Sprintf("Level colour changed to **#%06X**.", clr),
	})
	return
}
