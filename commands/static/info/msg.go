package info

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/starshine-sys/bcr"
)

var (
	msgMap   = map[discord.MessageID]discord.ChannelID{}
	msgMapMu sync.Mutex
)

func (bot *Bot) message(ctx *bcr.Context) (err error) {
	var msg *discord.Message

	if sf, err := discord.ParseSnowflake(ctx.Args[0]); err == nil {
		msgMapMu.Lock()
		chID, ok := msgMap[discord.MessageID(sf)]
		msgMapMu.Unlock()

		if ok {
			msg, err = ctx.State.Message(chID, discord.MessageID(sf))
		} else {
			m, err := ctx.Sendf("I couldn't find that message with just the ID.")
			if err != nil {
				return err
			}
			time.AfterFunc(5*time.Second, func() {
				ctx.State.DeleteMessage(ctx.Channel.ID, m.ID, "")
			})
			return err
		}
	} else {
		msg, err = ctx.ParseMessage(ctx.Args[0])
	}

	if err != nil || msg == nil {
		m, err := ctx.Sendf("Couldn't find that message.")
		if err != nil {
			return err
		}
		time.AfterFunc(5*time.Second, func() {
			ctx.State.DeleteMessage(ctx.Channel.ID, m.ID, "")
		})
		return err
	}

	if msg.GuildID != ctx.Message.GuildID {
		m, err := ctx.Sendf("That message isn't from this server.")
		if err != nil {
			return err
		}
		time.AfterFunc(5*time.Second, func() {
			ctx.State.DeleteMessage(ctx.Channel.ID, m.ID, "")
		})
		return err
	}

	perms, _ := ctx.State.Permissions(msg.ChannelID, ctx.Author.ID)
	if !perms.Has(discord.PermissionViewChannel) || !perms.Has(discord.PermissionReadMessageHistory) {
		m, err := ctx.Sendf("Nice try, but you can't see that channel.")
		if err != nil {
			return err
		}
		time.AfterFunc(5*time.Second, func() {
			ctx.State.DeleteMessage(ctx.Channel.ID, m.ID, "")
		})
		return err
	}

	clr := discord.Color(bcr.ColourBlurple)
	name := msg.Author.Username
	if !msg.WebhookID.IsValid() {
		member, err := bot.Member(msg.GuildID, msg.Author.ID)
		if err == nil {
			clr = discord.MemberColor(*ctx.Guild, member)
			if member.Nick != "" {
				name = member.Nick
			}
		}
	}

	e := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: name,
			Icon: msg.Author.AvatarURLWithType(discord.PNGImage),
		},
		Description: msg.Content,
		Color:       clr,
		Timestamp:   msg.Timestamp,
		Fields: []discord.EmbedField{{
			Name:  "Source",
			Value: "[Jump to message](" + msg.URL() + ")",
		}},
	}

	if len(msg.Attachments) > 0 {
		if isImage(msg.Attachments[0].Filename) {
			e.Image = &discord.EmbedImage{
				URL: msg.Attachments[0].URL,
			}
		}
	}

	ch, err := ctx.State.Channel(msg.ChannelID)
	if err == nil {
		e.Footer = &discord.EmbedFooter{
			Text: "#" + ch.Name,
		}
	}

	m, err := ctx.Send("", e)
	if err != nil {
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "insert into command_responses (message_id, user_id) values ($1, $2)", m.ID, ctx.Author.ID)
	return
}

func (bot *Bot) messageCreate(ev *gateway.MessageCreateEvent) {
	msgMapMu.Lock()
	msgMap[ev.ID] = ev.ChannelID
	msgMapMu.Unlock()
}

func isImage(name string) bool {
	switch {
	case strings.HasSuffix(name, ".jpg"), strings.HasSuffix(name, ".jpeg"), strings.HasSuffix(name, ".gif"), strings.HasSuffix(name, ".png"), strings.HasSuffix(name, ".webp"):
		return true
	default:
		return false
	}
}
