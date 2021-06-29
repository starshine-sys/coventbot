package roles

import (
	"fmt"
	"sort"
	"strings"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/jackc/pgx/v4"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) categories(ctx *bcr.Context) (err error) {
	if len(ctx.Args) == 0 {
		cats, err := bot.guildCategories(ctx.Guild.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		if len(cats) == 0 {
			_, err = ctx.Reply("There are no role categories.")
			return err
		}

		entries := []string{}

		for _, c := range cats {
			s := c.Name
			s += fmt.Sprintf(": %v role", len(c.Roles))
			if len(c.Roles) != 1 {
				s += "s"
			}
			s += "\n"

			entries = append(entries, s)
		}

		embeds := bcr.StringPaginator("Categories", bcr.ColourBlurple, entries, 10)
		for i := range embeds {
			embeds[i].Fields = []discord.EmbedField{{
				Name:  "_ _",
				Value: fmt.Sprintf("Use `%vroles <category>` to show all roles in a category.\nUse `%vrole` to give yourself a role, and `%vderole` to remove it.", ctx.Prefix, ctx.Prefix, ctx.Prefix),
			}}
		}

		_, err = ctx.PagedEmbed(embeds, false)
		return err
	}

	cat, err := bot.categoryName(ctx.Guild.ID, ctx.RawArgs)
	if err != nil {
		if errors.Cause(err) != pgx.ErrNoRows {
			return bot.Report(ctx, err)
		}

		_, err = ctx.Replyc(bcr.ColourRed, "Couldn't find a category with that name. Try ``%vroles`` for a list.", ctx.Prefix)
		return
	}

	if cat.RequireRole.IsValid() {
		if ctx.Member == nil {
			_, err = ctx.Replyc(bcr.ColourRed, "You don't have permission to use this category.")
			return
		}

		perm := false
		for _, r := range ctx.Member.RoleIDs {
			if r == cat.RequireRole {
				perm = true
				break
			}
		}

		if !perm {
			_, err = ctx.Replyc(bcr.ColourRed, "You don't have permission to use this category.")
			return
		}
	}

	roles, err := bot.roles(ctx.Guild.ID, cat.Roles)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(roles) == 0 {
		_, err = ctx.Reply("There are no roles in %v.", cat.Name)
		return
	}

	sort.Slice(roles, sortByName(roles))

	s := []string{}
	for _, r := range roles {
		s = append(s, r.Mention())
	}

	desc := ""
	if cat.Description != "" {
		desc = cat.Description + "\n\n"
	}

	desc += strings.Join(s, ", ")

	e := discord.Embed{
		Title:       cat.Name,
		Description: desc,
		Color:       cat.Colour,
		Fields: []discord.EmbedField{{
			Name:  "_ _",
			Value: fmt.Sprintf("Use `%vrole` to give yourself a role, and `%vderole` to remove it.", ctx.Prefix, ctx.Prefix),
		}},
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("Category ID: %v", cat.ID),
		},
	}

	if e.Color == 0 {
		e.Color = bcr.ColourBlurple
	}

	_, err = ctx.Send("", &e)
	return
}
