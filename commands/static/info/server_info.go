// SPDX-License-Identifier: AGPL-3.0-only
package info

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/common"
)

func (bot *Bot) serverInfo(ctx *bcr.Context) (err error) {
	err = ctx.State.Typing(ctx.Channel.ID)
	if err != nil {
		bot.Sugar.Errorf("Error triggering typing: %v", err)
	}

	g, err := ctx.State.GuildWithCount(ctx.Message.GuildID)
	if err != nil {
		return
	}

	owner, err := ctx.State.User(g.OwnerID)
	if err != nil {
		return
	}

	var animatedEmoji, staticEmoji int
	for _, e := range g.Emojis {
		if e.Animated {
			animatedEmoji++
		} else {
			staticEmoji++
		}
	}

	channels, err := ctx.State.Channels(ctx.Message.GuildID)
	if err != nil {
		return
	}

	var textCount, textLocked, voiceCount, voiceLocked, total, categories int

	for _, ch := range channels {
		switch ch.Type {
		case discord.GuildText, discord.GuildNews:
			total++
			textCount++
			for _, p := range ch.Overwrites {
				if p.ID == discord.Snowflake(ch.GuildID) && p.Deny.Has(discord.PermissionViewChannel) {
					textLocked++
				}
			}
		case discord.GuildVoice:
			total++
			voiceCount++
			for _, p := range ch.Overwrites {
				if p.ID == discord.Snowflake(ch.GuildID) && p.Deny.Has(discord.PermissionViewChannel) {
					voiceLocked++
				}
			}
		case discord.GuildCategory:
			categories++
		}
	}

	var humans, bots int64

	for _, m := range bot.Members(ctx.Guild.ID) {
		if m.User.Bot {
			bots++
		} else {
			humans++
		}
	}

	e := discord.Embed{
		Title: fmt.Sprintf("Info for %v", g.Name),
		Thumbnail: &discord.EmbedThumbnail{
			URL: g.IconURL(),
		},

		Fields: []discord.EmbedField{
			{
				Name:   "Owner",
				Value:  owner.Tag(),
				Inline: true,
			},
			{
				Name:   "Members",
				Value:  fmt.Sprintf("Total: %v\nHumans: %v\nBots: %v", humanize.Comma(int64(g.ApproximateMembers)), humanize.Comma(humans), humanize.Comma(bots)),
				Inline: true,
			},
			{
				Name:   "Level",
				Value:  fmt.Sprintf("%v (%v boosts)", g.NitroBoost, g.NitroBoosters),
				Inline: true,
			},
			{
				Name:   "Emoji",
				Value:  fmt.Sprintf("%v total\n%v static\n%v animated", len(g.Emojis), staticEmoji, animatedEmoji),
				Inline: true,
			},
			{
				Name:   "Roles",
				Value:  fmt.Sprint(len(g.Roles)),
				Inline: true,
			},
			{
				Name: "Channels",
				Value: fmt.Sprintf(`%v (in %v categories)
<:textchannel:770274583223336990> %v (%v locked)
<:voicechannel:770274509012992020> %v (%v locked)`, total, categories, textCount, textLocked, voiceCount, voiceLocked),
				Inline: true,
			},
			{
				Name:   "Created",
				Value:  fmt.Sprintf("<t:%v>\n(%v)", g.ID.Time().Unix(), common.FormatTime(g.ID.Time())),
				Inline: true,
			},
			{
				Name:   "Features",
				Value:  strings.Join(guildFeaturesToString(g.Features), ", "),
				Inline: false,
			},
		},

		Color:     bcr.ColourBlurple,
		Timestamp: discord.NewTimestamp(g.ID.Time()),
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v | Created", g.ID),
		},
	}

	_, err = ctx.Send("", e)
	return
}

// guildFeaturesToString converts a []discord.GuildFeature to a []string, and prettifies the names
func guildFeaturesToString(g []discord.GuildFeature) (s []string) {
	for _, f := range g {
		switch f {
		case discord.VIPRegions:
			s = append(s, "384kbps Voice")
		case discord.VanityURL:
			s = append(s, "Vanity URL")
		case discord.News:
			s = append(s, "News Channels")
		default:
			s = append(s,
				strings.Title(strings.ToLower(
					strings.ReplaceAll(
						string(f), "_", " "))))
		}
	}
	if len(s) == 0 {
		s = append(s, "None")
	}
	return s
}
