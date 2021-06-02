package moderation

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/json/option"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) ban(ctx *bcr.Context) (err error) {
	var (
		target   *discord.User
		isMember bool
	)

	member, err := ctx.ParseMember(ctx.Args[0])
	if err == nil {
		if !bot.aboveUser(ctx, member) {
			_, err = ctx.Send("You're not high enough in the role hierarchy to do that.", nil)
			return
		}
		isMember = true
		target = &member.User
	} else {
		target, err = ctx.ParseUser(ctx.Args[0])
		if err != nil {
			_, err = ctx.Send("Couldn't find a user with that name.", nil)
			return
		}
	}

	// check bot perms
	if p, _ := ctx.State.Permissions(ctx.Channel.ID, ctx.Bot.ID); !p.Has(discord.PermissionBanMembers) {
		_, err = ctx.Send("I do not have the **Ban Members** permission.", nil)
		return
	}

	reason := "N/A"
	if len(ctx.Args) > 1 {
		reason = strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))
	}

	if isMember {
		_, err = ctx.NewDM(target.ID).Content(fmt.Sprintf("You were banned from %v.\nReason: %v", ctx.Guild.Name, reason)).Send()
		if err != nil {
			ctx.Send("I was unable to DM the user about their ban.", nil)
		}
	}

	err = ctx.State.Ban(ctx.Message.GuildID, target.ID, api.BanData{
		DeleteDays: option.NewUint(0),
		Reason: option.NewString(
			fmt.Sprintf("%v#%v: %v", ctx.Author.Username, ctx.Author.Discriminator, reason)),
	})
	if err != nil {
		_, err = ctx.Send("I could not ban that user.", nil)
		return
	}

	err = bot.ModLog.Ban(ctx.Message.GuildID, target.ID, ctx.Author.ID, reason)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Sendf("Banned **%v#%v**", target.Username, target.Discriminator)
	return
}

func (bot *Bot) unban(ctx *bcr.Context) (err error) {
	u, err := ctx.ParseUser(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("I couldn't find that user.", nil)
		return
	}

	// check bot perms
	if p, _ := ctx.State.Permissions(ctx.Channel.ID, ctx.Bot.ID); !p.Has(discord.PermissionBanMembers) {
		_, err = ctx.Send("I do not have the **Ban Members** permission.", nil)
		return
	}

	reason := "N/A"
	if len(ctx.Args) > 1 {
		reason = strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))
	}

	bans, err := ctx.State.Bans(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	var isBanned bool
	for _, b := range bans {
		if b.User.ID == u.ID {
			isBanned = true
			break
		}
	}

	if !isBanned {
		_, err = ctx.Send("That user is not banned.", nil)
		return
	}

	err = ctx.State.Unban(ctx.Message.GuildID, u.ID)
	if err != nil {
		_, err = ctx.Sendf("I was unable to unban %v#%v.", u.Username, u.Discriminator)
		return
	}

	err = bot.ModLog.Unban(ctx.Message.GuildID, u.ID, ctx.Author.ID, reason)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Sendf("Unbanned **%v#%v**", u.Username, u.Discriminator)
	return
}
