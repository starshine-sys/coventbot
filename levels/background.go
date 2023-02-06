// SPDX-License-Identifier: AGPL-3.0-only
package levels

import (
	"context"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) background(ctx *bcr.Context) (err error) {
	sc, err := bot.getGuildConfig(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}
	if !sc.LevelsEnabled {
		_, err = ctx.Send("Levels are disabled on this server.")
		return
	}

	if len(ctx.Message.Attachments) == 0 {
		if ctx.RawArgs == "clear" || ctx.RawArgs == "--clear" || ctx.RawArgs == "-clear" {
			_, err = bot.DB.Pool.Exec(context.Background(), "update levels set background = $1 where server_id = $2 and user_id = $3", "", ctx.Message.GuildID, ctx.Author.ID)
			if err != nil {
				return bot.Report(ctx, err)
			}

			_, err = ctx.Reply("Level background cleared!")
			return
		}

		uc, err := bot.getUser(ctx.Message.GuildID, ctx.Author.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		if uc.Background == "" {
			_, err = ctx.Reply("You don't have a background set. Set one with `%vlvl background` and attaching an image.", ctx.Prefix)
			return err
		}

		e := discord.Embed{
			Description: fmt.Sprintf("To clear your level background, use `%vlvl background clear`", ctx.Prefix),
			Image: &discord.EmbedImage{
				URL: uc.Background,
			},
			Color: bcr.ColourBlurple,
		}

		_, err = ctx.Send("", e)
		return err
	}

	if !hasAnySuffix(ctx.Message.Attachments[0].Filename, ".png", ".jpeg", ".jpg") {
		_, err = ctx.Replyc(bcr.ColourRed, "You didn't give a valid image type (PNG or JPG).")
		return
	}

	if ctx.Message.Attachments[0].Size > 1024*1024 {
		_, err = ctx.Replyc(bcr.ColourRed, "The background image can't be bigger than 1 MB.")
		return
	}

	url := ctx.Message.Attachments[0].URL

	_, err = bot.DB.Pool.Exec(context.Background(), "update levels set background = $1 where server_id = $2 and user_id = $3", url, ctx.Message.GuildID, ctx.Author.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Level background changed to the attached image!\nNote that if the trigger message is deleted, your background will need to be reset.")
	return
}

func (bot *Bot) serverBackground(ctx *bcr.Context) (err error) {
	sc, err := bot.getGuildConfig(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}
	if !sc.LevelsEnabled {
		_, err = ctx.Send("Levels are disabled on this server.")
		return
	}

	if len(ctx.Message.Attachments) == 0 {
		if ctx.RawArgs == "clear" || ctx.RawArgs == "--clear" || ctx.RawArgs == "-clear" {
			_, err = bot.DB.Pool.Exec(context.Background(), "update server_levels set background = $1 where id = $2", "", ctx.Message.GuildID)
			if err != nil {
				return bot.Report(ctx, err)
			}

			_, err = ctx.Reply("Server level background cleared!")
			return
		}

		uc, err := bot.getUser(ctx.Message.GuildID, ctx.Author.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		if uc.Background == "" {
			_, err = ctx.Reply("This server doesn't have a default background set. Set one with `%vlvl background server` and attaching an image.", ctx.Prefix)
			return err
		}

		e := discord.Embed{
			Description: fmt.Sprintf("To clear this server's default level background, use `%vlvl background server clear`", ctx.Prefix),
			Image: &discord.EmbedImage{
				URL: uc.Background,
			},
			Color: bcr.ColourBlurple,
		}

		_, err = ctx.Send("", e)
		return err
	}

	if !hasAnySuffix(ctx.Message.Attachments[0].Filename, ".png", ".jpeg", ".jpg") {
		_, err = ctx.Replyc(bcr.ColourRed, "You didn't give a valid image type (PNG or JPG).")
		return
	}

	if ctx.Message.Attachments[0].Size > 1024*1024 {
		_, err = ctx.Replyc(bcr.ColourRed, "The background image can't be bigger than 1 MB.")
		return
	}

	url := ctx.Message.Attachments[0].URL

	_, err = bot.DB.Pool.Exec(context.Background(), "update server_levels set background = $1 where id = $2", url, ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Default level background changed to the attached image!\nNote that if the trigger message is deleted, the default background will need to be reset.")
	return
}

func hasAnySuffix(s string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}

	return false
}
