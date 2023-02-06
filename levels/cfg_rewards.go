// SPDX-License-Identifier: AGPL-3.0-only
package levels

import (
	"context"
	"strconv"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) cmdAddReward(ctx *bcr.Context) (err error) {
	lvl, err := strconv.ParseInt(ctx.Args[0], 0, 0)
	if err != nil {
		_, err = ctx.Send("Couldn't parse your input as an integer.")
		return
	}
	role, err := ctx.ParseRole(ctx.Args[1])
	if err != nil {
		_, err = ctx.Send("Couldn't find a role with that name or ID.")
		return
	}

	err = bot.addReward(ctx.Message.GuildID, lvl, role.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Set the reward for level `%v` to `%v`.", lvl, role.Name)
	return
}

func (bot *Bot) cmdDelReward(ctx *bcr.Context) (err error) {
	lvl, err := strconv.ParseInt(ctx.Args[0], 0, 0)
	if err != nil {
		_, err = ctx.Send("Couldn't parse your input as an integer.")
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "delete from level_rewards where server_id = $1 and lvl = $2", ctx.Message.GuildID, lvl)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Removed the reward for level `%v`.", lvl)
	return
}
