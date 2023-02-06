// SPDX-License-Identifier: AGPL-3.0-only
package tickets

import (
	"context"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

// Category ...
type Category struct {
	Name string

	CategoryID discord.ChannelID

	ServerID discord.GuildID

	PerUserLimit    int
	Count           uint
	CanCreatorClose bool

	LogChannel discord.ChannelID

	Mention     string
	Description string
}

func (bot *Bot) new(ctx *bcr.Context) (err error) {
	user := ctx.Author

	cat := &Category{}
	err = pgxscan.Get(context.Background(), bot.DB.Pool, cat, "select * from ticket_categories where server_id = $1 and name ilike $2 order by category_id limit 1", ctx.Message.GuildID, ctx.Args[0])
	if err != nil {
		_, err = ctx.Sendf("Category not found.")
		return
	}

	if cat.PerUserLimit != -1 {
		var currentTickets int
		err = bot.DB.Pool.QueryRow(context.Background(), "select count(*) from tickets where owner_id = $1 and category_id = $2", ctx.Author.ID, cat.CategoryID).Scan(&currentTickets)
		if currentTickets >= cat.PerUserLimit {
			_, err = ctx.NewDM(ctx.Author.ID).Content("You have reached the ticket limit for this category.").Send()
			return
		}
	}

	if perms, _ := ctx.State.Permissions(cat.CategoryID, ctx.Author.ID); len(ctx.Args) > 1 && perms.Has(discord.PermissionManageMessages) {
		member, err := ctx.ParseMember(ctx.Args[1])
		if err != nil {
			_, err = ctx.Send("Member not found.")
			return err
		}
		user = member.User
	}

	ch, err := ctx.State.CreateChannel(ctx.Message.GuildID, api.CreateChannelData{
		Name:           fmt.Sprintf("%v-%04d", cat.Name, cat.Count+1),
		Type:           discord.GuildText,
		CategoryID:     cat.CategoryID,
		AuditLogReason: api.AuditLogReason("Creating ticket channel in category " + cat.Name),
	})
	if err != nil {
		bot.Sugar.Errorf("Error creating channel: %v", err)
		_, err = ctx.Send("There was an error creating the ticket channel. Are you sure I have the manage channels permission there, and are you sure the category isn't full?")
		return
	}

	err = ctx.State.EditChannelPermission(ch.ID, discord.Snowflake(user.ID), api.EditChannelPermissionData{
		Type:  discord.OverwriteMember,
		Allow: discord.PermissionViewChannel | discord.PermissionSendMessages | discord.PermissionReadMessageHistory | discord.PermissionAddReactions,
	})
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "update ticket_categories set count = $1 where category_id = $2", cat.Count+1, cat.CategoryID)
	if err != nil {
		bot.Sugar.Errorf("Error updating count: %v", err)
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "insert into tickets (channel_id, category_id, owner_id, users) values ($1, $2, $3, $4)", ch.ID, cat.CategoryID, user.ID, []uint64{uint64(user.ID)})
	if err != nil {
		return bot.Report(ctx, err)
	}

	mention := strings.NewReplacer(
		"{mention}", user.Mention(),
		"{channel}", ch.Mention(),
		"{here}", "@here",
		"{everyone}", "@everyone",
	).Replace(cat.Mention)

	desc := strings.NewReplacer(
		"{mention}", user.Mention(),
		"{channel}", ch.Mention(),
		"{here}", "@here",
		"{everyone}", "@everyone",
	).Replace(cat.Description)

	e := discord.Embed{
		Title: cat.Name,
		Color: ctx.Router.EmbedColor,
	}
	if cat.CanCreatorClose {
		e.Description = desc + "\n\nReact with :x: to close this ticket."
	} else {
		e.Description = desc
	}

	if e.Description != "" {
		e.Title = ""
	}

	m, err := ctx.State.SendMessage(ch.ID, mention, e)
	if err != nil {
		bot.Sugar.Errorf("Error sending message: %v", err)
		return err
	}

	err = ctx.State.PinMessage(m.ChannelID, m.ID, "")
	if err == nil {
		msgs, err := ctx.State.Messages(ch.ID, 100)
		if err == nil {
			for _, m := range msgs {
				if m.Author.ID == ctx.Bot.ID && m.Content == "" && len(m.Embeds) == 0 {
					ctx.State.DeleteMessage(m.ChannelID, m.ID, "")
					break
				}
			}
		}
	}

	if cat.CanCreatorClose {
		ctx.State.React(m.ChannelID, m.ID, "❌")
	}

	_, err = bot.DB.Pool.Exec(context.Background(), `insert into triggers
(message_id, emoji, command)
values ($1, $2, $3) on conflict (message_id, emoji) do
update set command = $3`, m.ID, "❌", []string{"tickets", "delete"})
	return err
}
