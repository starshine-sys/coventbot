package mirror

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) cmdImport(ctx *bcr.Context) (err error) {
	modlog, err := ctx.ParseChannel(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "Channel not found.")
		return
	}
	if modlog.GuildID != ctx.Guild.ID || modlog.Type != discord.GuildText {
		_, err = ctx.Replyc(bcr.ColourRed, "Channel not found.")
		return
	}

	msgs, err := ctx.State.Session.Messages(modlog.ID, 0)
	if err != nil {
		return bot.Report(ctx, err)
	}

	var oldCount int
	bot.DB.Pool.QueryRow(context.Background(), "select count(*) from mod_log where server_id = $1", ctx.Guild.ID).Scan(&oldCount)

	for _, m := range msgs {
		m.GuildID = ctx.Message.GuildID
		bot.messageCreate(&gateway.MessageCreateEvent{Message: m})
	}

	var newCount int
	bot.DB.Pool.QueryRow(context.Background(), "select count(*) from mod_log where server_id = $1", ctx.Guild.ID).Scan(&newCount)

	if oldCount == newCount {
		_, err = ctx.Replyc(bcr.ColourRed, "No mod logs were imported.")
		return
	}

	_, err = ctx.Reply("Success, imported %v entries!", newCount-oldCount)
	return
}
