package names

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) usernames(ctx *bcr.Context) (err error) {
	u := ctx.Author
	if ctx.RawArgs != "" {
		user, err := ctx.ParseUser(ctx.RawArgs)
		if err == nil {
			u = *user
		}
	}

	usernames := make([]string, 0)
	err = bot.DB.Pool.QueryRow(context.Background(), "select array(select name from usernames where user_id = $1 order by time desc)", u.ID).Scan(&usernames)
	if err != nil {
		_, err = ctx.Send("Couldn't find any username history for that user.", nil)
		return
	}

	var e []discord.Embed
	var buf string
	var count int

	for i, n := range usernames {
		if count >= 10 {
			e = append(e, discord.Embed{
				Title:       fmt.Sprintf("Username history for %v#%v", u.Username, u.Discriminator),
				Description: buf,
				Color:       ctx.Router.EmbedColor,
				Footer: &discord.EmbedFooter{
					Text: fmt.Sprintf("User ID: %v", u.ID),
				},
			})

			buf = ""
			count = 0
		}

		buf += fmt.Sprintf("%v. %v\n", len(usernames)-i, n)
		count++
	}

	e = append(e, discord.Embed{
		Title:       fmt.Sprintf("Username history for %v#%v", u.Username, u.Discriminator),
		Description: buf,
		Color:       ctx.Router.EmbedColor,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("User ID: %v", u.ID),
		},
	})

	_, err = ctx.PagedEmbed(e, false)
	return
}

func (bot *Bot) nicknames(ctx *bcr.Context) (err error) {
	u := ctx.Author
	if ctx.RawArgs != "" {
		user, err := ctx.ParseUser(ctx.RawArgs)
		if err == nil {
			u = *user
		}
	}

	nicknames := make([]string, 0)
	err = bot.DB.Pool.QueryRow(context.Background(), "select array(select name from nicknames where user_id = $1 and server_id = $2 order by time desc)", u.ID, ctx.Message.GuildID).Scan(&nicknames)
	if err != nil {
		fmt.Println(err)

		_, err = ctx.Send("Couldn't find any nickname history for that user.", nil)
		return
	}

	var e []discord.Embed
	var buf string
	var count int

	for i, n := range nicknames {
		if count >= 10 {
			e = append(e, discord.Embed{
				Title:       fmt.Sprintf("Nickname history for %v#%v", u.Username, u.Discriminator),
				Description: buf,
				Color:       ctx.Router.EmbedColor,
				Footer: &discord.EmbedFooter{
					Text: fmt.Sprintf("User ID: %v", u.ID),
				},
			})

			buf = ""
			count = 0
		}

		buf += fmt.Sprintf("%v. %v\n", len(nicknames)-i, n)
		count++
	}

	e = append(e, discord.Embed{
		Title:       fmt.Sprintf("Nickname history for %v#%v", u.Username, u.Discriminator),
		Description: buf,
		Color:       ctx.Router.EmbedColor,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("User ID: %v", u.ID),
		},
	})

	_, err = ctx.PagedEmbed(e, false)
	return
}
