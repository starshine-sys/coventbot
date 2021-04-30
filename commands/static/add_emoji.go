package static

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	gourl "net/url"
	"path"
	"regexp"
	"strings"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/starshine-sys/bcr"
)

var emojiMatch = regexp.MustCompile("<(?P<animated>a)?:(?P<name>\\w+):(?P<emoteID>\\d{15,})>")

func (bot *Bot) addEmoji(ctx *bcr.Context) (err error) {
	var (
		name     string
		fileType string

		msg string
		url string
	)

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&msg, "msg", "", "Use the given message link as source")

	err = fs.Parse(ctx.Args)
	if err != nil {
		return
	}
	ctx.Args = fs.Args()

	// first we try attachments
	for _, a := range ctx.Message.Attachments {
		if !isImage(a.Filename) {
			continue
		}
		url = a.URL
		name = imageFileName(a.Filename)
		if len(ctx.Args) > 0 {
			name = ctx.Args[0]
		}
		fileType = imageFileType(a.Filename)
		break
	}
	if url != "" {
		goto url
	}

	if len(ctx.Args) > 0 {
		if emojiMatch.MatchString(ctx.Args[0]) {
			extension := ".png"
			groups := emojiMatch.FindStringSubmatch(ctx.Args[0])
			if groups[1] == "a" {
				extension = ".gif"
			}
			name = groups[2]
			url = fmt.Sprintf("https://cdn.discordapp.com/emojis/%v%v", groups[3], extension)
			fileType = imageFileType(url)
			goto url
		}

		if _, err := gourl.Parse(ctx.Args[0]); err == nil {
			url = ctx.Pop()
			if ctx.Peek() != "" {
				name = ctx.Pop()
			}
			goto url
		}
	}

	if msg != "" {
		return bot.addEmojiMsg(ctx, msg)
	}

url:
	if url == "" {
		_, err = ctx.Send("No URL given.", nil)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		_, err = ctx.Sendf("Internal error occurred:\n```%v```", err)
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		_, err = ctx.Sendf("Internal error occurred:\n```%v```", err)
		return err
	}

	if len(ctx.Args) > 1 {
		name = ctx.Args[1]
	}

	fileType = imageFileType(url)

	img := api.Image{
		ContentType: fileType,
		Content:     b,
	}

	if name == "" {
		u, err := gourl.Parse(url)
		if err != nil {
			_, err = ctx.Send("Couldn't parse the URL!", nil)
			return err
		}

		name = imageFileName(path.Base(u.Path))
	}

	ced := api.CreateEmojiData{
		Name:  name,
		Image: img,
	}

	emoji, err := ctx.State.CreateEmoji(ctx.Message.GuildID, ced)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Sendf("Added emoji %v with name \"%v\".", emoji.String(), emoji.Name)
	return err
}

func isImage(filename string) bool {
	filename = strings.TrimSuffix(filename, "?v=1")

	return strings.HasSuffix(filename, ".png") ||
		strings.HasSuffix(filename, ".jpg") ||
		strings.HasSuffix(filename, ".jpeg") ||
		strings.HasSuffix(filename, ".gif")
}

func imageFileName(filename string) string {
	filename = strings.TrimSuffix(filename, "?v=1")
	filename = strings.TrimSuffix(filename, ".gif")
	filename = strings.TrimSuffix(filename, ".jpg")
	filename = strings.TrimSuffix(filename, ".jpeg")
	filename = strings.TrimSuffix(filename, ".png")
	return filename
}

func imageFileType(filename string) string {
	filename = strings.TrimSuffix(filename, "?v=1")
	if strings.HasSuffix(filename, ".png") {
		return "image/png"
	}
	if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
		return "image/jpeg"
	}
	if strings.HasSuffix(filename, ".gif") {
		return "image/gif"
	}
	return "unknown"
}
