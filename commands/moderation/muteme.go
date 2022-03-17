package moderation

import (
	"fmt"
	"strings"
	"time"

	"codeberg.org/eviedelta/detctime/durationparser"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"gitlab.com/1f320/x/duration"
)

func (bot *Bot) muteme(ctx *bcr.Context) (err error) {
	dur, err := durationparser.Parse(ctx.RawArgs)
	if err != nil {
		dur, err = durationparser.Parse(ctx.Args[0])
		if err != nil {
			_, err = ctx.Replyc(bcr.ColourRed, "Couldn't parse your input as a valid duration.")
			return
		}
	}

	roles, err := bot.muteRoles(ctx.Guild.ID)
	if err != nil || !roles.MuteRole.IsValid() {
		_, err = ctx.Replyc(bcr.ColourRed, "This server doesn't have a mute role.")
		return
	}

	if dur < time.Minute {
		_, err = ctx.Replyc(bcr.ColourRed, "Duration needs to be at least one minute.")
		return
	} else if dur > 30*24*time.Hour {
		_, err = ctx.Replyc(bcr.ColourRed, "Can't mute yourself for more than 30 days.")
		return
	}

	msg, err := bot.mutemeMessage(ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	msg = strings.NewReplacer(
		"{mention}", ctx.Author.Mention(),
		"{tag}", ctx.Author.Tag(),
		"{duration}", duration.Format(dur),
		"{action}", "mute",
	).Replace(msg)

	err = ctx.State.AddRole(ctx.Guild.ID, ctx.Author.ID, roles.MuteRole, api.AddRoleData{
		AuditLogReason: api.AuditLogReason(
			fmt.Sprintf("Self-mute by %v for %v",
				ctx.Author.Tag(),
				duration.Format(dur),
			)),
	})
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "I couldn't mute you. Contact the mods for help.")
		return
	}

	reason := fmt.Sprintf("Self-mute from %v ago expired", duration.Format(dur))

	_, err = bot.Scheduler.Add(time.Now().Add(dur), &changeRoles{
		GuildID:        ctx.Message.GuildID,
		UserID:         ctx.Author.ID,
		RemoveRoles:    []discord.RoleID{roles.MuteRole},
		AuditLogReason: reason,
	})
	if err != nil {
		return bot.Report(ctx, err)
	}

	if msg == "" {
		return
	}

	ctx.NewDM(ctx.Author.ID).Content("**" + ctx.Guild.Name + "**: " + msg).Send()
	_, err = ctx.Send(msg)
	return err
}

func (bot *Bot) pauseme(ctx *bcr.Context) (err error) {
	dur, err := durationparser.Parse(ctx.RawArgs)
	if err != nil {
		dur, err = durationparser.Parse(ctx.Args[0])
		if err != nil {
			_, err = ctx.Replyc(bcr.ColourRed, "Couldn't parse your input as a valid duration.")
			return
		}
	}

	roles, err := bot.muteRoles(ctx.Guild.ID)
	if err != nil || !roles.PauseRole.IsValid() {
		_, err = ctx.Replyc(bcr.ColourRed, "This server doesn't have a pause role.")
		return
	}

	if dur < time.Minute {
		_, err = ctx.Replyc(bcr.ColourRed, "Duration needs to be at least one minute.")
		return
	} else if dur > 30*24*time.Hour {
		_, err = ctx.Replyc(bcr.ColourRed, "Can't pause yourself for more than 30 days.")
		return
	}

	msg, err := bot.mutemeMessage(ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	msg = strings.NewReplacer(
		"{mention}", ctx.Author.Mention(),
		"{tag}", ctx.Author.Tag(),
		"{duration}", duration.Format(dur),
		"{action}", "pause",
	).Replace(msg)

	err = ctx.State.AddRole(ctx.Guild.ID, ctx.Author.ID, roles.PauseRole, api.AddRoleData{
		AuditLogReason: api.AuditLogReason(
			fmt.Sprintf("Self-pause by %v for %v",
				ctx.Author.Tag(),
				duration.Format(dur),
			)),
	})
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "I couldn't pause you. Contact the mods for help.")
		return
	}

	reason := fmt.Sprintf("Self-pause from %v ago expired", duration.Format(dur))

	_, err = bot.Scheduler.Add(time.Now().Add(dur), &changeRoles{
		GuildID:        ctx.Message.GuildID,
		UserID:         ctx.Author.ID,
		RemoveRoles:    []discord.RoleID{roles.PauseRole},
		AuditLogReason: reason,
	})
	if err != nil {
		return bot.Report(ctx, err)
	}

	if msg == "" {
		return
	}

	ctx.NewDM(ctx.Author.ID).Content("**" + ctx.Guild.Name + "**: " + msg).Send()
	_, err = ctx.Send(msg)
	return err
}

func (bot *Bot) cmdMutemeMessage(ctx *bcr.Context) (err error) {
	if ctx.RawArgs == "" {
		current, err := bot.mutemeMessage(ctx.Guild.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		_, err = ctx.Reply("```\n" + current + "\n```")
		return err
	}

	if ctx.RawArgs == "clear" || ctx.RawArgs == "-clear" || ctx.RawArgs == "--clear" {
		err = bot.setMutemeMessage(ctx.Guild.ID, "")
		if err != nil {
			return bot.Report(ctx, err)
		}

		_, err = ctx.Reply("`%vmuteme` message cleared!", ctx.Prefix)
		return err
	}

	err = bot.setMutemeMessage(ctx.Guild.ID, ctx.RawArgs)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("`%vmuteme` message updated!", ctx.Prefix)
	return err
}
