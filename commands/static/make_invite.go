package static

import (
	"flag"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/utils/json/option"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) makeInvite(ctx *bcr.Context) (err error) {
	var existing bool

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.BoolVar(&existing, "existing", false, "Reuse an existing invite if one exists.")
	fs.Parse(ctx.Args)

	ctx.Args = fs.Args()

	channel := ctx.Channel
	if len(ctx.Args) > 0 {
		channel, err = ctx.ParseChannel(ctx.Args[0])
		if err != nil {
			_, err = ctx.Send("Channel not found!", nil)
			return
		}
	}

	inv, err := ctx.Session.CreateInvite(channel.ID, api.CreateInviteData{
		MaxAge:    option.NewUint(0),
		MaxUses:   0,
		Temporary: false,
		Unique:    !existing,
	})
	if err != nil {
		_, err = ctx.Send("There was an error creating the invite. Please make sure I have permission to create an invite.", nil)
		return
	}

	_, err = ctx.Sendf("Created new invite **%v** for %v.\nLink: https://discord.gg/%v", inv.Code, inv.Channel.Mention(), inv.Code)
	return
}
