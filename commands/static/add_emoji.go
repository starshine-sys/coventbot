package static

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

var emojiMatch = regexp.MustCompile("<(?P<animated>a)?:(?P<name>\\w+):(?P<emoteID>\\d{15,})>")

func (bot *Bot) addEmoji(ctx *bcr.Context) (err error) {
	if ctx.RawArgs == "-h" || ctx.RawArgs == "" {
		e := &discord.Embed{
			Title: "`ADDEMOJI`",
			Description: "`<source> [name]`" + `
Available sources are listed below. Name is optional when using the msg, existing, or attachment sources.

` + "`-msg <link/channel-message>`" + `
> Add an emoji used in the given message. The bot needs to be able to see the channel and read message history in it.
> If the message contains more than one emoji, shows a selection menu.
` + "`-existing <emoji>`" + `
> Add an existing emoji given as input.
` + "`-url <url>`" + `
> Add the image file at the given URL.
` + "`-attachment`" + `
> Add the image file attached to the message.`,
			Color: ctx.Router.EmbedColor,
		}
		_, err = ctx.Send("", e)
		return
	}

	var (
		name     string
		fileType string

		msg        string
		existing   string
		url        string
		attachment bool
	)

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.BoolVar(&attachment, "attachment", false, "Use the given attachment as source")
	fs.StringVar(&msg, "msg", "", "Use the given message link as source")
	fs.StringVar(&url, "url", "", "Use the given URL as source")
	fs.StringVar(&existing, "existing", "", "Use the given emoji as source")

	err = fs.Parse(ctx.Args)
	if err != nil {
		_, err = ctx.Send("Invalid source. Please double-check your input.", nil)
		return
	}
	ctx.Args = fs.Args()

	if !attachment && msg == "" && existing == "" && url == "" {
		_, err = ctx.Send("No source given. Please double-check your input.", nil)
		return
	}
	if attachment && url != "" || existing != "" && url != "" {
		_, err = ctx.Send("Too many sources given. Please double-check your input.", nil)
		return
	}

	if msg != "" {
		return bot.addEmojiMsg(ctx, msg)
	}

	if url != "" {
		if isImage(url) {
			fileType = imageFileType(url)
		}
	}

	if attachment {
		if len(ctx.Message.Attachments) == 0 {
			_, err = ctx.Send("`-attachment` flag specified, but the message has no attachments.", nil)
			return
		}

		for _, a := range ctx.Message.Attachments {
			if !isImage(a.Filename) {
				continue
			}
			url = a.URL
			name = imageFileName(a.Filename)
			fileType = imageFileType(a.Filename)
			break
		}

		if len(name) > 31 {
			name = name[:31]
		}
	}

	if existing != "" {
		if !emojiMatch.MatchString(existing) {
			_, err = ctx.Send("No valid custom emoji given.", nil)
			return
		}

		extension := ".png"
		groups := emojiMatch.FindStringSubmatch(existing)
		if groups[1] == "a" {
			extension = ".gif"
		}
		name = groups[2]
		url = fmt.Sprintf("https://cdn.discordapp.com/emojis/%v%v", groups[3], extension)
		fileType = imageFileType(url)
	}

	if url != "" {
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

		if len(ctx.Args) > 0 {
			name = ctx.Args[0]
		}

		img := api.Image{
			ContentType: fileType,
			Content:     b,
		}

		if name == "" {
			_, err = ctx.Send("No name given. Please double-check your input.", nil)
			return err
		}

		ced := api.CreateEmojiData{
			Name:  name,
			Image: img,
		}

		emoji, err := ctx.Session.CreateEmoji(ctx.Message.GuildID, ced)
		if err != nil {
			_, err = ctx.Sendf("Error:\n```%v```", err)
			return err
		}

		_, err = ctx.Sendf("Added emoji %v with name \"%v\".", emoji.String(), emoji.Name)
		return err
	}

	return
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
