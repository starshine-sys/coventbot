package admin

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/gateway/shard"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) status(ctx *bcr.Context) (err error) {
	s := bot.Settings()

	if ctx.RawArgs == "" {
		_, err = ctx.Send("", discord.Embed{
			Title: "Status",
			Description: fmt.Sprintf(`The bot's status is currently set to `+"`%v`"+`
Available values are:
`+"`online`: online\n`idle`: idle\n`dnd`: do not disturb", s.Status),
			Color: ctx.Router.EmbedColor,
		})
		return
	}

	status := strings.ToLower(ctx.RawArgs)
	if status == "do-not-disturb" || status == "donotdisturb" {
		status = "dnd"
	}
	if status != "online" && status != "idle" && status != "dnd" {
		_, err = ctx.Send("No valid status given. Valid statuses are: `online`, `idle,` `dnd`.")
		return
	}

	s.Status = gateway.Status(status)

	err = bot.SetSettings(s)
	if err != nil {
		return bot.Report(ctx, err)
	}

	go bot.Router.ShardManager.ForEach(func(s shard.Shard) {
		state := s.(*state.State)

		bot.updateStatus(state)
	})

	_, err = ctx.Send("", discord.Embed{
		Description: fmt.Sprintf("Set status to `%v`.", status),
		Color:       ctx.Router.EmbedColor,
	})
	return
}
