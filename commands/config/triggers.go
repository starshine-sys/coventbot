package config

import (
	"context"
	"fmt"
	"strings"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/starshine-sys/bcr"
)

// Trigger ...
type Trigger struct {
	MessageID discord.MessageID

	Emoji string

	Command []string
}

func (bot *Bot) addTrigger(ctx *bcr.Context) (err error) {
	msg, err := ctx.ParseMessage(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("Couldn't find that message.", nil)
		return
	}

	if msg.GuildID != ctx.Message.GuildID {
		_, err = ctx.Send("That message isn't in this server.", nil)
		return
	}

	if perms, _ := ctx.State.Permissions(msg.ChannelID, ctx.Author.ID); !perms.Has(discord.PermissionManageMessages) {
		_, err = ctx.Send("You're not a mod, you can't do that!", nil)
		return
	}

	emoji := ctx.Args[1]

	command := ctx.Args[2:]

	err = ctx.State.React(msg.ChannelID, msg.ID, discord.APIEmoji(emoji))
	if err != nil {
		_, err = ctx.Send("I couldn't react to the given message with that emoji. Either I don't have the **Add Reactions** permission in the channel, or you didn't give a valid emoji.", nil)
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), `insert into triggers
(message_id, emoji, command)
values ($1, $2, $3) on conflict (message_id, emoji) do
update set command = $3`, msg.ID, emoji, command)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.SendEmbed(bcr.SED{
		Message: fmt.Sprintf("Added %v as a trigger on [the given message](https://discord.com/channels/%v/%v/%v), pointing to ``%v``.", emoji, msg.GuildID, msg.ChannelID, msg.ID, strings.Join(command, " ")),
	})
	return err
}

func (bot *Bot) delTrigger(ctx *bcr.Context) (err error) {
	msg, err := ctx.ParseMessage(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("Couldn't find that message.", nil)
		return
	}

	if msg.GuildID != ctx.Message.GuildID {
		_, err = ctx.Send("That message isn't in this server.", nil)
		return
	}

	if perms, _ := ctx.State.Permissions(msg.ChannelID, ctx.Author.ID); !perms.Has(discord.PermissionManageMessages) {
		_, err = ctx.Send("You're not a mod, you can't do that!", nil)
		return
	}

	emoji := ctx.Args[1]

	_, err = bot.DB.Pool.Exec(context.Background(), "delete from triggers where message_id = $1 and emoji = $2", msg.ID, emoji)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("Trigger deleted.", nil)
	return
}

func (bot *Bot) triggerReactionAdd(ev *gateway.MessageReactionAddEvent) {
	if !ev.GuildID.IsValid() || ev.Member == nil {
		return
	}

	if ev.Member.User.Bot {
		return
	}

	var t Trigger
	err := pgxscan.Get(context.Background(), bot.DB.Pool, &t, "select * from triggers where message_id = $1 and emoji = $2", ev.MessageID, ev.Emoji.APIString())
	if err != nil {
		if errors.Cause(err) != pgx.ErrNoRows {
			bot.Sugar.Errorf("Error getting trigger: %v", err)
		}
		return
	}

	if len(t.Command) == 0 {
		return
	}

	err = bot.State.DeleteUserReaction(ev.ChannelID, ev.MessageID, ev.UserID, ev.Emoji.APIString())
	if err != nil {
		bot.Sugar.Errorf("Error deleting reaction: %v", err)
	}

	// create a new context, yes this is messy as hell
	ctx := &bcr.Context{
		Router: bot.Router,
		State:  bot.Router.State,
		Bot:    bot.Router.Bot,
	}

	msg, err := bot.State.Message(ev.ChannelID, ev.MessageID)
	if err != nil {
		bot.Sugar.Errorf("Error getting message: %v", err)
		return
	}

	ctx.Message = *msg
	ctx.Member = ev.Member
	ctx.Author = ev.Member.User

	ctx.Command = t.Command[0]

	ctx.Args = []string{}
	ctx.RawArgs = ""
	if len(t.Command) > 1 {
		ctx.Args = t.Command[1:]
		ctx.RawArgs = strings.Join(t.Command[1:], " ")
	}
	ctx.InternalArgs = ctx.Args

	ctx.AdditionalParams = make(map[string]interface{})

	channel, err := bot.State.Channel(ev.ChannelID)
	if err != nil {
		return
	}
	ctx.Channel = channel

	err = bot.Router.Execute(ctx)
	if err != nil {
		bot.Sugar.Errorf("Error executing command `%v`: %v", t.Command, err)
	}
}

func (bot *Bot) triggerMessageDelete(ev *gateway.MessageDeleteEvent) {
	bot.DB.Pool.Exec(context.Background(), "delete from triggers where message_id = $1", ev.ID)
}
