package highlights

import (
	"fmt"
	"sort"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) toggleHl(ctx *bcr.Context) (err error) {
	conf, err := bot.hlConfig(ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	conf.HlEnabled = !conf.HlEnabled

	err = bot.setConfig(conf)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if conf.HlEnabled {
		_, err = ctx.Reply("Enabled highlights for this server!")
	} else {
		_, err = ctx.Reply("Disabled highlights for this server!")
	}
	return
}

func (bot *Bot) modBlockHl(ctx *bcr.Context) (err error) {
	conf, err := bot.hlConfig(ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	ch, err := ctx.ParseChannel(ctx.RawArgs)
	if err != nil || ch.GuildID != ctx.Guild.ID || (ch.Type != discord.GuildText && ch.Type != discord.GuildNews && ch.Type != discord.GuildCategory) {
		_, err = ctx.Replyc(bcr.ColourRed, "Channel not found.")
		return
	}

	for _, id := range conf.Blocked {
		if discord.ChannelID(id) == ch.ID {
			_, err = ctx.Replyc(bcr.ColourRed, "That channel is already blocked.")
			return
		}
		if discord.ChannelID(id) == ch.CategoryID {
			_, err = ctx.Replyc(bcr.ColourRed, "That channel is technically already blocked (in blocked category <#%v>).", id)
			return
		}
	}

	conf.Blocked = append(conf.Blocked, uint64(ch.ID))
	err = bot.setConfig(conf)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if ch.Type == discord.GuildCategory {
		_, err = ctx.Replyc(bcr.ColourGreen, "Blocked highlights in category %v!", ch.Name)
	} else {
		_, err = ctx.Replyc(bcr.ColourGreen, "Blocked highlights in %v!", ch.Mention())
	}
	return
}

func (bot *Bot) showHlConfig(ctx *bcr.Context) (err error) {
	conf, err := bot.hlConfig(ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	e := discord.Embed{
		Title:       "Highlight config for " + ctx.Guild.Name,
		Color:       ctx.Router.EmbedColor,
		Description: fmt.Sprintf("Highlights enabled? %v", conf.HlEnabled),
	}

	channels, err := ctx.State.Channels(ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	var blockedChannels []discord.Channel

	for _, id := range conf.Blocked {
		for _, ch := range channels {
			if ch.ID == discord.ChannelID(id) {
				blockedChannels = append(blockedChannels, ch)
				break
			}
		}
	}

	sort.Slice(blockedChannels, func(i, j int) bool {
		return blockedChannels[i].Name < blockedChannels[j].Name
	})

	var channelString, categories string
	for _, ch := range blockedChannels {
		if ch.Type == discord.GuildCategory {
			categories += ch.Name + "\n"
		} else {
			channelString += ch.Mention() + "\n"
		}
	}

	if categories != "" {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Blocked categories",
			Value: categories,
		})
	}
	if channelString != "" {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Blocked channels",
			Value: channelString,
		})
	}

	_, err = ctx.Send("", e)
	return
}

func (bot *Bot) modHlUnblock(ctx *bcr.Context) (err error) {
	conf, err := bot.hlConfig(ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	ch, err := ctx.ParseChannel(ctx.RawArgs)
	if err != nil || ch.GuildID != ctx.Guild.ID || (ch.Type != discord.GuildText && ch.Type != discord.GuildNews && ch.Type != discord.GuildCategory) {
		_, err = ctx.Replyc(bcr.ColourRed, "Channel not found.")
		return
	}

	isBlocked := false
	for _, id := range conf.Blocked {
		if discord.ChannelID(id) == ch.ID {
			isBlocked = true
			break
		}
	}

	if !isBlocked {
		_, err = ctx.Replyc(bcr.ColourRed, "That channel isn't blocked from highlights.")
		return
	}

	if len(conf.Blocked) == 1 {
		conf.Blocked = []uint64{}
	} else {
		for i := range conf.Blocked {
			if conf.Blocked[i] == uint64(ch.ID) {
				if i == 0 {
					conf.Blocked = conf.Blocked[1:]
				} else if i == len(conf.Blocked)-1 {
					conf.Blocked = conf.Blocked[:len(conf.Blocked)-1]
				} else {
					conf.Blocked = append(conf.Blocked[:i], conf.Blocked[i+1:]...)
				}
				break
			}
		}
	}

	err = bot.setConfig(conf)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Replyc(bcr.ColourGreen, "Unblocked %v from highlights!", ch.Mention())
	return
}
