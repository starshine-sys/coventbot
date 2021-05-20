package levels

import (
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

func (bot *Bot) messageCreate(m *gateway.MessageCreateEvent) {
	if !m.GuildID.IsValid() || m.Author.Bot || m.Author.DiscordSystem || m.Member == nil {
		return
	}

	sc, err := bot.getGuildConfig(m.GuildID)
	if err != nil {
		bot.Sugar.Errorf("Error getting guild config: %v", err)
		return
	}

	if !sc.LevelsEnabled {
		return
	}

	uc, err := bot.getUser(m.GuildID, m.Author.ID)
	if err != nil {
		bot.Sugar.Errorf("Error getting user: %v", err)
		return
	}

	if uc.NextTime.After(time.Now().UTC()) {
		return
	}

	// check blocked channels/categories
	for _, blocked := range sc.BlockedChannels {
		if m.ChannelID == discord.ChannelID(blocked) {
			return
		}
	}
	if ch, err := bot.State.Channel(m.ChannelID); err == nil {
		for _, blocked := range sc.BlockedCategories {
			if ch.CategoryID == discord.ChannelID(blocked) {
				return
			}
		}
	}

	// check blocked roles
	for _, blocked := range sc.BlockedRoles {
		for _, r := range m.Member.RoleIDs {
			if discord.RoleID(blocked) == r {
				return
			}
		}
	}

	// increment the user's xp!
	newXP, err := bot.incrementXP(m.GuildID, m.Author.ID, sc.BetweenXP)
	if err != nil {
		bot.Sugar.Errorf("Error updating XP for user: %v", err)
		return
	}

	// only check for rewards on level up
	oldLvl := currentLevel(uc.XP)
	newLvl := currentLevel(newXP)

	if oldLvl >= newLvl {
		return
	}

	reward := bot.getReward(m.GuildID, newLvl)
	if reward == nil {
		return
	}

	if !reward.RoleReward.IsValid() {
		return
	}

	err = bot.State.AddRole(m.GuildID, m.Author.ID, reward.RoleReward)
	if err != nil {
		bot.Sugar.Errorf("Error adding role to user: %v", err)
		return
	}

	if sc.RewardText != "" {
		txt := strings.NewReplacer("{lvl}", fmt.Sprint(newLvl)).Replace(sc.RewardText)

		ch, err := bot.State.CreatePrivateChannel(m.Author.ID)
		if err == nil {
			bot.State.SendText(ch.ID, txt)
		}
	}
}
