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

	p, err := bot.setPending(ctx.Message.GuildID, ctx.Author.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if !p.Pending {
		return ctx.SendX("You have already passed the gatekeeper.")
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
