package quotes

import (
	"fmt"
	"sort"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) list(ctx *bcr.Context) (err error) {
	if !bot.quotesEnabled(ctx.Guild.ID) {
		_, err = ctx.Send("Quotes aren't enabled on this server, sorry :(\nAsk a server admin to enable it!")
		return
	}

	quotes, err := bot.quotes(ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	title := "Quotes"

	userString, _ := ctx.Flags.GetString("user")
	channelString, _ := ctx.Flags.GetString("channel")
	server, _ := ctx.Flags.GetBool("server")

	if server {
		orig := quotes
		quotes = nil

		for _, q := range orig {
			if q.ServerID == ctx.Message.GuildID {
				quotes = append(quotes, q)
			}
		}
	}

	if userString != "" {
		var u *discord.User
		m, err := ctx.ParseMember(userString)
		if err == nil {
			u = &m.User
		} else {
			u, err = ctx.ParseUser(userString)
			if err != nil {
				_, err = ctx.Sendf("No user named ``%v`` found.", bcr.EscapeBackticks(userString))
				return err
			}
		}

		orig := quotes
		quotes = nil

		for _, q := range orig {
			if q.UserID == u.ID {
				quotes = append(quotes, q)
			}
		}

		title += " by " + u.Username
	}

	if channelString != "" {
		ch, err := ctx.ParseChannel(channelString)
		if err != nil {
			_, err = ctx.Sendf("No channel named ``%v`` found.", bcr.EscapeBackticks(channelString))
			return err
		}

		if ch.GuildID != ctx.Channel.GuildID || ch.Type != discord.GuildText {
			_, err = ctx.Sendf("No channel named ``%v`` found.", bcr.EscapeBackticks(channelString))
			return err
		}

		orig := quotes
		quotes = nil

		for _, q := range orig {
			if q.ChannelID == ch.ID {
				quotes = append(quotes, q)
			}
		}

		title += " in #" + ch.Name
	}

	if len(quotes) == 0 {
		_, err = ctx.Sendf("Couldn't find any quotes matching your criteria :(")
		return
	}

	// sorting
	sortByID, _ := ctx.Flags.GetBool("sort-by-message")
	rev, _ := ctx.Flags.GetBool("reversed")

	if sortByID {
		sort.Slice(quotes, func(i, j int) bool {
			return quotes[i].MessageID < quotes[j].MessageID
		})
	}

	if rev {
		for i, j := 0, len(quotes)-1; i < j; i, j = i+1, j-1 {
			quotes[i], quotes[j] = quotes[j], quotes[i]
		}
	}

	s := []string{}

	for _, q := range quotes {
		s = append(s, fmt.Sprintf("`%v` [(jump)](https://discord.com/channels/%v/%v/%v) by <@!%v>\n", q.HID, q.ServerID, q.ChannelID, q.MessageID, q.UserID))
	}

	_, err = bot.PagedEmbed(ctx,
		bcr.StringPaginator(fmt.Sprintf("%v (%v)", title, len(quotes)), bcr.ColourBlurple, s, 15), 10*time.Minute,
	)
	return err
}
