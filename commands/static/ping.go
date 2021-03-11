package static

import (
	"fmt"
	"time"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) ping(ctx *bcr.Context) (err error) {
	// this will return 0ms in the first minute after the bot is restarted
	// can't do much about that though
	heartbeat := ctx.Session.Gateway.PacerLoop.EchoBeat.Time().Sub(ctx.Session.Gateway.PacerLoop.SentBeat.Time()).Round(time.Millisecond)

	t := time.Now()
	m, err := ctx.Send("Pong!", nil)
	if err != nil {
		return err
	}

	_, err = ctx.Edit(m, fmt.Sprintf(
		"Ping: %v | Message: %v", heartbeat, time.Since(t).Round(time.Millisecond),
	), nil)
	return err
}
