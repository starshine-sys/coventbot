// SPDX-License-Identifier: AGPL-3.0-only
package admin

import (
	"context"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/voice"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) lurk(ctx *bcr.Context) error {
	ch, err := ctx.ParseChannel(ctx.Args[0])
	if err != nil {
		return ctx.SendfX("error: %v", err)
	}

	if ch.Type != discord.GuildVoice && ch.Type != discord.GuildStageVoice {
		return ctx.SendfX("channel %v must be a voice channel", ch.Mention())
	}

	vs, ok := bot.VoiceSessions.Get(ctx.Message.GuildID)
	if !ok {
		vs, err = voice.NewSession(ctx.State)
		if err != nil {
			bot.Sugar.Error("error creating voice session:", err)
			return ctx.SendfX("error: %v", err)
		}

		bot.VoiceSessions.Set(ctx.Message.GuildID, vs)
	}

	unmute, _ := ctx.Flags.GetBool("unmute")
	undeafen, _ := ctx.Flags.GetBool("undeafen")

	err = vs.JoinChannel(context.Background(), ch.ID, !unmute, !undeafen)
	if err != nil {
		bot.Sugar.Error("error connecting to voice channel:", err)
		return ctx.SendfX("error: %v", err)
	}

	return ctx.SendfX("Now lurking in %v!", ch.Mention())
}

func (bot *Bot) unlurk(ctx *bcr.Context) error {
	id := ctx.Message.GuildID
	if len(ctx.Args) > 0 {
		sf, err := discord.ParseSnowflake(ctx.RawArgs)
		if err == nil {
			id = discord.GuildID(sf)
		}
	}

	vs, ok := bot.VoiceSessions.Get(id)
	if !ok {
		return ctx.SendX("There isn't a voice session for this guild.")
	}

	err := vs.Leave(context.Background())
	if err != nil {
		return ctx.SendfX("error: %v", err)
	}

	bot.VoiceSessions.Remove(id)
	return ctx.SendX("Stopped lurking in this guild!")
}
