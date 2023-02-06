// SPDX-License-Identifier: AGPL-3.0-only
package admin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

type AllowedGuild struct {
	ID       discord.GuildID
	Reason   string
	AddedBy  discord.UserID
	AddedFor discord.UserID
	AddedAt  time.Time
}

func (bot *Bot) getAllowedGuilds() (guilds []AllowedGuild, err error) {
	return guilds, pgxscan.Select(context.Background(), bot.DB.Pool, &guilds, "select * from allowed_guilds order by added_at")
}

func (bot *Bot) listAllowedGuilds(ctx *bcr.Context) (err error) {
	guilds, err := bot.getAllowedGuilds()
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(guilds) == 0 {
		return ctx.SendfX("There are no allowed guilds, %v is probably not private.", bot.Router.Bot.Username)
	}

	fields := make([]discord.EmbedField, 0, len(guilds))
	for _, g := range guilds {
		fields = append(fields, discord.EmbedField{
			Name: g.ID.String(),
			Value: fmt.Sprintf(`**Added by:** %v
**Added for:** %v
**Added at:** <t:%v>
**Reason:** %v`, g.AddedBy.Mention(), g.AddedFor.Mention(), g.AddedAt.Unix(), g.Reason),
		})
	}

	_, err = bot.PagedEmbed(ctx,
		bcr.FieldPaginator("Allowed guilds", "", bcr.ColourBlurple, fields, 5), 10*time.Minute)
	return err
}

func (bot *Bot) addAllowedGuild(ctx *bcr.Context) (err error) {
	sf, err := discord.ParseSnowflake(ctx.Args[0])
	if err != nil {
		return ctx.SendfX("Couldn't parse ``%v`` as a valid snowflake.", bcr.EscapeBackticks(ctx.Args[0]))
	}

	if bot.isGuildAllowed(discord.GuildID(sf)) {
		return ctx.SendfX("Guild %v is already allowed!", sf)
	}

	u, err := ctx.ParseUser(ctx.Args[1])
	if err != nil {
		return ctx.SendfX("Couldn't find a user with the ID ``%v``.", bcr.EscapeBackticks(ctx.Args[1]))
	}

	reason := "No reason given"
	if len(ctx.Args) > 2 {
		reason = strings.TrimSpace(
			strings.TrimPrefix(
				strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0])), ctx.Args[1],
			),
		)

		if reason == ctx.RawArgs || strings.Contains(reason, ctx.Args[1]) {
			reason = strings.Join(ctx.Args[2:], " ")
		}
	}

	yes, _ := ctx.ConfirmButton(ctx.Author.ID, bcr.ConfirmData{
		Embeds: []discord.Embed{{
			Description: fmt.Sprintf("Are you sure you want to allow %v to join the guild %v?\n**For:** %v (%v)\n**Reason:** %v", bot.Router.Bot.Username, sf, u.Tag(), u.ID, reason),
			Color:       bcr.ColourBlurple,
		}},
		YesPrompt: "Allow this guild",
		NoPrompt:  "Do not allow this guild",
	})
	if !yes {
		return ctx.SendX("Cancelled.")
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "insert into allowed_guilds (id, reason, added_by, added_for) values ($1, $2, $3, $4)", sf, reason, ctx.Author.ID, u.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	return ctx.SendfX("Success! Guild %v can now be joined by %v.", sf, bot.Router.Bot.Username)
}
