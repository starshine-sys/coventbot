package static

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/disintegration/imaging"
	"github.com/dustin/go-humanize"
	"github.com/fogleman/gg"
	"github.com/starshine-sys/bcr"
	bcr2 "github.com/starshine-sys/bcr/v2"

	// imports for user-supplied images
	_ "image/gif"
	_ "image/jpeg"
)

//go:embed pride
var prideFS embed.FS

var possibleFlags []string

func init() {
	entries, err := prideFS.ReadDir("pride")
	if err != nil {
		panic(err)
	}

	for _, e := range entries {
		if e.Type().IsDir() {
			continue
		}

		possibleFlags = append(possibleFlags, strings.TrimSuffix(e.Name(), filepath.Ext(e.Name())))
	}

	sort.Strings(possibleFlags)
}

func (bot *Bot) prideSlash(ctx *bcr2.CommandContext) (err error) {
	if err := ctx.Defer(); err != nil {
		return err
	}

	url := ctx.User.AvatarURLWithType(discord.PNGImage) + "?size=1024"
	flagName := ctx.Options.Find("flag").String()

	if id := ctx.Options.Find("pk-member").String(); id != "" {
		m, err := bot.PK.Member(strings.ToLower(id))
		if err == nil && m.AvatarURL != "" {
			url = m.AvatarURL
		}
	}

	var flagFile fs.DirEntry
	entries, _ := prideFS.ReadDir("pride")
	for _, de := range entries {
		if de.Type().IsDir() {
			continue
		}

		if strings.ToLower(strings.TrimSuffix(de.Name(), filepath.Ext(de.Name()))) == flagName {
			flagFile = de
			break
		}
	}

	if flagFile == nil {
		return ctx.Reply(fmt.Sprintf("The following flags are available: %v", strings.Join(possibleFlags, ", ")))
	}

	buf, err := bot.prideFlag(url, flagFile)
	if err != nil {
		return bot.ReportInteraction(ctx, err)
	}

	return ctx.ReplyComplex(api.InteractionResponseData{
		Files: []sendpart.File{{Name: ctx.User.ID.String() + ".png", Reader: buf}},
	})
}

func (bot *Bot) prideFlag(url string, file fs.DirEntry) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	pfp, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	flagb, err := prideFS.Open(filepath.Join("pride/", file.Name()))
	if err != nil {
		return nil, err
	}
	defer flagb.Close()

	flag, err := png.Decode(flagb)
	if err != nil {
		return nil, err
	}

	size := pfp.Bounds().Max.X - pfp.Bounds().Min.X
	img := gg.NewContext(size, size)

	img.DrawImageAnchored(pfp, size/2, size/2, 0.5, 0.5)

	flag = imaging.Resize(flag, size, size, imaging.Linear)

	img.DrawImageAnchored(flag, size/2, size/2, 0.5, 0.5)

	buf := new(bytes.Buffer)

	err = img.EncodePNG(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (bot *Bot) pride(ctx *bcr.Context) (err error) {
	url := ctx.Author.AvatarURLWithType(discord.PNGImage) + "?size=1024"
	flagName := strings.ToLower(ctx.RawArgs)

	var flagFile fs.DirEntry
	entries, _ := prideFS.ReadDir("pride")
	for _, de := range entries {
		if de.Type().IsDir() {
			continue
		}

		if strings.ToLower(strings.TrimSuffix(de.Name(), filepath.Ext(de.Name()))) == flagName {
			flagFile = de
			break
		}
	}

	if flagFile == nil {
		return ctx.SendfX("The following flags are available: %v", strings.Join(possibleFlags, ", "))
	}

	filename := strings.ToLower(ctx.User().Username)
	if len(ctx.Message.Attachments) > 0 {
		if strings.HasSuffix(ctx.Message.Attachments[0].Filename, ".png") ||
			strings.HasSuffix(ctx.Message.Attachments[0].Filename, ".gif") ||
			strings.HasSuffix(ctx.Message.Attachments[0].Filename, ".jpg") ||
			strings.HasSuffix(ctx.Message.Attachments[0].Filename, ".jpeg") {
			url = ctx.Message.Attachments[0].URL
			filename = strings.TrimSuffix(ctx.Message.Attachments[0].Filename, filepath.Ext(ctx.Message.Attachments[0].Filename))

			if ctx.Message.Attachments[0].Size > 1*1024*1024 {
				return ctx.SendfX("That file is too big, sorry. (%v > 1 MB)", humanize.Bytes(ctx.Message.Attachments[0].Size))
			}
		}
	}

	buf, err := bot.prideFlag(url, flagFile)
	if err != nil {
		return bot.Report(ctx, err)
	}

	filename = filename + "-" + flagName + ".png"
	return ctx.SendFiles("", sendpart.File{
		Name:   filename,
		Reader: buf,
	})
}
