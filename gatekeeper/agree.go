// SPDX-License-Identifier: AGPL-3.0-only
package gatekeeper

import (
	"fmt"
	"strings"

	"github.com/starshine-sys/bcr"
)

const defaultAgreeResp = "{mention}, check your DMs!"

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

	for _, r := range ctx.Member.RoleIDs {
		if r == settings.MemberRole {
			return ctx.SendX("You have already passed the gateway.")
		}
	}

	p, err := bot.setPending(ctx.Message.GuildID, ctx.Author.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if !p.Pending {
		return ctx.SendX("You have already passed the gateway.")
	}

	url := fmt.Sprintf("%v/gatekeeper/%v", bot.Config.HTTPBaseURL, p.Key)

	_, err = ctx.NewDM(ctx.Author.ID).Content(
		fmt.Sprintf("Please solve the captcha at the following link to verify that you're a human: <%v>", url),
	).Send()
	if err != nil {
		_, err = ctx.Sendf("%v, could not DM you a captcha link. Please make sure you have DMs enabled, if not, please message this server's moderators.", ctx.Author.Mention())
		return
	}

	if settings.GatekeeperLog.IsValid() {
		_, err = ctx.State.SendMessage(settings.GatekeeperLog, fmt.Sprintf("User %v / %v has passed the agree step of the gateway, now doing verification.", ctx.Author.Tag(), ctx.Author.Mention()))
		if err != nil {
			bot.Sugar.Errorf("sending gatekeeper log message: %v", err)
		}
	}

	resp, err := bot.DB.GuildStringGet(ctx.Message.GuildID, "gateway:agree_response")
	if err != nil || resp == "" {
		return ctx.SendX(
			strings.ReplaceAll(defaultAgreeResp, "{mention}", ctx.Author.Mention()),
		)
	}

	return ctx.SendX(
		strings.ReplaceAll(resp, "{mention}", ctx.Author.Mention()),
	)
}
