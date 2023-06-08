// SPDX-License-Identifier: AGPL-3.0-only
package approval

import (
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) approve(ctx *bcr.Context) (err error) {
	m, err := ctx.ParseMember(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Sendf("That member could not be found.")
		return
	}

	s, err := bot.serverSettings(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(s.ApproveAddRoles) == 0 {
		_, err = ctx.Sendf("Approval is not set up on this server.")
		return
	}

	for _, r := range s.ApproveAddRoles {
		err = ctx.State.AddRole(ctx.Message.GuildID, m.User.ID, discord.RoleID(r), api.AddRoleData{
			AuditLogReason: "Gatekeeper: approve member",
		})
		if err != nil {
			return bot.Report(ctx, err)
		}
	}
	for _, r := range s.ApproveRemoveRoles {
		err = ctx.State.RemoveRole(ctx.Message.GuildID, m.User.ID, discord.RoleID(r), "Gatekeeper: approve member")
		if err != nil {
			return bot.Report(ctx, err)
		}
	}

	if s.ApproveWelcomeChannel.IsValid() && s.ApproveWelcomeMessage != "" {
		_, err = ctx.NewMessage(s.ApproveWelcomeChannel).Content(
			strings.NewReplacer("{mention}", m.Mention()).Replace(s.ApproveWelcomeMessage),
		).Send()
		if err != nil {
			return bot.Report(ctx, err)
		}
	}

	_, err = ctx.Sendf("Approved **%v**.", m.User.Tag())
	return
}
