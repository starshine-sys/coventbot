package tickets

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) delete(ctx *bcr.Context) (err error) {
	ch := struct {
		OwnerID         discord.UserID
		CanCreatorClose bool
		LogChannel      discord.ChannelID
	}{}

	err = bot.DB.Pool.QueryRow(context.Background(), "select t.owner_id, c.can_creator_close, c.log_channel from tickets as t, ticket_categories as c where t.category_id = c.category_id and t.channel_id = $1", ctx.Channel.ID).Scan(&ch.OwnerID, &ch.CanCreatorClose, &ch.LogChannel)
	if err != nil {
		bot.Sugar.Errorf("Error: %v", err)
		return bot.Report(ctx, err)
	}

	if ch.CanCreatorClose && ch.OwnerID == ctx.Author.ID {
		goto canClose
	}

	if perms, _ := ctx.State.Permissions(ctx.Channel.ID, ctx.Author.ID); perms.Has(discord.PermissionManageMessages) {
		goto canClose
	}

	_, err = ctx.Send("You don't have permission to close this ticket.", nil)
	return

canClose:

	msgs, err := ctx.State.MessagesAfter(ctx.Channel.ID, 0, 0)
	if err != nil {
		_, err = ctx.Send("I couldn't fetch all messages in this channel, aborting.", nil)
		return
	}

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].ID < msgs[j].ID
	})

	var buf []string

	owner, err := ctx.State.User(ch.OwnerID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	buf = append(buf, fmt.Sprintf(`#%v (%v)
Owner: %v#%v (%v)
Created %v, closed %v
Messages: %v
`, ctx.Channel.Name, ctx.Channel.ID, owner.Username, owner.Discriminator, owner.ID, ctx.Channel.ID.Time().UTC().Format("2006-01-02 15:05:05"), time.Now().UTC().Format("2006-01-02 15:05:05"), len(msgs)))

	users := []string{}

	for _, m := range msgs {
		b := fmt.Sprintf(`--------------------------------------------------------------------------------
[%v] %v#%v (%v)
%v`, m.Timestamp.Time().UTC().Format("2006-01-02 15:05:05"), m.Author.Username, m.Author.Discriminator, m.Author.ID, m.Content)

		var isInUsers bool
		for _, u := range users {
			if m.Author.Mention() == u {
				isInUsers = true
				break
			}
		}
		if !isInUsers {
			users = append(users, m.Author.Mention())
		}

		if len(m.Embeds) > 0 {
			bt, err := json.Marshal(&m.Embeds[0])
			if err == nil {
				b += "\n\n" + string(bt)
			}
		}

		if len(m.Attachments) > 0 {
			b += "\n\nAttachments:\n"
			for _, a := range m.Attachments {
				b += a.URL + "\n"
			}
		}

		buf = append(buf, b+"\n")
	}

	text := strings.Join(buf, "\n")

	e := &discord.Embed{
		Title: fmt.Sprintf("%v closed", ctx.Channel.Name),

		Fields: []discord.EmbedField{
			{
				Name:   "Owner",
				Value:  fmt.Sprintf("%v#%v\n%v", owner.Username, owner.Discriminator, owner.Mention()),
				Inline: true,
			},
			{
				Name:   "Messages",
				Value:  fmt.Sprint(len(msgs)),
				Inline: true,
			},
			{
				Name:   "Participants",
				Value:  strings.Join(users, "\n"),
				Inline: true,
			},
		},

		Color: ctx.Router.EmbedColor,
	}

	_, err = ctx.NewMessage(ch.LogChannel).Embed(e).AddFile(
		fmt.Sprintf("transcript-%v.txt", ctx.Channel.Name), strings.NewReader(text),
	).Send()
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("", &discord.Embed{
		Description: "Ticket will be deleted in 5 seconds.",
		Color:       bcr.ColourRed,
	})
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	bot.State.DeleteChannel(ctx.Channel.ID)

	_, err = bot.DB.Pool.Exec(context.Background(), "delete from tickets where channel_id = $1", ctx.Channel.ID)
	return err
}
