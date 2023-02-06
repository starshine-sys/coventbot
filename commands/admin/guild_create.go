// SPDX-License-Identifier: AGPL-3.0-only
package admin

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) guildCreate(ev *gateway.GuildCreateEvent) {
	if !bot.loadedAllowedGuilds {
		_, err := bot.DB.Pool.Exec(context.Background(), "insert into allowed_guilds (id, reason, added_by, added_for) values ($1, $2, $3, $3) on conflict (id) do nothing", ev.ID, "Automatically added in first session", bot.Router.Bot.ID)
		if err != nil {
			bot.Sugar.Errorf("Error saving guild %v: %v", ev.ID, err)
		}

		bot.Sugar.Infof("Added guild %v (%v) to allowed guilds list", ev.Name, ev.ID)

		return
	}

	if bot.isGuildAllowed(ev.ID) {
		return
	}

	bot.Sugar.Warnf("Leaving guild %v (%v) because it is not on the allow list.", ev.Name, ev.ID)

	s, _ := bot.Router.StateFromGuildID(ev.ID)

	if ev.OwnerID.IsValid() {
		ch, err := s.CreatePrivateChannel(ev.OwnerID)
		if err == nil {
			_, err = s.SendEmbeds(ch.ID, discord.Embed{
				Description: fmt.Sprintf("I have left your server **%v** as it's not on my list of allowed servers. Please contact my owner (<@%v>) if you think this is in error.", ev.Name, bot.Router.BotOwners[0]),
				Color:       bcr.ColourRed,
			})
			if err != nil {
				bot.Sugar.Errorf("Error sending leave message to owner of %v: %v", ev.OwnerID, err)
			}
		}
	}

	err := s.LeaveGuild(ev.ID)
	if err != nil {
		bot.Sugar.Errorf("Error leaving guild %v: %v", err)
	}
}

func (bot *Bot) isGuildAllowed(id discord.GuildID) (allowed bool) {
	err := bot.DB.Pool.QueryRow(context.Background(), "select exists(select * from allowed_guilds where id = $1)", id).Scan(&allowed)
	if err != nil {
		bot.Sugar.Errorf("Error checking if guild %v is allowed: %v", id, err)
		return true
	}
	return allowed
}
