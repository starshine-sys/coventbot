package reactroles

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) clear(ctx *bcr.Context) (err error) {
	if len(ctx.Args) == 0 {
		m, err := ctx.SendEmbed(bcr.SED{
			Message: "Warning: this will delete **all** reaction roles for this server. Are you sure you want to continue?",
			Color:   bcr.ColourRed,
		})
		if err != nil {
			return err
		}

		yes, timeout := ctx.YesNoHandler(*m, ctx.Author.ID)
		if !yes || timeout {
			_, err = ctx.Edit(m, "", &discord.Embed{
				Description: "Cancelled.",
				Color:       bcr.ColourBlurple,
			})
			return err
		}

		ct, err := bot.DB.Pool.Exec(context.Background(), "delete from react_roles where server_id = $1", ctx.Message.GuildID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		_, err = ctx.Edit(m, "", &discord.Embed{
			Description: fmt.Sprintf("Success! Deleted reaction roles from %v message(s)", ct.RowsAffected()),
			Color:       bcr.ColourBlurple,
		})
		return err
	}

	m, err := ctx.ParseMessage(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("Couldn't parse your input as a message.", nil)
		return
	}
	if m.GuildID != ctx.Message.GuildID {
		_, err = ctx.Send("The given message isn't in this server.", nil)
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "delete from react_roles where message_id = $1", m.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("Removed reaction roles from that message.", nil)
	return
}
