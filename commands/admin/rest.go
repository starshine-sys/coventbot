// SPDX-License-Identifier: AGPL-3.0-only
package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) rest(ctx *bcr.Context) (err error) {
	url := ctx.RawArgs
	if strings.HasPrefix(ctx.RawArgs, "/") {
		url = api.Endpoint + strings.TrimPrefix(ctx.RawArgs, "/")
	}

	bot.Sugar.Warnf("User %v (%v) is making GET request to %v", ctx.Author.Tag(), ctx.Author.ID, url)

	respond := func(warning bool, v ...interface{}) error {
		url := url
		if len(url) > 128 {
			url = url[:128] + "..."
		}

		title := "✅ Request to `" + url + "`"
		colour := discord.Color(bcr.ColourBlurple)
		if warning {
			title = "⚠️ Request to `" + url + "`"
			colour = bcr.ColourOrange
		}

		desc := fmt.Sprintln(v...)
		if len(desc) > 2000 {
			return ctx.SendFiles("", sendpart.File{
				Name:   "output.json",
				Reader: strings.NewReader(desc),
			})
		}

		e := discord.Embed{
			Title:       title,
			Color:       colour,
			Description: "```json\n" + desc + "```",
			Timestamp:   discord.NowTimestamp(),
		}

		_, err := ctx.Send("", e)
		return err
	}

	resp, err := ctx.State.Client.Request("GET", url)
	if err != nil {
		return respond(true, err)
	}
	defer resp.GetBody().Close()

	body, err := io.ReadAll(resp.GetBody())
	if err != nil {
		return respond(true, err)
	}

	var buf bytes.Buffer
	err = json.Indent(&buf, body, "", "    ")
	if err != nil {
		return respond(false, string(body))
	}

	return respond(false, buf.String())
}
