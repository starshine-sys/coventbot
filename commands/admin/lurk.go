package admin

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/voice"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) lurk(ctx *bcr.Context) error {
	ch, err := ctx.ParseChannel(ctx.RawArgs)
	if err != nil {
		return ctx.SendfX("error: %v", err)
	}

	if ch.Type != discord.GuildVoice && ch.Type != discord.GuildStageVoice {
		return ctx.SendfX("channel %v must be a voice channel", ch.Mention())
	}

	vs, err := voice.NewSession(ctx.State)
	if err != nil {
		bot.Sugar.Error("error creating voice session:", err)
		return ctx.SendfX("error: %v", err)
	}

	err = vs.JoinChannel(context.Background(), ch.ID, true, true)
	if err != nil {
		bot.Sugar.Error("error connecting to voice channel:", err)
		return ctx.SendfX("error: %v", err)
	}

	return ctx.SendfX("Now lurking in %v!", ch.Mention())
}
