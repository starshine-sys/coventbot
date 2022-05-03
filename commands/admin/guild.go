package admin

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/common"
)

func (bot *Bot) serverInfo(ctx *bcr.Context) (err error) {
	sf, err := discord.ParseSnowflake(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("Couldn't parse your input as a snowflake.")
		return
	}

	g, err := ctx.State.GuildWithCount(discord.GuildID(sf))
	if err != nil {
		_, err = ctx.Send("I'm not in that server, so I can't show info for it.")
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

	channels, err := ctx.State.Channels(discord.GuildID(sf))
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

	e := discord.Embed{
		Title: fmt.Sprintf("Info for %v", g.Name),
		Thumbnail: &discord.EmbedThumbnail{
			URL: g.IconURL(),
		},

		Fields: []discord.EmbedField{
			{
				Name:   "Owner",
				Value:  owner.Username + "#" + owner.Discriminator,
				Inline: true,
			},
			{
				Name:   "Members",
				Value:  fmt.Sprint(g.ApproximateMembers),
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
				Name:   "Features",
				Value:  strings.Join(guildFeaturesToString(g.Features), "\n"),
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
		case discord.InviteSplash:
			s = append(s, "Invite Splash")
		case discord.VIPRegions:
			s = append(s, "384kbps Voice")
		case discord.VanityURL:
			s = append(s, "Vanity URL")
		case discord.Verified:
			s = append(s, "Verified")
		case discord.Partnered:
			s = append(s, "Partnered")
		case discord.Public:
			s = append(s, "Public")
		case discord.Commerce:
			s = append(s, "Commerce")
		case discord.News:
			s = append(s, "News Channels")
		case discord.Discoverable:
			s = append(s, "Discoverable")
		case discord.Featurable:
			s = append(s, "Featurable")
		case discord.AnimatedIcon:
			s = append(s, "Animated Icon")
		case discord.Banner:
			s = append(s, "Banner")
		}
	}
	if len(s) == 0 {
		s = append(s, "None")
	}
	return s
}
