package starboard

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) settings(ctx *bcr.Context) (err error) {
	var b strings.Builder

	settings, err := bot.DB.Starboard(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	if !settings.StarboardChannel.IsValid() {
		b.WriteString(fmt.Sprintf("**There is no starboard channel set.** Set one with `%v%v channel`\n", ctx.Prefix, ctx.Command))
	} else {
		b.WriteString(fmt.Sprintf("The starboard channel is <#%v>. Use `%v%v channel -clear` to clear it\n", settings.StarboardChannel, ctx.Prefix, ctx.Command))
	}

	b.WriteString(fmt.Sprintf("The current starboard limit is %v stars\nThe current starboard emoji is %v", settings.StarboardLimit, settings.StarboardEmoji))

	_, err = ctx.Send("", &discord.Embed{
		Title:       "Starboard settings",
		Description: b.String(),
		Color:       ctx.Router.EmbedColor,
	})
	return err
}

func (bot *Bot) setChannel(ctx *bcr.Context) (err error) {
	var id discord.ChannelID

	if ctx.RawArgs == "-clear" {
		id = 0
	} else {
		ch, err := ctx.ParseChannel(ctx.RawArgs)
		if err != nil {
			_, err = ctx.Send("Channel not found.", nil)
			return err
		}

		if ch.GuildID != ctx.Message.GuildID {
			_, err = ctx.Send("The given channel isn't in this server.", nil)
			return err
		}

		id = ch.ID
	}

	settings, err := bot.DB.Starboard(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	if settings.StarboardChannel == id {
		_, err = ctx.Send("The given channel is already the starboard channel.", nil)
		return err
	}

	settings.StarboardChannel = id
	err = bot.DB.SetStarboard(ctx.Message.GuildID, settings)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	if id == 0 {
		_, err = ctx.Send("Starboard channel reset.", nil)
		return
	}
	_, err = ctx.Sendf("Starboard channel changed to %v.", id.Mention())
	return
}

func (bot *Bot) setEmoji(ctx *bcr.Context) (err error) {
	settings, err := bot.DB.Starboard(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	if settings.StarboardEmoji == ctx.RawArgs {
		_, err = ctx.Send("The given emoji is already the starboard emoji.", nil)
		return err
	}

	settings.StarboardEmoji = ctx.RawArgs
	err = bot.DB.SetStarboard(ctx.Message.GuildID, settings)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	_, err = ctx.NewMessage().Content(fmt.Sprintf("Starboard emoji changed to %v.", ctx.RawArgs)).BlockMentions().Send()
	return
}

func (bot *Bot) setLimit(ctx *bcr.Context) (err error) {
	i, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("Could not parse your input as a number.", nil)
		return
	}

	settings, err := bot.DB.Starboard(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	if settings.StarboardLimit == i {
		_, err = ctx.Send("The given limit is already the starboard limit.", nil)
		return err
	}

	settings.StarboardLimit = i
	err = bot.DB.SetStarboard(ctx.Message.GuildID, settings)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	_, err = ctx.Sendf("Starboard limit changed to %v stars.", i)
	return
}
