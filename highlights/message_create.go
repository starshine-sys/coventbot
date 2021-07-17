package highlights

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/caneroj1/stemmer"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) messageCreate(m *gateway.MessageCreateEvent) {
	if !m.GuildID.IsValid() || m.Author.Bot {
		return
	}

	s, _ := bot.Router.StateFromGuildID(m.GuildID)

	guildConf, err := bot.hlConfig(m.GuildID)
	if err != nil {
		bot.Sugar.Errorf("Error getting guild config: %v", err)
		return
	}

	if !guildConf.HlEnabled {
		return
	}

	g, err := s.Guild(m.GuildID)
	if err != nil {
		bot.Sugar.Errorf("Error getting guild: %v", err)
		return
	}

	g.Roles, err = s.Roles(m.GuildID)
	if err != nil {
		bot.Sugar.Errorf("Error getting roles: %v", err)
		return
	}

	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		bot.Sugar.Errorf("Error getting channel: %v", err)
		return
	}

	bot.UserExpirationMu.Lock()
	bot.UserExpiration[userExpirationKey{m.Author.ID, m.GuildID}] = time.Now().Add(5 * time.Minute)
	bot.UserExpirationMu.Unlock()

	for _, id := range guildConf.Blocked {
		if id == uint64(ch.ID) || id == uint64(ch.CategoryID) {
			return
		}
	}

	var parent *discord.Channel
	if ch.Type == discord.GuildNewsThread || ch.Type == discord.GuildPrivateThread || ch.Type == discord.GuildPublicThread {
		parent, err = s.Channel(ch.CategoryID)
		if err != nil {
			bot.Sugar.Errorf("Error getting channel: %v", err)
			return
		}

		for _, id := range guildConf.Blocked {
			if id == uint64(parent.CategoryID) {
				return
			}
		}
	}

	if m.Content == "" {
		return
	}

	split := strings.Fields(m.Content)
	stemmer.StemConcurrent(&split)

	var messages []discord.Message = nil

	hls, err := bot.guildHighlights(m.GuildID)
	if err != nil {
		bot.Sugar.Errorf("Error getting highlights: %v", err)
		return
	}

	for _, hl := range hls {
		if hl.UserID == m.Author.ID {
			continue
		}

		bot.UserExpirationMu.Lock()
		exp, ok := bot.UserExpiration[userExpirationKey{hl.UserID, m.GuildID}]
		if ok {
			if time.Now().Before(exp) {
				bot.UserExpirationMu.Unlock()
				continue
			}
			delete(bot.UserExpiration, userExpirationKey{hl.UserID, m.GuildID})
		}
		bot.UserExpirationMu.Unlock()

		member, err := bot.Member(m.GuildID, hl.UserID)
		if err != nil {
			continue
		}

		if !discord.CalcOverwrites(*g, *ch, member).Has(discord.PermissionViewChannel) {
			continue
		}

		for _, id := range hl.Blocked {
			if id == uint64(ch.ID) || id == uint64(ch.CategoryID) || id == uint64(m.Author.ID) {
				continue
			}

			if parent != nil {
				if id == uint64(parent.CategoryID) {
					continue
				}
			}
		}

		words := stemmer.StemMultiple(hl.Highlights)
		for i, word := range words {
			bot.WordExpirationMu.Lock()
			exp, ok := bot.WordExpiration[wordExpirationKey{hl.UserID, m.GuildID, word}]
			if ok {
				if time.Now().Before(exp) {
					bot.WordExpirationMu.Unlock()
					continue
				}
				delete(bot.WordExpiration, wordExpirationKey{hl.UserID, m.GuildID, word})
			}
			bot.WordExpirationMu.Unlock()

			for _, msg := range split {
				if word == msg {
					if messages == nil {
						messages, err = s.Messages(m.ChannelID, 5)
						if err != nil {
							bot.Sugar.Errorf("Error fetching messages: %v", err)
							return
						}
						for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
							messages[i], messages[j] = messages[j], messages[i]
						}
					}

					bot.WordExpirationMu.Lock()
					bot.WordExpiration[wordExpirationKey{hl.UserID, m.GuildID, word}] = time.Now().Add(time.Minute)
					bot.WordExpirationMu.Unlock()

					err = bot.makeAndSendHighlight(s, g, member.User.ID, m.ChannelID, hl.Highlights[i], messages)
					if err != nil {
						bot.Sugar.Errorf("Error sending message to %v (%v): %v", member.User.Tag(), member.User.ID, err)
						continue
					}
				}
			}
		}
	}
}

func (bot *Bot) makeAndSendHighlight(s *state.State, g *discord.Guild, userID discord.UserID, channelID discord.ChannelID, word string, messages []discord.Message) (err error) {
	data := api.SendMessageData{
		Content: fmt.Sprintf("In **%v** %v, you were mentioned with highlight word \"%v\"", g.Name, channelID.Mention(), word),
		Embeds: []discord.Embed{{
			Title: strings.Title(word),
			Fields: []discord.EmbedField{{
				Name:  "Source",
				Value: fmt.Sprintf("[Jump!](https://discord.com/channels/%v/%v/%v)", g.ID, channelID, messages[len(messages)-1].ID),
			}},
			Footer: &discord.EmbedFooter{
				Icon: messages[len(messages)-1].Author.AvatarURLWithType(discord.PNGImage),
				Text: messages[len(messages)-1].Author.Username,
			},
			Timestamp: discord.NowTimestamp(),
			Color:     bcr.ColourBlurple,
		}},
		AllowedMentions: &api.AllowedMentions{
			Parse: []api.AllowedMentionType{},
		},
	}

	for _, m := range messages {
		content := m.Content
		if len(m.Content) > 200 {
			content = m.Content[:200] + "..."
		}

		data.Embeds[0].Description += fmt.Sprintf("**`[%v]` %v:** %v\n\n", m.Timestamp.Time().UTC().Format("15:04:05"), m.Author.Username, content)
	}

	data.Embeds[0].Description = strings.TrimSpace(data.Embeds[0].Description)

	ch, err := s.CreatePrivateChannel(userID)
	if err != nil {
		return err
	}

	msg, err := s.SendMessageComplex(ch.ID, data)
	if err != nil {
		return err
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "insert into highlight_delete_queue (message_id, channel_id) values ($1, $2)", msg.ID, ch.ID)
	return err
}
