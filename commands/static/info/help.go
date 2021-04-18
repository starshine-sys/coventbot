package info

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/etc"
)

func (bot *Bot) help(ctx *bcr.Context) (err error) {
	if ctx.RawArgs != "" {
		return ctx.Help(ctx.Args)
	}

	e := discord.Embed{
		Title: "Help",
		Description: fmt.Sprintf(`For help with a command, use "%vhelp <command>"
For a list of all commands, use "%vhelp commands"

Note that this bot is currently **very** undocumented; use at your own risk.
Either way, here are some commonly used commands:`, ctx.Router.Prefixes[0], ctx.Router.Prefixes[0]),

		Fields: []discord.EmbedField{
			{
				Name:   "Information",
				Value:  "`info`: show user info\n`avatar`: show a user's avatar\n`roleinfo`: show role info\n`getinvite`: show invite info\n`serverinfo`: show server info\n`idtime`: get time from an ID",
				Inline: true,
			},
			{
				Name:   "Bot info",
				Value:  "`about`: info about the bot\n`invite`: bot invite link\n`ping`: bot latency",
				Inline: true,
			},
			{
				Name:  "​",
				Value: "​",
			},
			{
				Name:   "Utility",
				Value:  "`addemoji`: add a custom emote\n`embedsource`: get JSON of embeds\n`enlarge`: enlarge a custom emote\n`meow`: show a meowmoji\n`poll` & `quickpoll`: run polls\n`roll`: roll dice\n`echo`: send messages as the bot\n`members`: show a filtered list of members",
				Inline: true,
			},
			{
				Name:   "Moderation",
				Value:  "`lock`: lock a channel from @everyone\n`makeinvite`: make an invite\n`massban`: ban multiple people at once\n`watchlist`: get warned when people join\n`nicknames`\n`usernames`: track name changes",
				Inline: true,
			},
		},

		Color: etc.ColourBlurple,
	}

	if bot.Config.Branding.SupportServer != "" {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Support server",
			Value: "For help and feedback, join the support server: " + bot.Config.Branding.SupportServer,
		})
	}

	_, err = ctx.Send("", &e)
	return
}
