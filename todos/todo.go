package todos

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) todo(ctx *bcr.Context) (err error) {
	todoCh := bot.getChannel(ctx.Author.ID)
	if !todoCh.IsValid() {
		_, err = ctx.Replyc(bcr.ColourRed, "You don't have a todo channel set! Set one with `%vtodo channel`.", ctx.Prefix)
		return
	}

	jumpLink := "https://discord.com/channels/"
	if !ctx.Message.GuildID.IsValid() {
		jumpLink += "@me/"
	} else {
		jumpLink += ctx.Message.GuildID.String() + "/"
	}
	jumpLink += fmt.Sprintf("%v/%v", ctx.Message.ChannelID, ctx.Message.ID)

	e := discord.Embed{
		Title:       "Todo",
		Color:       bcr.ColourBlurple,
		Description: ctx.RawArgs,

		Fields: []discord.EmbedField{{
			Name:  "Source",
			Value: fmt.Sprintf("[Jump!](%v)", jumpLink),
		}},

		Timestamp: discord.NowTimestamp(),
	}
	if ctx.Guild == nil {
		e.Author = &discord.EmbedAuthor{
			Name: fmt.Sprintf("DM with %v", ctx.Author.Username),
			Icon: ctx.Author.AvatarURL(),
		}
	} else {
		e.Author = &discord.EmbedAuthor{
			Name: ctx.Guild.Name,
			Icon: ctx.Guild.IconURL(),
		}
	}

	msg, err := ctx.State.SendEmbed(todoCh, e)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "I couldn't send the todo message. Are you sure I have write perms in your todo channel?")
		return
	}

	t := Todo{
		UserID:      ctx.Author.ID,
		Description: ctx.RawArgs,

		OrigMID:       ctx.Message.ID,
		OrigChannelID: ctx.Message.ChannelID,
		OrigServerID:  ctx.Message.GuildID,

		MID:       msg.ID,
		ChannelID: msg.ChannelID,
		ServerID:  msg.GuildID,
	}

	out, err := bot.newTodo(t)
	if err != nil {
		return bot.Report(ctx, err)
	}

	e.Footer = &discord.EmbedFooter{
		Text: fmt.Sprintf("ID: %v", out.ID),
	}

	_, err = ctx.State.EditEmbed(msg.ChannelID, msg.ID, e)
	if err != nil {
		return bot.Report(ctx, err)
	}

	ctx.State.React(msg.ChannelID, msg.ID, "✅")

	_, err = ctx.Reply("New todo added with ID #%v.\n[Link](https://discord.com/channels/%v/%v/%v)", out.ID, msg.GuildID, msg.ChannelID, msg.ID)
	return
}
