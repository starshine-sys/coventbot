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
			_, err = ctx.Send("No server ID provided!", nil)
			return
		}
		guildID, err = discord.ParseSnowflake(ctx.Args[0])
		if err != nil {
			_, err = ctx.Send("No valid server ID provided!", nil)
			return
		}
	}

	tags, err := bot.DB.Tags(discord.GuildID(guildID))
	if err != nil {
		return bot.Report(ctx, err)
	}

	var tagNames []string
	for _, t := range tags {
		tagNames = append(tagNames, t.Name)
	}
	if len(tagNames) == 0 {
		tagNames = []string{"This server has no tags."}
	}

	// we need to grab the server name
	g, err := ctx.State.Guild(discord.GuildID(guildID))
	if err != nil {
		_, err = ctx.Send("I'm not in the given server, so it has no tags.", nil)
		return
	}

	prefixes, err := bot.DB.Prefixes(discord.GuildID(guildID))
	if err != nil || len(prefixes) == 0 {
		prefixes = ctx.Router.Prefixes
	}

	s := fmt.Sprintf(`**%v**
Check out `+"`%vhelp`"+` for a list of built-in commands.
This server's prefixes are: %v

This server has the following tags:`, g.Name, prefixes[0], strings.Join(prefixes, ", "))

	_, err = ctx.NewDM(ctx.Author.ID).Content(s).Send()
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
