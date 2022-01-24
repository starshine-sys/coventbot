package levels

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codeberg.org/eviedelta/detctime/durationparser"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) nolevelsList(ctx *bcr.Context) (err error) {
	if len(ctx.Args) > 0 {
		return bot.nolevelsAdd(ctx)
	}

	list, err := bot.guildNolevels(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if len(list) == 0 {
		_, err = ctx.Send("There are no noleveled users.")
		return
	}

	var s []string
	for i, l := range list {
		buf := fmt.Sprintf("%v. %v: ", i+1, l.UserID.Mention())
		if l.Expires {
			buf += fmt.Sprintf("expires <t:%v>\n", l.Expiry.Unix())
		} else {
			buf += "does not expire\n"
		}
		s = append(s, buf)
	}

	_, err = bot.PagedEmbed(ctx,
		bcr.StringPaginator("User blacklist", bcr.ColourBlurple, s, 10), 10*time.Minute)
	return
}

func (bot *Bot) nolevelsAdd(ctx *bcr.Context) (err error) {
	u, err := ctx.ParseUser(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("User not found.")
		return
	}

	expiry := time.Now().UTC()
	expires := false

	if len(ctx.Args) > 1 {
		dur, err := durationparser.Parse(ctx.Args[1])
		if err != nil {
			_, err = ctx.Send("Couldn't parse your input as a valid duration.")
			return err
		}
		expiry = expiry.Add(dur)
		expires = true
	}

	err = bot.setNolevels(ctx.Message.GuildID, u.ID, expires, expiry)
	if err != nil {
		return bot.Report(ctx, err)
	}

	text := fmt.Sprintf("Noleveled %v ", u.Mention())
	if expires {
		text += fmt.Sprintf("until <t:%v>.", expiry.Unix())
	} else {
		text += "indefinitely."
	}

	if sc, _ := bot.getGuildConfig(ctx.Message.GuildID); sc.NolevelsLog.IsValid() {
		_, err = ctx.State.SendEmbeds(sc.NolevelsLog, discord.Embed{
			Title:       "User noleveled",
			Description: text,
			Fields: []discord.EmbedField{{
				Name:  "Responsible moderator",
				Value: fmt.Sprintf("%v#%v (%v)", ctx.Author.Username, ctx.Author.Discriminator, ctx.Author.Mention()),
			}},
			Color: bcr.ColourBlurple,
		})
		if err != nil {
			bot.Sugar.Errorf("Error logging nolevels: %v", err)
		}
	}

	_, err = ctx.Reply(text)
	return
}

func (bot *Bot) nolevelsRemove(ctx *bcr.Context) (err error) {
	u, err := ctx.ParseUser(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("User not found.")
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "delete from nolevels where user_id = $1 and server_id = $2", u.ID, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	if sc, _ := bot.getGuildConfig(ctx.Message.GuildID); sc.NolevelsLog.IsValid() {
		_, err = ctx.State.SendEmbeds(sc.NolevelsLog, discord.Embed{
			Title:       "User unblacklisted",
			Description: "Removed nolevels from " + u.Mention(),
			Fields: []discord.EmbedField{{
				Name:  "Responsible moderator",
				Value: fmt.Sprintf("%v#%v (%v)", ctx.Author.Username, ctx.Author.Discriminator, ctx.Author.Mention()),
			}},
			Color: bcr.ColourBlurple,
		})
		if err != nil {
			bot.Sugar.Errorf("Error logging nolevels: %v", err)
		}
	}

	_, err = ctx.Reply("Unblacklisted " + u.Mention())
	return
}

func (bot *Bot) nolevelLoop() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	for {
		select {
		case <-sc:
			break
		default:
		}

		list := []Nolevels{}

		err := pgxscan.Select(context.Background(), bot.DB.Pool, &list, "select n.*, l.nolevels_log as log_channel from nolevels as n, server_levels as l where n.server_id = l.id and n.expires = true and n.expiry < $1 limit 5", time.Now().UTC())
		if err != nil {
			bot.Sugar.Errorf("Error getting nolevels: %v", err)
			time.Sleep(time.Second)
			continue
		}

		for _, n := range list {
			_, err = bot.DB.Pool.Exec(context.Background(), "delete from nolevels where server_id = $1 and user_id = $2", n.ServerID, n.UserID)
			if err != nil {
				bot.Sugar.Errorf("Error deleting nolevel entry: %v", err)
				continue
			}

			s, _ := bot.Router.StateFromGuildID(n.ServerID)

			if n.LogChannel.IsValid() {
				_, err = s.SendEmbeds(n.LogChannel, discord.Embed{
					Title:       "User nolevel expired",
					Description: fmt.Sprintf("The blacklist of %v expired.", n.UserID.Mention()),
					Color:       bcr.ColourBlurple,
				})
				if err != nil {
					bot.Sugar.Errorf("Error sending nolevels log: %v", err)
				}
			}
		}

		time.Sleep(time.Second)
	}
}
