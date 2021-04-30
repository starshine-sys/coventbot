package static

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/starshine-sys/bcr"
)

type emoji struct {
	name     string
	animated bool
	url      string
}

const emojiBaseURL = "https://cdn.discordapp.com/emojis/"

var (
	idRegex   = regexp.MustCompile(`(?P<channel_id>[0-9]{15,20})-(?P<message_id>[0-9]{15,20})$`)
	linkRegex = regexp.MustCompile(`https?://(?:(ptb|canary|www)\.)?discord(?:app)?\.com/channels/(?:[0-9]{15,20}|@me)/(?P<channel_id>[0-9]{15,20})/(?P<message_id>[0-9]{15,20})/?$`)
)

func (bot *Bot) addEmojiMsg(ctx *bcr.Context, url string) (err error) {
	var groups []string

	if idRegex.MatchString(url) {
		groups = idRegex.FindStringSubmatch(url)
	} else if linkRegex.MatchString(url) {
		groups = linkRegex.FindStringSubmatch(url)
		groups = groups[1:]
	}

	if len(groups) == 0 {
		_, err = ctx.Send("Invalid message link/ID given.", nil)
		return
	}

	channel, _ := discord.ParseSnowflake(groups[1])
	msgID, _ := discord.ParseSnowflake(groups[2])

	msg, err := ctx.State.Message(discord.ChannelID(channel), discord.MessageID(msgID))
	if err != nil {
		_, err = ctx.Send("Couldn't find that message. Are you sure I have access to the channel?", nil)
		return
	}

	// find emojis
	emojis := emojiMatch.FindAllString(msg.Content, -1)
	if emojis == nil {
		_, err = ctx.Send("That message has no custom emoji.", nil)
		return
	}

	emojiObjects := make([]emoji, 0)
	emojiEmbeds := make([]discord.Embed, 0)

	// loop through all emojis
	for _, e := range emojis {
		em := emoji{}
		ext := ".png"
		groups := emojiMatch.FindStringSubmatch(e)
		if groups[1] == "a" {
			ext = ".gif"
			em.animated = true
		}
		em.name = groups[2]
		em.url = emojiBaseURL + groups[3] + ext
		emojiObjects = append(emojiObjects, em)

		emojiEmbeds = append(emojiEmbeds, discord.Embed{
			Title: em.name,
			Image: &discord.EmbedImage{
				URL: em.url,
			},
			Footer: &discord.EmbedFooter{
				Text: fmt.Sprintf("Animated: %v | React with ✅ to select this emoji.", em.animated),
			},
		})
	}

	var (
		e      emoji
		v      interface{}
		c      context.Context
		cancel context.CancelFunc

		embedMsg *discord.Message
	)

	if len(emojiObjects) == 1 {
		e = emojiObjects[0]
		goto gotEmoji
	}

	embedMsg, err = ctx.PagedEmbed(emojiEmbeds, false)
	if err != nil {
		return err
	}
	err = ctx.State.React(embedMsg.ChannelID, embedMsg.ID, discord.APIEmoji("✅"))
	if err != nil {
		bot.Sugar.Errorf("Error adding reaction: %v", err)
	}

	c, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	v = ctx.State.WaitFor(c, func(v interface{}) bool {
		ev, ok := v.(*gateway.MessageReactionAddEvent)
		if !ok {
			return false
		}

		return ev.MessageID == embedMsg.ID && ev.UserID == ctx.Author.ID && ev.Emoji.APIString() == "✅"
	})

	if v == nil {
		_, err = ctx.Send("Timed out.", nil)
		return
	}

	e = emojiObjects[ctx.AdditionalParams["page"].(int)]

gotEmoji:

	resp, err := http.Get(e.url)
	if err != nil {
		return bot.Report(ctx, err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return bot.Report(ctx, err)
	}

	ced := api.CreateEmojiData{
		Name: e.name,
		Image: api.Image{
			ContentType: imageFileType(e.url),
			Content:     b,
		},
	}

	emoji, err := ctx.State.CreateEmoji(ctx.Message.GuildID, ced)
	if err != nil {
		return bot.Report(ctx, err)
	}

	_, err = ctx.Sendf("Added emoji %v with name \"%v\".", emoji.String(), emoji.Name)
	return err
}
