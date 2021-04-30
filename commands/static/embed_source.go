package static

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) embedSource(ctx *bcr.Context) (err error) {
	if !linkRegex.MatchString(ctx.RawArgs) {
		_, err = ctx.Send("You didn't give a valid message ID/link.", nil)
		return
	}

	groups := linkRegex.FindStringSubmatch(ctx.RawArgs)
	channelID, _ := discord.ParseSnowflake(groups[2])
	msgID, _ := discord.ParseSnowflake(groups[3])

	msg, err := bot.State.Message(discord.ChannelID(channelID), discord.MessageID(msgID))
	if err != nil {
		fmt.Println(err)
		_, err = ctx.Send("Could not find that message. Are you sure I have access to that channel?", nil)
		return
	}

	if len(msg.Embeds) == 0 {
		_, err = ctx.Send("That message has no embeds.", nil)
	}

	var embedJSON [][]byte

	for _, e := range msg.Embeds {
		b, err := json.MarshalIndent(e, "", "    ")
		if err != nil {
			return bot.Report(ctx, err)
		}
		embedJSON = append(embedJSON, b)
	}

	for _, e := range embedJSON {
		if len(e) > 2000 {
			_, err = ctx.NewMessage().AddFile("embed.json", bytes.NewReader(e)).Send()
			if err != nil {
				return err
			}
		} else {
			_, err = ctx.Send("", &discord.Embed{
				Title:       "Source",
				Description: "```" + string(e) + "```",
				Color:       ctx.Router.EmbedColor,
			})
			if err != nil {
				return err
			}
		}
	}

	return
}
