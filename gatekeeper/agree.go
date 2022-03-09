package gatekeeper

import (
	"fmt"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) agree(ctx *bcr.Context) (err error) {
	settings, err := bot.serverSettings(ctx.Message.GuildID)
	if err != nil {
		bot.Sugar.Errorf("Error getting server settings: %v", err)
		return bot.Report(ctx, err)
	}

	if !settings.MemberRole.IsValid() {
		_, err = ctx.Send("This server does not use a gateway.")
		return
	}

	if !bot.isPending(ctx.Message.GuildID, ctx.Author.ID) {
		_, err = ctx.Send("You are not a pending user.")
		return
	}

	p, err := bot.pendingUser(ctx.Message.GuildID, ctx.Author.ID)
	if err != nil {
		bot.Sugar.Errorf("Error getting user: %v", err)
		return bot.Report(ctx, err)
	}

	url := fmt.Sprintf("%v/gatekeeper/%v", bot.Config.HTTPBaseURL, p.Key)

	_, err = ctx.NewDM(ctx.Author.ID).Content(
		fmt.Sprintf("Please solve the captcha at the following link to verify that you're a human: <%v>", url),
	).Send()
	if err != nil {
		_, err = ctx.Sendf("%v, could not DM you a captcha link. Please make sure you have DMs enabled, if not, please message this server's moderators.", ctx.Author.Mention())
		return
	}

	_, err = ctx.Sendf("%v, check your DMs!", ctx.Author.Mention())
	return
}
