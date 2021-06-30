package quotes

import (
	"regexp"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/jackc/pgx/v4"
	"github.com/starshine-sys/bcr"
)

var idMatch = regexp.MustCompile(`(?i)^[a-zA-Z]{5}$`)

func (bot *Bot) quote(ctx *bcr.Context) (err error) {
	if !bot.quotesEnabled(ctx.Guild.ID) {
		_, err = ctx.Send("Quotes aren't enabled on this server, sorry :(\nAsk a server admin to enable it!", nil)
		return
	}

	if len(ctx.Args) == 0 {
		q, err := bot.serverQuote(ctx.Guild.ID)
		if err != nil {
			if errors.Cause(err) == pgx.ErrNoRows {
				_, err = ctx.Send("This server has no quotes yet. React with 💬 to a message to add it as a quote!", nil)
				return err
			}

			return bot.Report(ctx, err)
		}

		e := q.Embed(bot.PK)
		_, err = ctx.Send("", &e)
		return err
	}

	if idMatch.MatchString(ctx.RawArgs) {
		q, err := bot.getQuote(ctx.RawArgs, ctx.Guild.ID)
		if err != nil {
			_, err = ctx.Send("No quote with that ID found! Note that a quote ID is 5 characters long.", nil)
			return err
		}

		e := q.Embed(bot.PK)
		_, err = ctx.Send("", &e)
		return err
	}

	var u discord.User
	m, err := ctx.ParseMember(ctx.RawArgs)
	if err != nil {
		user, err := ctx.ParseUser(ctx.RawArgs)
		if err != nil {
			_, err = ctx.Send("No user with that name found.", nil)
		}
		u = *user
	} else {
		u = m.User
	}

	q, err := bot.userQuote(ctx.Guild.ID, u.ID)
	if err != nil {
		_, err = ctx.Send("That user doesn't have any quotes, sorry :(", nil)
		return err
	}

	e := q.Embed(bot.PK)
	_, err = ctx.Send("", &e)
	return
}
