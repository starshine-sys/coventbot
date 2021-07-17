package highlights

import (
	"fmt"
	"strings"

	"github.com/caneroj1/stemmer"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

// MaxHighlights ...
const (
	MaxHighlights = 50
	MaxHlLength   = 80
)

func (bot *Bot) addHl(ctx *bcr.Context) (err error) {
	hlConf, err := bot.hlConfig(ctx.Guild.ID)
	if err == nil {
		if !hlConf.HlEnabled {
			_, err = ctx.Replyc(bcr.ColourRed, "Highlights aren't enabled on this server :(")
			return
		}
	}

	hl := strings.ToLower(ctx.RawArgs)

	if len(ctx.Message.Mentions) > 0 {
		_, err = ctx.Replyc(bcr.ColourRed, "Can't have any mentions in a highlight.")
		return
	}

	hls, err := bot.userHighlights(ctx.Guild.ID, ctx.Author.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(hls.Highlights) >= MaxHighlights {
		_, err = ctx.Replyc(bcr.ColourRed, "You have reached the limit of %v highlights.", MaxHighlights)
		return
	}

	if len(hl) > MaxHlLength {
		_, err = ctx.Replyc(bcr.ColourRed, "Your highlight is too long (%v > %v characters)", len(hl), MaxHlLength)
		return
	}

	stems := stemmer.StemMultiple(hls.Highlights)
	stem := stemmer.Stem(hl)
	for i, s := range stems {
		if s == stem {
			_, err = ctx.Replyc(bcr.ColourRed, "That word is already highlighted for you (%v = %v)", hl, hls.Highlights[i])
			return
		}
	}

	hls.Highlights = append(hls.Highlights, hl)
	err = bot.setUserHighlights(hls)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("", discord.Embed{
		Author: &discord.EmbedAuthor{
			Icon: ctx.Author.AvatarURLWithType(discord.PNGImage),
			Name: ctx.DisplayName() + "'s highlights",
		},
		Color:       bcr.ColourGreen,
		Description: fmt.Sprintf("Added ``%v`` to your highlights.", bcr.EscapeBackticks(hl)),
	})
	return
}

func (bot *Bot) delHl(ctx *bcr.Context) (err error) {
	hlConf, err := bot.hlConfig(ctx.Guild.ID)
	if err == nil {
		if !hlConf.HlEnabled {
			_, err = ctx.Replyc(bcr.ColourRed, "Highlights aren't enabled on this server :(")
			return
		}
	}

	hl := strings.ToLower(ctx.RawArgs)

	hls, err := bot.userHighlights(ctx.Guild.ID, ctx.Author.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	found := false
	for _, word := range hls.Highlights {
		if word == hl {
			found = true
			break
		}
	}

	if !found {
		_, err = ctx.Replyc(bcr.ColourRed, "That word isn't highlighted for you.")
		return
	}

	if len(hls.Highlights) == 1 {
		hls.Highlights = []string{}
	} else {
		for i := range hls.Highlights {
			if hls.Highlights[i] == hl {
				if i == 0 {
					hls.Highlights = hls.Highlights[1:]
				} else if i == len(hls.Highlights)-1 {
					hls.Highlights = hls.Highlights[:len(hls.Highlights)-1]
				} else {
					hls.Highlights = append(hls.Highlights[:i], hls.Highlights[i+1:]...)
				}
				break
			}
		}
	}

	err = bot.setUserHighlights(hls)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Send("", discord.Embed{
		Author: &discord.EmbedAuthor{
			Icon: ctx.Author.AvatarURLWithType(discord.PNGImage),
			Name: ctx.DisplayName() + "'s highlights",
		},
		Color:       bcr.ColourGreen,
		Description: fmt.Sprintf("Removed ``%v`` from your highlights.", bcr.EscapeBackticks(hl)),
	})
	return
}

func (bot *Bot) listHl(ctx *bcr.Context) (err error) {
	hlConf, err := bot.hlConfig(ctx.Guild.ID)
	if err == nil {
		if !hlConf.HlEnabled {
			_, err = ctx.Replyc(bcr.ColourRed, "Highlights aren't enabled on this server :(")
			return
		}
	}

	user := &ctx.Author
	if len(ctx.Args) > 0 {
		if check, _ := bot.HelperRole.Check(ctx); !check {
			_, err = ctx.Send("Only moderators can see other users' highlights.")
			return
		}

		user, err = ctx.ParseUser(ctx.RawArgs)
		if err != nil {
			_, err = ctx.Send("User not found.")
			return
		}
	}

	hl, err := bot.userHighlights(ctx.Guild.ID, user.ID)
	if err != nil {
		s := "That user has no highlights."
		if len(ctx.Args) == 0 {
			s = "You have no highlights."
		}
		_, err = ctx.Send(s)
		return
	}

	if len(hl.Highlights) == 0 {
		hl.Highlights = []string{"None"}
	}

	e := discord.Embed{
		Author: &discord.EmbedAuthor{
			Icon: user.AvatarURLWithType(discord.PNGImage),
			Name: user.Username + "'s highlights",
		},
		Color:       ctx.Router.EmbedColor,
		Description: strings.Join(hl.Highlights, "\n"),
	}

	if len(hl.Blocked) != 0 {
		channels, err := ctx.State.Channels(ctx.Guild.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		blockedUsers := ""
		blockedChannels := ""

		for _, id := range hl.Blocked {
			isCh := false
			for _, ch := range channels {
				if ch.ID == discord.ChannelID(id) {
					blockedChannels += fmt.Sprintf("<#%v>\n", id)
					isCh = true
					break
				}

				if !isCh {
					blockedUsers += fmt.Sprintf("<@!%v>\n", id)
				}
			}
		}

		if blockedChannels != "" {
			e.Fields = append(e.Fields, discord.EmbedField{
				Name:  "Blocked channel(s)",
				Value: blockedChannels,
			})
		}

		if blockedUsers != "" {
			e.Fields = append(e.Fields, discord.EmbedField{
				Name:  "Blocked user(s)",
				Value: blockedUsers,
			})
		}
	}

	_, err = ctx.Send("", e)
	return
}

func (bot *Bot) hlBlock(ctx *bcr.Context) (err error) {
	hlConf, err := bot.hlConfig(ctx.Guild.ID)
	if err == nil {
		if !hlConf.HlEnabled {
			_, err = ctx.Replyc(bcr.ColourRed, "Highlights aren't enabled on this server :(")
			return
		}
	}

	var isUser bool
	var id uint64

	u, err := ctx.ParseUser(ctx.RawArgs)
	if err == nil {
		isUser = true
		id = uint64(u.ID)
	} else {
		ch, err := ctx.ParseChannel(ctx.RawArgs)
		if err != nil || ch.GuildID != ctx.Guild.ID || (ch.Type != discord.GuildText && ch.Type != discord.GuildNews && ch.Type != discord.GuildCategory) {
			_, err = ctx.Replyc(bcr.ColourRed, "Couldn't parse your input as a user or channel.")
			return err
		}

		id = uint64(ch.ID)
	}

	hls, err := bot.userHighlights(ctx.Guild.ID, ctx.Author.ID)
	if err != nil {
		hls = Highlight{
			UserID:   ctx.Author.ID,
			ServerID: ctx.Message.GuildID,
		}
	}

	for _, blocked := range hls.Blocked {
		if id == blocked {
			if isUser {
				_, err = ctx.Replyc(bcr.ColourRed, "You've already blocked that user.")
			} else {
				_, err = ctx.Replyc(bcr.ColourRed, "You've already blocked that channel.")
			}
			return
		}
	}

	hls.Blocked = append(hls.Blocked, id)
	err = bot.setUserHighlights(hls)
	if err != nil {
		return bot.Report(ctx, err)
	}

	e := discord.Embed{
		Author: &discord.EmbedAuthor{
			Icon: ctx.Author.AvatarURLWithType(discord.PNGImage),
			Name: ctx.DisplayName() + "'s highlights",
		},
		Color: bcr.ColourGreen,
	}

	if isUser {
		e.Description = fmt.Sprintf("Blocked <@!%v> from your highlights!", id)
	} else {
		e.Description = fmt.Sprintf("Blocked <#%v> from your highlights!", id)
	}

	_, err = ctx.Send("", e)
	return
}

func (bot *Bot) hlUnblock(ctx *bcr.Context) (err error) {
	hlConf, err := bot.hlConfig(ctx.Guild.ID)
	if err == nil {
		if !hlConf.HlEnabled {
			_, err = ctx.Replyc(bcr.ColourRed, "Highlights aren't enabled on this server :(")
			return
		}
	}

	var isUser bool
	var id uint64

	u, err := ctx.ParseUser(ctx.RawArgs)
	if err == nil {
		isUser = true
		id = uint64(u.ID)
	} else {
		ch, err := ctx.ParseChannel(ctx.RawArgs)
		if err != nil || ch.GuildID != ctx.Guild.ID || (ch.Type != discord.GuildText && ch.Type != discord.GuildNews && ch.Type != discord.GuildCategory) {
			_, err = ctx.Replyc(bcr.ColourRed, "Couldn't parse your input as a user or channel.")
			return err
		}

		id = uint64(ch.ID)
	}

	hls, err := bot.userHighlights(ctx.Guild.ID, ctx.Author.ID)
	if err != nil || len(hls.Blocked) == 0 {
		_, err = ctx.Replyc(bcr.ColourRed, "You have no blocked channels or users.")
	}

	isBlocked := false
	for _, blocked := range hls.Blocked {
		if id == blocked {
			isBlocked = true
			break
		}
	}

	if !isBlocked {
		if isUser {
			_, err = ctx.Replyc(bcr.ColourRed, "That user isn't blocked from your highlights.")
		} else {
			_, err = ctx.Replyc(bcr.ColourRed, "That channel isn't blocked from your highlights.")
		}
		return
	}

	if len(hls.Blocked) == 1 {
		hls.Blocked = []uint64{}
	} else {
		for i := range hls.Blocked {
			if hls.Blocked[i] == id {
				if i == 0 {
					hls.Blocked = hls.Blocked[1:]
				} else if i == len(hls.Blocked)-1 {
					hls.Blocked = hls.Blocked[:len(hls.Blocked)-1]
				} else {
					hls.Blocked = append(hls.Blocked[:i], hls.Blocked[i+1:]...)
				}
				break
			}
		}
	}

	err = bot.setUserHighlights(hls)
	if err != nil {
		return bot.Report(ctx, err)
	}

	e := discord.Embed{
		Author: &discord.EmbedAuthor{
			Icon: ctx.Author.AvatarURLWithType(discord.PNGImage),
			Name: ctx.DisplayName() + "'s highlights",
		},
		Color: bcr.ColourGreen,
	}

	if isUser {
		e.Description = fmt.Sprintf("Unblocked <@!%v> from your highlights!", id)
	} else {
		e.Description = fmt.Sprintf("Unblocked <#%v> from your highlights!", id)
	}

	_, err = ctx.Send("", e)
	return
}

func (bot *Bot) hlTest(ctx *bcr.Context) (err error) {
	hlConf, err := bot.hlConfig(ctx.Guild.ID)
	if err == nil {
		if !hlConf.HlEnabled {
			_, err = ctx.Replyc(bcr.ColourRed, "Highlights aren't enabled on this server :(")
			return
		}
	}

	e := discord.Embed{
		Author: &discord.EmbedAuthor{
			Icon: ctx.Author.AvatarURLWithType(discord.PNGImage),
			Name: ctx.DisplayName() + "'s highlights",
		},
		Color: bcr.ColourBlurple,
	}

	hls, err := bot.userHighlights(ctx.Guild.ID, ctx.Author.ID)
	if err != nil || len(hls.Highlights) == 0 {
		_, err = ctx.Replyc(bcr.ColourRed, "You have no highlights.")
		return
	}

	msg := strings.Fields(ctx.RawArgs)
	stemmer.StemConcurrent(&msg)

	for i, stem := range stemmer.StemMultiple(hls.Highlights) {
		matched := false
		for _, word := range msg {
			if stem == word {
				matched = true
				break
			}
		}

		if matched {
			e.Description += fmt.Sprintf("✅ %v\n", hls.Highlights[i])
		} else {
			e.Description += fmt.Sprintf("❌ %v\n", hls.Highlights[i])
		}
	}

	e.Description = strings.TrimSpace(e.Description)

	_, err = ctx.Send("", e)
	return
}
