package approval

import (
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) setChannel(ctx *bcr.Context) (err error) {
	var id discord.ChannelID

	if ctx.RawArgs == "-clear" {
		id = 0
	} else {
		ch, err := ctx.ParseChannel(ctx.RawArgs)
		if err != nil {
			_, err = ctx.Send("Channel not found.", nil)
			return err
		}

		if ch.GuildID != ctx.Message.GuildID {
			_, err = ctx.Send("The given channel isn't in this server.", nil)
			return err
		}

		id = ch.ID
	}

	settings, err := bot.serverSettings(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	if settings.ApproveWelcomeChannel == id {
		_, err = ctx.Send("The given channel is already the welcome channel.", nil)
		return err
	}

	settings.ApproveWelcomeChannel = id
	err = bot.setSettings(settings)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	if id == 0 {
		_, err = ctx.Send("Welcome channel reset.", nil)
		return
	}
	_, err = ctx.Sendf("Welcome channel changed to %v.", id.Mention())
	return
}

func (bot *Bot) setMessage(ctx *bcr.Context) (err error) {
	settings, err := bot.serverSettings(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	if settings.ApproveWelcomeMessage == ctx.RawArgs {
		_, err = ctx.Send("The given welcome message is already set.", nil)
		return err
	}

	settings.ApproveWelcomeMessage = ctx.RawArgs
	err = bot.setSettings(settings)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	_, err = ctx.Send("Welcome message changed!", nil)
	return
}

func (bot *Bot) setAddRoles(ctx *bcr.Context) (err error) {
	var roles []uint64
	var roleNames []string

	if ctx.RawArgs == "-clear" {
		roles = []uint64{}
	} else {
		addRoles, n := ctx.GreedyRoleParser(ctx.Args)
		if n == 0 {
			_, err = ctx.Sendf("No roles found!")
			return
		}

		for _, r := range addRoles {
			roles = append(roles, uint64(r.ID))
			roleNames = append(roleNames, r.Name)
		}
	}

	settings, err := bot.serverSettings(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	settings.ApproveAddRoles = roles
	err = bot.setSettings(settings)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	if len(roles) == 0 {
		_, err = ctx.Send("Cleared approval add roles.", nil)
	} else {
		_, err = ctx.Sendf("Now adding the role(s) %v on approval.", strings.Join(roleNames, ", "))
	}
	return
}

func (bot *Bot) setRemoveRoles(ctx *bcr.Context) (err error) {
	var roles []uint64
	var roleNames []string

	if ctx.RawArgs == "-clear" {
		roles = []uint64{}
	} else {
		addRoles, n := ctx.GreedyRoleParser(ctx.Args)
		if n == 0 {
			_, err = ctx.Sendf("No roles found!")
			return
		}

		for _, r := range addRoles {
			roles = append(roles, uint64(r.ID))
			roleNames = append(roleNames, r.Name)
		}
	}

	settings, err := bot.serverSettings(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	settings.ApproveRemoveRoles = roles
	err = bot.setSettings(settings)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

	if len(roles) == 0 {
		_, err = ctx.Send("Cleared approval remove roles.", nil)
	} else {
		_, err = ctx.Sendf("Now removing the role(s) %v on approval.", strings.Join(roleNames, ", "))
	}
	return
}
