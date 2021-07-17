package moderation

import "github.com/starshine-sys/bcr"

func (bot *Bot) muteRole(ctx *bcr.Context) (err error) {
	if ctx.RawArgs == "" {
		roles, err := bot.muteRoles(ctx.Guild.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		if !roles.MuteRole.IsValid() {
			_, err = ctx.Reply("This server has no mute role set.")
		} else {
			_, err = ctx.Reply("This server's mute role is %v.", roles.MuteRole.Mention())
		}
		return err
	}

	role, err := ctx.ParseRole(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "Role not found.")
		return
	}

	botMember, err := bot.Member(ctx.Guild.ID, ctx.Bot.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	var highestUserRolePos, highestBotRolePos int
	for _, r := range ctx.Guild.Roles {
		for _, rID := range botMember.RoleIDs {
			if rID == r.ID {
				if r.Position > highestBotRolePos {
					highestBotRolePos = r.Position
				}
			}
		}

		for _, rID := range ctx.Member.RoleIDs {
			if rID == r.ID {
				if r.Position > highestUserRolePos {
					highestUserRolePos = r.Position
				}
			}
		}
	}

	if role.Position >= highestUserRolePos {
		_, err = ctx.Replyc(bcr.ColourRed, "You can't manage this role.")
		return
	}
	if role.Position >= highestBotRolePos {
		_, err = ctx.Replyc(bcr.ColourRed, "I can't manage this role.")
		return
	}

	err = bot.setMuteRole(ctx.Guild.ID, role.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Mute role set to %v!", role.Mention())
	return
}

func (bot *Bot) pauseRole(ctx *bcr.Context) (err error) {
	if ctx.RawArgs == "" {
		roles, err := bot.muteRoles(ctx.Guild.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		if !roles.PauseRole.IsValid() {
			_, err = ctx.Reply("This server has no pause role set.")
		} else {
			_, err = ctx.Reply("This server's pause role is %v.", roles.PauseRole.Mention())
		}
		return err
	}

	role, err := ctx.ParseRole(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "Role not found.")
		return
	}

	botMember, err := bot.Member(ctx.Guild.ID, ctx.Bot.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	var highestUserRolePos, highestBotRolePos int
	for _, r := range ctx.Guild.Roles {
		for _, rID := range botMember.RoleIDs {
			if rID == r.ID {
				if r.Position > highestBotRolePos {
					highestBotRolePos = r.Position
				}
			}
		}

		for _, rID := range ctx.Member.RoleIDs {
			if rID == r.ID {
				if r.Position > highestUserRolePos {
					highestUserRolePos = r.Position
				}
			}
		}
	}

	if role.Position >= highestUserRolePos {
		_, err = ctx.Replyc(bcr.ColourRed, "You can't manage this role.")
		return
	}
	if role.Position >= highestBotRolePos {
		_, err = ctx.Replyc(bcr.ColourRed, "I can't manage this role.")
		return
	}

	err = bot.setPauseRole(ctx.Guild.ID, role.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Pause role set to %v!", role.Mention())
	return
}
