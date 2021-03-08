package tags

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) list(ctx *bcr.Context) (err error) {
	guildID := discord.Snowflake(ctx.Message.GuildID)

	if !ctx.Message.GuildID.IsValid() {
		if len(ctx.Args) < 1 {
			_, err = ctx.Send(":x: No server ID provided!", nil)
			return
		}
		guildID, err = discord.ParseSnowflake(ctx.Args[0])
		if err != nil {
			_, err = ctx.Send(":x: No valid server ID provided!", nil)
			return
		}
	}

	tags, err := bot.DB.Tags(discord.GuildID(guildID))
	if err != nil {
		_, err = ctx.Send(":x: Internal error occurred.", nil)
		return
	}

	var tagNames []string
	for _, t := range tags {
		tagNames = append(tagNames, t.Name)
	}
	if len(tagNames) == 0 {
		tagNames = []string{"This server has no tags."}
	}

	// we need to grab the server name
	g, err := ctx.Session.Guild(discord.GuildID(guildID))
	if err != nil {
		_, err = ctx.Send(":x: Internal error occurred.", nil)
		return
	}

	_, err = ctx.NewDM(ctx.Author.ID).Content(fmt.Sprintf("Check out `%vhelp` for a list of built-in commands.\n**%v** has the following tags:", ctx.Router.Prefixes[0], g.Name)).Send()
	if err != nil {
		return err
	}

	_, err = ctx.NewDM(ctx.Author.ID).Content(strings.Join(tagNames, ", ")).Send()
	return
}
