package levels

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/pkg/errors"
	"github.com/starshine-sys/bcr"
)

type importedLevel struct {
	UserID discord.UserID `json:"user_id"`
	XP     int64          `json:"exp"`
	Colour discord.Color  `json:"-"`
}

func (i *importedLevel) UnmarshalJSON(src []byte) (err error) {
	raw := struct {
		UserID    discord.UserID `json:"user_id"`
		XP        int64          `json:"exp"`
		CardColor string         `json:"card_color"`
	}{}

	err = json.Unmarshal(src, &raw)
	if err != nil {
		return err
	}
	i.UserID = raw.UserID
	i.XP = raw.XP

	if raw.CardColor != "" {
		clr, err := strconv.ParseUint(raw.CardColor, 16, 64)
		if err == nil {
			i.Colour = discord.Color(clr)
		}
	}
	return nil
}

func (bot *Bot) importLevels(ctx *bcr.Context) (err error) {
	if len(ctx.Message.Attachments) < 1 && len(ctx.Args) < 1 {
		return ctx.SendX("Please attach a JSON export, or provide a link to one.")
	}
	url := ctx.RawArgs
	if len(ctx.Message.Attachments) > 0 && strings.HasSuffix(ctx.Message.Attachments[0].URL, ".json") {
		url = ctx.Message.Attachments[0].URL
	}

	resp, err := http.Get(url)
	if err != nil {
		bot.Sugar.Errorf("downloading export file: %v", err)
		return bot.Report(ctx, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return ctx.SendfX("Encountered an error code while downloading the export file: %v\nPlease try again.", resp.StatusCode)
	}

	levels := make([]importedLevel, 0)
	err = json.NewDecoder(resp.Body).Decode(&levels)
	if err != nil {
		bot.Sugar.Errorf("decoding export file: %v", err)
		return ctx.SendfX("There was an error decoding the export file: %v", bcr.AsCode(err.Error()))
	}

	yes, timeout := ctx.ConfirmButton(ctx.Author.ID, bcr.ConfirmData{
		Message:   "Are you sure you want to import level data? This will overwrite **all** existing data!",
		YesPrompt: "Yes, overwrite data",
		NoPrompt:  "Cancel",
	})
	if timeout {
		return ctx.SendX("Prompt timed out.")
	}
	if !yes {
		return ctx.SendX("Cancelled.")
	}

	err = bot.importLevelsDB(ctx.Message.GuildID, levels)
	if err != nil {
		return bot.Report(ctx, err)
	}

	return ctx.SendfX("Success! Imported %v user levels.", len(levels))
}

func (bot *Bot) importLevelsDB(guildID discord.GuildID, levels []importedLevel) error {
	ctx := context.Background()

	tx, err := bot.DB.Pool.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "beginning transaction")
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "delete from levels where server_id = $1", guildID)
	if err != nil {
		return errors.Wrap(err, "deleting existing rows")
	}

	for _, lvl := range levels {
		_, err = tx.Exec(ctx, `insert into levels (server_id, user_id, xp, colour)
		values ($1, $2, $3, $4) on conflict (server_id, user_id) do update
		set xp = $3, colour = $4`, guildID, lvl.UserID, lvl.XP, lvl.Colour)
		if err != nil {
			return errors.Wrap(err, "inserting level row")
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "committing transaction")
	}
	return nil
}
