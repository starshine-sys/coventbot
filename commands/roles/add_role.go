// SPDX-License-Identifier: AGPL-3.0-only
package roles

import (
	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/jackc/pgx/v4"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) addRole(ctx *bcr.Context) (err error) {
	if ctx.Member == nil {
		return
	}

	r, err := ctx.ParseRole(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "No role called ``%v`` found.", ctx.RawArgs)
	}

	for _, id := range ctx.Member.RoleIDs {
		if r.ID == id {
			_, err = ctx.Replyc(bcr.ColourRed, "You already have that role.")
			return
		}
	}

	cat, err := bot.categoryRole(ctx.Guild.ID, r.ID)
	if err != nil {
		if errors.Cause(err) != pgx.ErrNoRows {
			return bot.Report(ctx, err)
		}

		_, err = ctx.Replyc(bcr.ColourRed, "That role isn't self-assignable.")
		return
	}

	if cat.RequireRole.IsValid() {
		perm := false
		for _, r := range ctx.Member.RoleIDs {
			if r == cat.RequireRole {
				perm = true
				break
			}
		}

		if !perm {
			_, err = ctx.Replyc(bcr.ColourRed, "You can't assign that to yourself.")
			return
		}
	}

	err = ctx.State.AddRole(ctx.Guild.ID, ctx.Author.ID, r.ID, api.AddRoleData{
		AuditLogReason: "Self-assigned role",
	})
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "I couldn't assign that role to you.")
		return
	}

	_, err = ctx.Reply("Added %v to %v.", r.Mention(), ctx.Author.Mention())
	return
}

func (bot *Bot) removeRole(ctx *bcr.Context) (err error) {
	if ctx.Member == nil {
		return
	}

	r, err := ctx.ParseRole(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "No role called ``%v`` found.", ctx.RawArgs)
	}

	hasRole := false
	for _, id := range ctx.Member.RoleIDs {
		if r.ID == id {
			hasRole = true
			break
		}
	}
	if !hasRole {
		_, err = ctx.Replyc(bcr.ColourRed, "You don't have that role.")
		return
	}

	cat, err := bot.categoryRole(ctx.Guild.ID, r.ID)
	if err != nil {
		if errors.Cause(err) != pgx.ErrNoRows {
			return bot.Report(ctx, err)
		}

		_, err = ctx.Replyc(bcr.ColourRed, "That role isn't self-assignable.")
		return
	}

	if cat.RequireRole.IsValid() {
		perm := false
		for _, r := range ctx.Member.RoleIDs {
			if r == cat.RequireRole {
				perm = true
				break
			}
		}

		if !perm {
			_, err = ctx.Replyc(bcr.ColourRed, "You can't assign that to yourself.")
			return
		}
	}

	err = ctx.State.RemoveRole(ctx.Guild.ID, ctx.Author.ID, r.ID, "Self-assigned role")
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "I couldn't remove that role from you.")
		return
	}

	_, err = ctx.Reply("Removed %v from %v.", r.Mention(), ctx.Author.Mention())
	return
}
