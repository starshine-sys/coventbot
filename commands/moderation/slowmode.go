package moderation

import (
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) discordSlowmode(ctx *bcr.Context) (err error) {
	ch := ctx.Channel
	if len(ctx.Args) > 1 {
		ch, err = ctx.ParseChannel(ctx.Args[1])
		if err != nil {
			_, err = ctx.Replyc(bcr.ColourRed, "Could not find the given channel.")
			return
		}
		if ch.GuildID != ctx.Channel.GuildID {
			_, err = ctx.Replyc(bcr.ColourRed, "The given channel must be in this server.")
			return
		}
	}

	duration, err := time.ParseDuration(ctx.Args[0])
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "You didn't give a valid slowmode duration.")
		return
	}

	if duration > 6*time.Hour || duration < 0 {
		_, err = ctx.Replyc(bcr.ColourRed, "The given duration must be between 0 seconds and 6 hours.")
		return
	}

	err = ctx.State.ModifyChannel(ch.ID, api.ModifyChannelData{
		UserRateLimit: option.NewNullableUint(uint(duration.Seconds())),
	})
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "There was an error changing the slowmode for the given channel. Are you sure I have the **Manage Channel** permission?")
		return
	}

	_, err = ctx.Reply("Changed the slowmode for %v to %s.", ch.Mention(), duration)
	return
}
