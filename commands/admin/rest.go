package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/starshine-sys/bcr"
)

const restTimeout = 10 * time.Second

func (bot *Bot) rest(ctx *bcr.Context) (err error) {
	respond := func(warning bool, v ...interface{}) error {
		url := ctx.RawArgs
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

	rctx, cancel := context.WithTimeout(context.Background(), restTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(rctx, "GET", ctx.RawArgs, nil)
	if err != nil {
		return respond(true, err)
	}

	resp, err := bot.Client.Do(req)
	if err != nil {
		return respond(true, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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
