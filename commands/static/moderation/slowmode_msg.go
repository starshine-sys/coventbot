package moderation

import (
	"context"
	"fmt"
	"time"

	"github.com/diamondburned/arikawa/v2/gateway"
)

func (bot *Bot) slowmodeMessage(m *gateway.MessageCreateEvent) {
	if !m.GuildID.IsValid() || m.Author.Bot || m.Member == nil {
		return
	}

	hasSlowmode, duration := bot.hasSlowmode(m.ChannelID)
	if !hasSlowmode {
		return
	}

	if ignore, _ := bot.slowmodeIgnore(m.GuildID, m.Member.RoleIDs); ignore {
		return
	}

	delete := bot.userSlowmode(m.ChannelID, m.Author.ID)
	if delete {
		err := bot.Router.State.DeleteMessage(m.ChannelID, m.ID)
		if err != nil {
			bot.Sugar.Errorf("Error deleting message: %v", err)
			return
		}

		var expiry time.Time
		err = bot.DB.Pool.QueryRow(context.Background(), "select expiry from user_slowmode where channel_id = $1 and user_id = $2", m.ChannelID, m.Author.ID).Scan(&expiry)
		if err != nil {
			bot.Sugar.Errorf("Error getting expiry time: %v", err)
			return
		}

		msg := fmt.Sprintf("You can send your next message in %v at %v.", m.ChannelID.Mention(), expiry.UTC().Format("15:04:05 UTC, January 02 2006"))

		ch, err := bot.State.CreatePrivateChannel(m.Author.ID)
		if err != nil {
			bot.Sugar.Errorf("Error creating private channel for %v: %v", m.Author.ID, err)
			return
		}

		bot.State.SendText(ch.ID, msg)
		return
	}

	expiry := time.Now().UTC().Add(duration)

	err := bot.setUserSlowmode(m.GuildID, m.ChannelID, m.Author.ID, expiry)
	if err != nil {
		bot.Sugar.Errorf("Error setting slowmode for %v in %v: %v", m.Author.ID, m.ChannelID, err)
	}
}
