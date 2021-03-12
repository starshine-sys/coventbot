package info

import (
	"fmt"
	"time"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) ping(ctx *bcr.Context) (err error) {
	t := time.Now()
	// this will return 0ms in the first minute after the bot is restarted
	// can't do much about that though
	heartbeat := ctx.Session.Gateway.PacerLoop.EchoBeat.Time().Sub(ctx.Session.Gateway.PacerLoop.SentBeat.Time()).Round(time.Millisecond)

	s := fmt.Sprintf("🏓 **Pong!**\nHeartbeat: %v", heartbeat)
	m, err := ctx.Send(s, nil)
	if err != nil {
		return err
	}

	latency := time.Since(t).Round(time.Millisecond)

	_, err = ctx.Edit(m, fmt.Sprintf("%v\nMessage: %v", s, latency), nil)
	return err
}
