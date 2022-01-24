package levels

import (
	"context"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) setlvl(ctx *bcr.Context) (err error) {
	u, err := ctx.ParseMember(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("Couldn't find a member with that name.")
		return
	}

	sc, err := bot.getGuildConfig(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	lvl, err := strconv.ParseInt(ctx.Args[1], 0, 0)
	if err != nil || lvl < 0 || lvl > 1000 {
		_, err = ctx.Send("Couldn't parse your input as a valid level.")
		return
	}

	var xp int64
	if lvl == 0 {
		xp = 0
	} else {
		xp = sc.CalculateExp(lvl)
	}

	_, err = bot.DB.Pool.Exec(context.Background(), `insert into levels
	(server_id, user_id, xp) values ($1, $2, $3)
	on conflict (server_id, user_id) do update
	set xp = $3`, ctx.Message.GuildID, u.User.ID, xp)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Updated %v's level to `%v`.", u.Mention(), humanize.Comma(sc.CalculateLevel(xp)))
	return
}
