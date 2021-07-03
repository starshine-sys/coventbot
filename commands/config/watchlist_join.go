package config

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/starshine-sys/tribble/etc"
)

func (bot *Bot) watchlistMemberAdd(m *gateway.GuildMemberAddEvent) {
	var ch discord.ChannelID
	err := bot.DB.Pool.QueryRow(context.Background(), "select watch_list_channel from servers where id = $1", m.GuildID).Scan(&ch)
	if err != nil {
		bot.Sugar.Errorf("Error getting watchlist channel for %v: %v", m.GuildID, err)
		return
	}

	// if there's no watch list channel set, return
	if !ch.IsValid() {
		return
	}

	// if the user isn't on the watch list, return
	if !bot.DB.IsWatchlisted(m.GuildID, m.User.ID) {
		return
	}

	e := discord.Embed{
		Title: "User on watch list joined",
		Color: etc.ColourOrange,
		Author: &discord.EmbedAuthor{
			Name: m.User.Username + "#" + m.User.Discriminator,
			Icon: m.User.AvatarURL(),
		},
		Thumbnail: &discord.EmbedThumbnail{
			URL: m.User.AvatarURL(),
		},
		Description: fmt.Sprintf("⚠️ **%v#%v** just joined the server and is on the watch list.", m.User.Username, m.User.Discriminator),
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("User ID: %v", m.User.ID),
		},
		Timestamp: discord.NowTimestamp(),
	}

	var reason string
	bot.DB.Pool.QueryRow(context.Background(), "select reason from watch_list_reasons where user_id = $1 and server_id = $2", m.User.ID, m.GuildID).Scan(&reason)

	if reason != "" {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Reason",
			Value: reason,
		})
	}

	s, _ := bot.Router.StateFromGuildID(m.GuildID)

	_, err = s.SendEmbeds(ch, e)
	if err != nil {
		bot.Sugar.Errorf("Error sending watch list warning: %v", err)
	}
	return
}
