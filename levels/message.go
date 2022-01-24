package levels

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) messageCreate(m *gateway.MessageCreateEvent) {
	if !m.GuildID.IsValid() || m.Author.Bot || m.Author.DiscordSystem || m.Member == nil {
		return
	}

	s, _ := bot.Router.StateFromGuildID(m.GuildID)

	sc, err := bot.getGuildConfig(m.GuildID)
	if err != nil {
		bot.Sugar.Errorf("Error getting guild config: %v", err)
		return
	}

	if !sc.LevelsEnabled {
		return
	}

	if bot.isBlacklisted(m.GuildID, m.Author.ID) {
		return
	}

	uc, err := bot.getUser(m.GuildID, m.Author.ID)
	if err != nil {
		bot.Sugar.Errorf("Error getting user: %v", err)
		return
	}

	if uc.NextTime.Add(sc.BetweenXP).After(time.Now().UTC()) {
		return
	}

	// check blocked channels/categories
	for _, blocked := range sc.BlockedChannels {
		if m.ChannelID == discord.ChannelID(blocked) {
			return
		}
	}
	if ch, err := s.Channel(m.ChannelID); err == nil {
		for _, blocked := range sc.BlockedCategories {
			if ch.ParentID == discord.ChannelID(blocked) {
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
	newXP, err := bot.incrementXP(m.GuildID, m.Author.ID, sc.CarlineCompatible)
	if err != nil {
		bot.Sugar.Errorf("Error updating XP for user: %v", err)
		return
	}

	// only check for rewards on level up
	oldLvl := sc.CalculateLevel(uc.XP)
	newLvl := sc.CalculateLevel(newXP)

	if oldLvl >= newLvl {
		return
	}

	if sc.LevelMessages == AllChannel || sc.LevelMessages == AllDM {
		err = bot.sendLevelMessage(s, m, sc, newLvl)
		if err != nil {
			bot.Sugar.Errorf("Error sending reward message: %v", err)
		}
	}

	reward := bot.getReward(m.GuildID, newLvl)
	if reward == nil {
		return
	}

	if !reward.RoleReward.IsValid() {
		return
	}

	// don't announce/log roles the user already has
	for _, r := range m.Member.RoleIDs {
		if r == reward.RoleReward {
			return
		}
	}

	err = s.AddRole(m.GuildID, m.Author.ID, reward.RoleReward, api.AddRoleData{
		AuditLogReason: api.AuditLogReason(fmt.Sprintf("Level reward for reaching level %v", newLvl)),
	})
	if err != nil {
		bot.Sugar.Errorf("Error adding role to user: %v", err)
		return
	}

	if sc.RewardLog.IsValid() {
		e := discord.Embed{
			Title:       "Level reward given",
			Description: fmt.Sprintf("%v reached level `%v`.", m.Author.Mention(), newLvl),
			Fields: []discord.EmbedField{
				{
					Name:  "Reward given",
					Value: reward.RoleReward.Mention(),
				},
				{
					Name:  "Message",
					Value: fmt.Sprintf("https://discord.com/channels/%v/%v/%v", m.GuildID, m.ChannelID, m.ID),
				},
			},
			Color: bcr.ColourBlurple,
		}

		s.SendEmbeds(sc.RewardLog, e)
	}

	if sc.LevelMessages == RewardsChannel {
		err = bot.sendLevelMessage(s, m, sc, newLvl)
		if err != nil {
			bot.Sugar.Errorf("Error sending reward message: %v", err)
		}
		return
	}

	if sc.LevelMessages != NoMessages && sc.RewardText != "" {
		var msgsDisabled bool
		bot.DB.Pool.QueryRow(context.Background(), "select disable_levelup_messages from user_config where user_id = $1", m.Author.ID).Scan(&msgsDisabled)

		if !msgsDisabled {
			txt := strings.NewReplacer("{lvl}", fmt.Sprint(newLvl)).Replace(sc.RewardText)

			ch, err := s.CreatePrivateChannel(m.Author.ID)
			if err == nil {
				s.SendMessage(ch.ID, txt)
			}
		}
	}
}

// this is only called if the level message setting isn't set to RewardsDM or NoMessages
func (bot *Bot) sendLevelMessage(s *state.State, m *gateway.MessageCreateEvent, sc Server, lvl int64) (err error) {
	if sc.LevelMessages == AllDM {
		var msgsDisabled bool
		bot.DB.Pool.QueryRow(context.Background(), "select disable_levelup_messages from user_config where user_id = $1", m.Author.ID).Scan(&msgsDisabled)
		if msgsDisabled {
			return
		}
	}

	e := discord.Embed{
		Thumbnail: &discord.EmbedThumbnail{
			URL: m.Author.AvatarURLWithType(discord.PNGImage) + "?size=256",
		},
		Title:       "Level up!",
		Description: fmt.Sprintf("%v has reached level %v!", m.Author.Mention(), lvl),
		Color:       bcr.ColourOrange,
		Footer: &discord.EmbedFooter{
			Text: m.Author.Tag(),
		},
		Timestamp: discord.NowTimestamp(),
	}

	chID := sc.LevelChannel
	if sc.LevelMessages == AllDM {
		ch, err := s.CreatePrivateChannel(m.Author.ID)
		// couldn't create a channel
		if err != nil {
			return nil
		}
		chID = ch.ID

		e.Description = fmt.Sprintf("Congratulations, you reached level %v!", lvl)
	}

	if !chID.IsValid() {
		return
	}

	_, err = s.SendEmbeds(chID, e)
	return
}
