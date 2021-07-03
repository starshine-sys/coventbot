package static

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) exportEmotes(ctx *bcr.Context) (err error) {
	emojis, err := ctx.State.Emojis(ctx.Guild.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	w := new(bytes.Buffer)

	zip := zip.NewWriter(w)

	msg, err := ctx.Reply("**Downloading emotes, please wait...**\nDownloaded 0/%v", len(emojis))
	if err != nil {
		return err
	}

	for i, e := range emojis {
		if i%10 == 0 {
			ctx.Edit(msg, "", true, discord.Embed{
				Color:       ctx.Router.EmbedColor,
				Description: fmt.Sprintf("**Downloading emotes, please wait...**\nDownloaded %v/%v", i, len(emojis)),
			})
		}

		fn := fmt.Sprintf("static/%v.%v.png", e.Name, e.ID)
		if e.Animated {
			fn = fmt.Sprintf("animated/%v.%v.gif", e.Name, e.ID)
		}

		w, err := zip.Create(fn)
		if err != nil {
			return bot.Report(ctx, err)
		}

		resp, err := http.Get(e.EmojiURL())
		if err != nil {
			return bot.Report(ctx, err)
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return bot.Report(ctx, err)
		}

		w.Write(b)
	}

	err = zip.Close()
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.NewMessage().Content(fmt.Sprintf("Here you go! Downloaded %v emotes.", len(emojis))).AddFile("emotes.zip", w).Send()
	return
}
