// SPDX-License-Identifier: AGPL-3.0-only
package static

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/starshine-sys/bcr"
	bcr2 "github.com/starshine-sys/bcr/v2"
)

const cooldownTime = time.Minute

type cdKey struct {
	user discord.UserID
	ch   discord.ChannelID
}

var cd = &cooldown{m: map[cdKey]time.Time{}}

type cooldown struct {
	m  map[cdKey]time.Time
	mu sync.Mutex
}

// Get returns true if the user is on cooldown
func (c *cooldown) Get(userID discord.UserID, ch discord.ChannelID) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	t, ok := c.m[cdKey{userID, ch}]
	if !ok {
		return false
	}

	return t.After(time.Now())
}

func (c *cooldown) Set(userID discord.UserID, ch discord.ChannelID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[cdKey{userID, ch}] = time.Now().Add(cooldownTime)
}

func (bot *Bot) linktoSlash(ctx *bcr2.CommandContext) error {
	sf, _ := ctx.Options.Find("channel").SnowflakeValue()
	reason := ctx.Options.Find("reason").String()

	ch, ok := ctx.Data.Resolved.Channels[discord.ChannelID(sf)]
	if !ok {
		return ctx.ReplyEphemeral("Somehow, the channel you gave wasn't sent to me.")
	}

	if cd.Get(ctx.User.ID, ch.ID) {
		return ctx.ReplyEphemeral("You're currently on cooldown.")
	}

	perms := discord.CalcOverwrites(*ctx.Guild, ch, *ctx.Member)
	if !perms.Has(discord.PermissionViewChannel) || !perms.Has(discord.PermissionSendMessages) {
		return ctx.ReplyEphemeral("You can't send messages to " + ch.Mention() + "!")
	}

	desc := ""
	if reason != "" {
		desc = fmt.Sprintf("\n**Topic**\n> %v", reason)
	}

	err := ctx.Reply("", discord.Embed{
		Color:       bcr.ColourBlurple,
		Description: fmt.Sprintf("Directing conversation to %v%v", ch.Mention(), desc),
		Footer: &discord.EmbedFooter{
			Icon: ctx.User.AvatarURLWithType(discord.PNGImage),
			Text: ctx.User.Tag(),
		},
		Timestamp: discord.NowTimestamp(),
	})
	if err != nil {
		return err
	}

	msg1, err := ctx.Original()
	if err != nil {
		return err
	}

	msg2, err := ctx.State.SendMessage(ch.ID, fmt.Sprintf("<https://discord.com/channels/%v/%v/%v>", ctx.Guild.ID, ctx.Channel.ID, msg1.ID), discord.Embed{
		Color:       bcr.ColourBlurple,
		Description: fmt.Sprintf("Conversation moved from %v by %v%v", ctx.Channel.Mention(), ctx.User.Mention(), desc),
		Footer: &discord.EmbedFooter{
			Icon: ctx.User.AvatarURLWithType(discord.PNGImage),
			Text: ctx.User.Tag(),
		},
		Timestamp: discord.NowTimestamp(),
	})
	if err != nil {
		return err
	}

	_, err = ctx.State.EditInteractionResponse(discord.AppID(ctx.State.Ready().User.ID), ctx.InteractionToken, api.EditInteractionResponseData{
		Content: option.NewNullableString(
			fmt.Sprintf("<https://discord.com/channels/%v/%v/%v>", ch.GuildID, ch.ID, msg2.ID),
		),
	})
	if err != nil {
		return err
	}

	cd.Set(ctx.User.ID, ch.ID)
	return nil
}

func (bot *Bot) linkto(ctx *bcr.Context) (err error) {
	var ch *discord.Channel
	var reason string

	ch, err = ctx.ParseChannel(ctx.Args[0])
	if err != nil {
		return ctx.SendEphemeral("You didn't give a valid channel.")
	}
	reason = strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))
	if reason == ctx.RawArgs && len(ctx.Args) > 1 {
		reason = strings.Join(ctx.Args[1:], " ")
	}

	if cd.Get(ctx.User().ID, ch.ID) {
		return ctx.SendEphemeral("You're currently on cooldown.")
	}

	if ctx.GetGuild() == nil {
		return ctx.SendEphemeral("This command can only be used in servers.")
	}

	if ch.GuildID != ctx.GetGuild().ID {
		return ctx.SendEphemeral("You can only link to channels in the same server.")
	}

	if ch.Type != discord.GuildNews && ch.Type != discord.GuildText {
		return ctx.SendEphemeral("You can only link to text channels.")
	}

	member, err := bot.Member(ctx.GetGuild().ID, ctx.User().ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	perms := discord.CalcOverwrites(*ctx.GetGuild(), *ch, member)
	if !perms.Has(discord.PermissionViewChannel) || !perms.Has(discord.PermissionSendMessages) {
		return ctx.SendEphemeral("You can't send messages to " + ch.Mention() + "!")
	}

	desc := ""
	if reason != "" {
		desc = fmt.Sprintf("\n**Topic**\n> %v", reason)
	}

	msg1, err := ctx.Send("", discord.Embed{
		Color:       bcr.ColourBlurple,
		Description: fmt.Sprintf("Directing conversation to %v%v", ch.Mention(), desc),
		Footer: &discord.EmbedFooter{
			Icon: ctx.User().AvatarURLWithType(discord.PNGImage),
			Text: ctx.User().Tag(),
		},
		Timestamp: discord.NowTimestamp(),
	})
	if err != nil {
		return
	}

	msg2, err := ctx.Session().SendMessage(ch.ID, fmt.Sprintf("<https://discord.com/channels/%v/%v/%v>", ctx.GetGuild().ID, ctx.GetChannel().ID, msg1.ID), discord.Embed{
		Color:       bcr.ColourBlurple,
		Description: fmt.Sprintf("Conversation moved from %v by %v%v", ctx.GetChannel().Mention(), ctx.User().Mention(), desc),
		Footer: &discord.EmbedFooter{
			Icon: ctx.User().AvatarURLWithType(discord.PNGImage),
			Text: ctx.User().Tag(),
		},
		Timestamp: discord.NowTimestamp(),
	})
	if err != nil {
		return
	}

	_, err = ctx.EditOriginal(api.EditInteractionResponseData{
		Content: option.NewNullableString(
			fmt.Sprintf("<https://discord.com/channels/%v/%v/%v>", ch.GuildID, ch.ID, msg2.ID),
		),
	})
	if err != nil {
		return
	}

	cd.Set(ctx.User().ID, ch.ID)
	return
}
