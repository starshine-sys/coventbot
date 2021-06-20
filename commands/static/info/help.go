package info

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) help(ctx *bcr.Context) (err error) {
	if ctx.RawArgs != "" {
		return ctx.Help(ctx.Args)
	}

	desc := fmt.Sprintf(`For help with a command, use "%vhelp <command>"
For a list of all commands, use "%vhelp commands"

Note that this bot is currently **very** undocumented; use at your own risk.`, ctx.Router.Prefixes[0], ctx.Router.Prefixes[0])

	if bot.Config.Branding.SupportServer != "" {
		desc += "\n\nFor help and feedback, join the support server: " + bot.Config.Branding.SupportServer
	}

	fields := []discord.EmbedField{
		{
			Name:  "Information",
			Value: "`userinfo`: show member/user info\n`avatar`: show a user's avatar\n`roleinfo`: show info about a role\n`getinvite`: get detailed info from an invite link\n`serverinfo`: show info about the current server\n`idtime`: get the time from an Discord snowflake",
		},
		{
			Name:  "Bot info",
			Value: "`about`: show info about the bot\n`invite`: get an invite link for the bot\n`ping`: get the bot's latency",
		},
		{
			Name:  "Utility",
			Value: "`addemoji`: add a custom emoji\n`embedsource`: get the JSON source from an embed\n`enlarge`: enlarge a custom emote\n`meow`: show a random meowmoji\n`poll` & `quickpoll`: run polls\n`roll`: roll for initiative\n`echo`: send messages as the bot\n`members`: show a filtered list of members",
		},
		{
			Name:  "Reminders and todos",
			Value: "`remindme`: set a reminder for yourself\n`reminders`: show your current reminders\n`delreminder`: delete one of your reminders\n`todo`: set a todo and send it to your todo channel\n`complete`: complete a todo",
		},
		{
			Name:  "Levels",
			Value: "`lvl`: show your or another user's level\n`lvl config`: configure levels\n`leaderboard`: show this server's leaderboard",
		},
		{
			Name:  "Moderation",
			Value: "`lock`: lock or unlock a channel from @everyone\n`makeinvite`: make a permanent invite for a channel\n`massban`: ban multiple people at once\n`watchlist`: get warned when people join\n`nicknames`\n`usernames`: track username and nickname changes",
		},
		{
			Name:  "Moderation (cont.)",
			Value: "`slowmode`: set slowmode (both Discord and custom)\n`channelban`\n`unchannelban`: block/unblock users from channels\n`warn`: warn members",
		},
	}

	_, err = ctx.PagedEmbed(
		bcr.FieldPaginator("Help", desc, bcr.ColourBlurple, fields, 2), false,
	)
	return
}
