package chanmirror

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jackc/pgx/v4"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) set(ctx *bcr.Context) (err error) {
	if len(ctx.Args) == 0 {
		m, err := bot.mirrors(ctx.Message.GuildID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		if len(m) == 0 {
			_, err = ctx.Reply("No channels are being mirrored.")
			return err
		}

		var s []string
		for _, ch := range m {
			s = append(s, fmt.Sprintf("%v → %v\n", ch.FromChannel.Mention(), ch.ToChannel.Mention()))
		}

		_, err = ctx.PagedEmbed(
			bcr.StringPaginator("Channel mirrors", bcr.ColourBlurple, s, 10), false,
		)
		return err
	}

	if len(ctx.Args) < 2 {
		_, err = ctx.Replyc(bcr.ColourRed, "You must give both a source and destination channel.")
		return
	}

	src, err := ctx.ParseChannel(ctx.Args[0])
	if err != nil || src.GuildID != ctx.Message.GuildID {
		_, err = ctx.Replyc(bcr.ColourRed, "Source channel not found.")
		return
	}

	var dest discord.ChannelID
	if ctx.Args[1] == "-clear" || ctx.Args[1] == "--clear" || ctx.Args[1] == "clear" {
		dest = 0
	} else {
		destCh, err := ctx.ParseChannel(ctx.Args[1])
		if err != nil || destCh.GuildID != ctx.Message.GuildID {
			_, err = ctx.Replyc(bcr.ColourRed, "Destination channel not found.")
			return err
		}

		dest = destCh.ID
	}

	if dest == 0 {
		_, err = bot.DB.Pool.Exec(context.Background(), "delete from channel_mirror where server_id = $1 and from_channel = $2", ctx.Message.GuildID, src.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}
		_, err = ctx.Reply("Mirror for %v removed!", src.Mention())
		return
	}

	m, err := bot.mirrorFor(src.ID)
	if err != nil {
		if errors.Cause(err) != pgx.ErrNoRows {
			return bot.Report(ctx, err)
		}
	} else {
		if m.ToChannel == dest {
			_, err = ctx.Replyc(bcr.ColourRed, "The channel you're trying to mirror to is already set for that channel.")
			return
		}
	}

	name := src.Name
	if len(src.Name) > 40 {
		name = src.Name[:40] + "..."
	}

	wh, err := ctx.State.CreateWebhook(dest, api.CreateWebhookData{
		Name: fmt.Sprintf("Channel mirror for #%v", name),
	})
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "I couldn't create a webhook in %v. Please make sure I have the Manage Webhooks permission in that channel, and the channel hasn't hit the maximum number of webhooks.")
		return
	}

	err = bot.setMirror(Mirror{
		ServerID:    ctx.Message.GuildID,
		FromChannel: src.ID,
		ToChannel:   dest,
		WebhookID:   wh.ID,
		Token:       wh.Token,
	})
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Mirror channel set! Messages in %v will now be mirrored to %v.", src.Mention(), dest.Mention())
	return
}
