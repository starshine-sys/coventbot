package info

import (
	"fmt"
	"regexp"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

var inviteRegexp = regexp.MustCompile(`https:\/\/discord.gg\/(\w+)`)

func (bot *Bot) inviteInfo(ctx *bcr.Context) (err error) {
	code := ctx.RawArgs
	// if the argument is a link, use that instead
	groups := inviteRegexp.FindStringSubmatch(ctx.RawArgs)
	if len(groups) == 2 {
		code = groups[1]
	}

	g, err := ctx.State.InviteWithCounts(code)
	if err != nil {
		_, err = ctx.Send("You did not give a valid invite.")
		return
	}

	e := discord.Embed{
		Color: ctx.Router.EmbedColor,

		Title:       fmt.Sprintf("Invite for %v", g.Guild.Name),
		Description: fmt.Sprintf("This invite points to the #%v channel (%v).", g.Channel.Name, g.Channel.Mention()),

		Thumbnail: &discord.EmbedThumbnail{
			URL: g.Guild.IconURL(),
		},

		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("Server ID: %v", g.Guild.ID),
		},
	}

	if g.Inviter != nil {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:   "Created by",
			Value:  fmt.Sprintf("%v#%v\n%v\nID: %v", g.Inviter.Username, g.Inviter.Discriminator, g.Inviter.Mention(), g.Inviter.ID),
			Inline: true,
		})
	}

	e.Fields = append(e.Fields, discord.EmbedField{
		Name:   "Members",
		Value:  fmt.Sprintf("👥 %v\n<:online2:826545116838756412> %v", g.ApproximateMembers, g.ApproximatePresences),
		Inline: true,
	})

	_, err = ctx.Send("", e)
	return
}
