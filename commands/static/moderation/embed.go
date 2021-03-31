package moderation

import (
	"encoding/json"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) embed(ctx *bcr.Context) (err error) {
	return bot.embedInner(ctx, ctx.Channel.ID, []byte(ctx.RawArgs))
}

func (bot *Bot) embedTo(ctx *bcr.Context) (err error) {
	ch, err := ctx.ParseChannel(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("Channel not found.", nil)
	}
	if ch.GuildID != ctx.Message.GuildID {
		_, err = ctx.Send("That channel is not in this server.", nil)
		return
	}

	g, err := ctx.State.Guild(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Send("Error getting guild info.", nil)
		return
	}

	if !discord.CalcOverwrites(*g, *ctx.Channel, *ctx.Member).Has(discord.PermissionViewChannel | discord.PermissionSendMessages) {
		_, err = ctx.Send("You do not have permission to send messages in that channel.", nil)
		return
	}

	args := strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))

	err = bot.embedInner(ctx, ch.ID, []byte(args))
	if err != nil {
		_, err = ctx.Sendf("Error sending message: %v", err)
		return
	}

	_, err = ctx.Send("✅", nil)
	return
}

func (bot *Bot) embedInner(ctx *bcr.Context, ch discord.ChannelID, input []byte) (err error) {
	var e discord.Embed
	err = json.Unmarshal(input, &e)
	if err != nil {
		ctx.Sendf("Error: %v", err)
		return
	}

	_, err = ctx.NewMessage(ch).Embed(&e).Send()
	return
}
