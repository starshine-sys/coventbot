package gatekeeper

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) settings(ctx *bcr.Context) (err error) {
	var b strings.Builder

	settings, err := bot.serverSettings(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if !settings.MemberRole.IsValid() {
		_, err = ctx.Send("The gateway is currently disabled.")
		return err
	}
	b.WriteString(fmt.Sprintf("Users will be given the <@&%v> role upon verification.\n", settings.MemberRole))

	if !settings.WelcomeChannel.IsValid() || settings.WelcomeMessage == "" {
		b.WriteString("No message will be sent upon verification.")
	} else {
		b.WriteString(fmt.Sprintf("This message will be sent in <#%v> upon verification:\n```%v```", settings.WelcomeChannel, settings.WelcomeMessage))
	}

	_, err = ctx.Send("", discord.Embed{
		Title:       "Gateway settings",
		Description: b.String(),
		Color:       ctx.Router.EmbedColor,
	})
	return err
}

func (bot *Bot) setChannel(ctx *bcr.Context) (err error) {
	var id discord.ChannelID

	if ctx.RawArgs == "-clear" {
		id = 0
	} else {
		ch, err := ctx.ParseChannel(ctx.RawArgs)
		if err != nil {
			_, err = ctx.Send("Channel not found.")
			return err
		}

		if ch.GuildID != ctx.Message.GuildID {
			_, err = ctx.Send("The given channel isn't in this server.")
			return err
		}

		id = ch.ID
	}

	settings, err := bot.serverSettings(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if settings.WelcomeChannel == id {
		_, err = ctx.Send("The given channel is already the welcome channel.")
		return err
	}

	settings.WelcomeChannel = id
	err = bot.setSettings(settings)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if id == 0 {
		_, err = ctx.Send("Welcome channel reset.")
		return
	}
	_, err = ctx.Sendf("Welcome channel changed to %v.", id.Mention())
	return
}

func (bot *Bot) setLog(ctx *bcr.Context) (err error) {
	var id discord.ChannelID

	if ctx.RawArgs == "-clear" {
		id = 0
	} else {
		ch, err := ctx.ParseChannel(ctx.RawArgs)
		if err != nil {
			_, err = ctx.Send("Channel not found.")
			return err
		}

		if ch.GuildID != ctx.Message.GuildID {
			_, err = ctx.Send("The given channel isn't in this server.")
			return err
		}

		id = ch.ID
	}

	settings, err := bot.serverSettings(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if settings.GatekeeperLog == id {
		_, err = ctx.Send("The given channel is already the log channel.")
		return err
	}

	settings.GatekeeperLog = id
	err = bot.setSettings(settings)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if id == 0 {
		_, err = ctx.Send("Log channel reset.")
		return
	}
	_, err = ctx.Sendf("Log channel changed to %v.", id.Mention())
	return
}

func (bot *Bot) setMessage(ctx *bcr.Context) (err error) {
	settings, err := bot.serverSettings(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if settings.WelcomeMessage == ctx.RawArgs {
		_, err = ctx.Send("The given welcome message is already set.")
		return err
	}

	settings.WelcomeMessage = ctx.RawArgs
	err = bot.setSettings(settings)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("Welcome message changed!")
	return
}

func (bot *Bot) setRole(ctx *bcr.Context) (err error) {
	var id discord.RoleID

	if ctx.RawArgs == "-clear" {
		id = 0
	} else {
		role, err := ctx.ParseRole(ctx.RawArgs)
		if err != nil {
			_, err = ctx.Send("Role not found.")
			return err
		}

		id = role.ID
	}

	settings, err := bot.serverSettings(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if settings.MemberRole == id {
		_, err = ctx.Send("The given role is already the member role.")
		return err
	}

	settings.MemberRole = id
	err = bot.setSettings(settings)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if id == 0 {
		_, err = ctx.Send("Member role reset.")
		return
	}
	_, err = ctx.Sendf("Member role changed to %v.", id.Mention())
	return
}
