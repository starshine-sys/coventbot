package tags

import (
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) tag(ctx *bcr.Context) (err error) {
	t, err := bot.DB.GetTag(ctx.Message.GuildID, ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("No tag with that name found.")
		return
	}

	tr := false

	if len(ctx.Message.Mentions) > 0 {
		tr = true
	}

	data := api.SendMessageData{
		Content: t.Response,
		AllowedMentions: &api.AllowedMentions{
			Parse:       []api.AllowedMentionType{},
			RepliedUser: option.Bool(&tr),
		},
	}

	if ctx.Message.Reference != nil {
		data.Reference = &discord.MessageReference{
			MessageID: ctx.Message.Reference.MessageID,
		}
	}

	_, err = ctx.State.SendMessageComplex(ctx.Message.ChannelID, data)
	return err
}
