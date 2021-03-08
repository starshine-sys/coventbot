package static

import (
	"net/http"
	"strings"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) echo(ctx *bcr.Context) (err error) {
	return bot.echoInner(ctx, ctx.Channel)
}

func (bot *Bot) echoTo(ctx *bcr.Context) (err error) {
	ctx.RawArgs = strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))

	ch, err := ctx.ParseChannel(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("Could not find that channel.", nil)
		return
	}

	if ch.GuildID != ctx.Message.GuildID {
		_, err = ctx.Send("That channel isn't in this server.", nil)
		return
	}

	return bot.echoInner(ctx, ch)
}

func (bot *Bot) echoInner(ctx *bcr.Context, ch *discord.Channel) (err error) {
	if ctx.RawArgs == "" && len(ctx.Message.Attachments) == 0 {
		_, err = ctx.Send("You need to give me something to say!", nil)
		return
	}

	perms, err := ctx.Session.Permissions(ch.ID, ctx.Author.ID)
	if err != nil {
		_, err = ctx.Sendf("Could not check your permissions: %v", err)
		return
	}

	msg := ctx.NewMessage(ch.ID).Content(ctx.RawArgs)
	var am *api.AllowedMentions

	if perms.Has(discord.PermissionMentionEveryone) {
		am = &api.AllowedMentions{
			Parse: []api.AllowedMentionType{api.AllowRoleMention, api.AllowEveryoneMention, api.AllowUserMention},
		}
	} else {
		am = &api.AllowedMentions{
			Parse: []api.AllowedMentionType{api.AllowUserMention},
			Roles: []discord.RoleID{},
		}
		roles, err := ctx.Session.Roles(ch.GuildID)
		if err == nil {
			for _, r := range roles {
				am.Roles = append(am.Roles, r.ID)
			}
		}
	}

	for _, a := range ctx.Message.Attachments {
		resp, err := http.Get(a.Filename)
		if err != nil {
			_, err = ctx.Sendf("Error occurred: %v", err)
			return err
		}
		defer resp.Body.Close()
		msg.AddFile(a.Filename, resp.Body)
	}

	_, err = msg.AllowedMentions(am).Send()
	return
}
