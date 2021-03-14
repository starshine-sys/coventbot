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
		_, err = ctx.Send(":x: I'm not in the given server, so it has no tags.", nil)
		return
	}

	_, err = ctx.NewDM(ctx.Author.ID).Content(fmt.Sprintf("Check out `%vhelp` for a list of built-in commands.\n**%v** has the following tags:", ctx.Router.Prefixes[0], g.Name)).Send()
	if err != nil {
		return err
	}

	var msgs []string
	var b strings.Builder
	for i, n := range tagNames {
		if b.Len() >= 1900 {
			msgs = append(msgs, b.String())
			b.Reset()
		}
		b.WriteString(n)
		if i != len(tagNames)-1 {
			b.WriteString(", ")
		}
	}
	msgs = append(msgs, b.String())

	for _, m := range msgs {
		_, err = ctx.NewDM(ctx.Author.ID).Content(m).Send()
		if err != nil {
			return err
		}
	}

	return
}
