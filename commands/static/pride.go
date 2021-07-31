package static

import (
	"bytes"
	"embed"
	"image"
	"image/png"
	"io/fs"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/disintegration/imaging"
	"github.com/dustin/go-humanize"
	"github.com/fogleman/gg"
	"github.com/starshine-sys/bcr"

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

func (bot *Bot) pride(ctx bcr.Contexter) (err error) {
	flagName := ""
	if v, ok := ctx.(*bcr.Context); ok {
		flagName = strings.ToLower(v.RawArgs)
	} else if sv, ok := ctx.(*bcr.SlashContext); ok {
		for _, o := range sv.CommandOptions {
			flagName = strings.ToLower(o.Value)
			break
		}
	}

	var flagFile fs.DirEntry
	entries, err := prideFS.ReadDir("pride")
	if err != nil {
		return bot.Report(ctx, err)
	}
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
	url := ctx.User().AvatarURLWithType(discord.PNGImage) + "?size=1024"
	if v, ok := ctx.(*bcr.Context); ok {
		if len(v.Message.Attachments) > 0 {
			if strings.HasSuffix(v.Message.Attachments[0].Filename, ".png") ||
				strings.HasSuffix(v.Message.Attachments[0].Filename, ".gif") ||
				strings.HasSuffix(v.Message.Attachments[0].Filename, ".jpg") ||
				strings.HasSuffix(v.Message.Attachments[0].Filename, ".jpeg") {
				url = v.Message.Attachments[0].URL
				filename = strings.TrimSuffix(v.Message.Attachments[0].Filename, filepath.Ext(v.Message.Attachments[0].Filename))

				if v.Message.Attachments[0].Size > 1*1024*1024 {
					return ctx.SendfX("That file is too big, sorry. (%v > 1 MB)", humanize.Bytes(v.Message.Attachments[0].Size))
				}
			}
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return bot.Report(ctx, err)
	}
	defer resp.Body.Close()

	pfp, _, err := image.Decode(resp.Body)
	if err != nil {
		return bot.Report(ctx, err)
	}

	flagb, err := prideFS.Open(filepath.Join("pride/", flagFile.Name()))
	if err != nil {
		return bot.Report(ctx, err)
	}
	defer flagb.Close()

	flag, err := png.Decode(flagb)
	if err != nil {
		return bot.Report(ctx, err)
	}

	size := pfp.Bounds().Max.X - pfp.Bounds().Min.X
	img := gg.NewContext(size, size)

	img.DrawImageAnchored(pfp, size/2, size/2, 0.5, 0.5)

	flag = imaging.Resize(flag, size, size, imaging.Linear)

	img.DrawImageAnchored(flag, size/2, size/2, 0.5, 0.5)

	buf := new(bytes.Buffer)

	err = img.EncodePNG(buf)
	if err != nil {
		return bot.Report(ctx, err)
	}

	filename = filename + "-" + flagName + ".png"
	return ctx.SendFiles("", sendpart.File{
		Name:   filename,
		Reader: buf,
	})
}
