// SPDX-License-Identifier: AGPL-3.0-only
package moderation

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/commands/moderation/modlog"
)

func (bot *Bot) warn(ctx *bcr.Context) (err error) {
	go ctx.State.Typing(ctx.Channel.ID)

	u, err := ctx.ParseMember(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("User not found.")
		return
	}

	reason := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))

	if u.User.ID == ctx.Bot.ID {
		_, err = ctx.Send("😭 Why would you do that?")
		return
	}

	if !bot.aboveUser(ctx, ctx.Member, u) {
		_, err = ctx.Send("You're not high enough in the hierarchy to do that.")
		return
	}

	err = bot.ModLog.Log(ctx.State, modlog.ActionWarn, ctx.Message.GuildID, u.User.ID, ctx.Author.ID, reason)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.NewDM(u.User.ID).Content(fmt.Sprintf("You were warned in %v.\nReason: %v", ctx.Guild.Name, reason)).Send()
	if err != nil {
		_, err = ctx.Send("The warning was logged, but I was unable to notify the user of their warning.")
		return
	}

	var count int
	err = bot.DB.Pool.QueryRow(context.Background(), "select count(*) from mod_log where user_id = $1 and server_id = $2 and action_type = 'warn'", u.User.ID, ctx.Message.GuildID).Scan(&count)
	if err != nil {
		count = 1
	}

	_, err = ctx.NewMessage().Content(fmt.Sprintf("**%v#%v** has been warned, this is their %v warning.", u.User.Username, u.User.Discriminator, humanize.Ordinal(count))).Send()
	return
}

func (bot *Bot) aboveUser(ctx *bcr.Context, mod *discord.Member, member *discord.Member) (above bool) {
	if ctx.Guild == nil {
		return false
	}

	if ctx.Guild.OwnerID == mod.User.ID {
		return true
	}

	var modRoles, memberRoles bcr.Roles
	for _, r := range ctx.Guild.Roles {
		for _, id := range mod.RoleIDs {
			if r.ID == id {
				modRoles = append(modRoles, r)
				break
			}
		}
		for _, id := range member.RoleIDs {
			if r.ID == id {
				memberRoles = append(memberRoles, r)
				break
			}
		}
	}

	if len(modRoles) == 0 {
		return false
	}
	if len(memberRoles) == 0 {
		return true
	}

	sort.Sort(modRoles)
	sort.Sort(memberRoles)

	return modRoles[0].Position > memberRoles[0].Position
}
