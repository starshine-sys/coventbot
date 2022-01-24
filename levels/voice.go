package levels

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/state/store"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) voiceLevelsLoop() {
	t := time.NewTicker(5 * time.Minute)

	for {
		<-t.C

		bot.Sugar.Debugf("Checking voice states for XP")
		bot.voiceLevels()
	}
}

func (bot *Bot) voiceLevels() {
	var guilds []Server
	err := pgxscan.Select(context.Background(), bot.DB.Pool, &guilds, "select * from server_levels where voice = true")
	if err != nil {
		bot.Sugar.Errorf("Error fetching guilds with voice XP enabled: %v", err)
		return
	}

	for _, gc := range guilds {
		s, _ := bot.Router.StateFromGuildID(gc.ID)

		vcs, err := s.VoiceStates(gc.ID)
		if err != nil {
			if err == store.ErrNotFound {
				continue
			}
			bot.Sugar.Errorf("Error getting voice states for %v: %v", gc.ID, err)
			continue
		}

		for _, vc := range vcs {
			err = bot.checkVoiceState(s, gc, vc)
			if err != nil {
				bot.Sugar.Errorf("Error processing voice state for u:%v/g:%v: %v", vc.UserID, gc.ID, err)
			}
		}
	}
}

func (bot *Bot) checkVoiceState(s *state.State, gc Server, vc discord.VoiceState) (err error) {
	// check blocked channels/categories
	for _, blocked := range gc.BlockedChannels {
		if vc.ChannelID == discord.ChannelID(blocked) {
			return
		}
	}
	if ch, err := s.Channel(vc.ChannelID); err == nil {
		for _, blocked := range gc.BlockedCategories {
			if ch.ParentID == discord.ChannelID(blocked) {
				return nil
			}
		}
	}

	if vc.Mute || vc.Deaf || vc.SelfDeaf || vc.SelfMute {
		bot.Sugar.Debugf("User %v is muted or deafened, skipping", vc.UserID)
		return
	}

	if bot.isBlacklisted(gc.ID, vc.UserID) {
		bot.Sugar.Debugf("User %v is noleveled, skipping", vc.UserID)
		return
	}

	if vc.Member == nil {
		bot.Sugar.Debugf("Member for %v is nil, fetching member", vc.UserID)
		m, err := bot.Member(gc.ID, vc.UserID)
		if err != nil {
			bot.Sugar.Errorf("Error getting member %v (g:%v): %v", vc.UserID, gc.ID, err)
		}
		vc.Member = &m
	}

	// check blocked roles
	for _, blocked := range gc.BlockedRoles {
		for _, r := range vc.Member.RoleIDs {
			if discord.RoleID(blocked) == r {
				return
			}
		}
	}

	uc, err := bot.getUser(gc.ID, vc.UserID)
	if err != nil {
		bot.Sugar.Errorf("Error getting user: %v", err)
		return
	}

	if uc.NextTime.Add(gc.BetweenXP).After(time.Now().UTC()) {
		bot.Sugar.Debugf("User %v typed within the last %v, not incrementing voice XP", vc.UserID, gc.BetweenXP)
		return nil
	}

	// increment the user's xp!
	bot.Sugar.Debugf("Incrementing voice XP for user %v", vc.UserID)
	newXP, err := bot.incrementXP(gc.ID, vc.UserID, gc.CarlineCompatible)
	if err != nil {
		bot.Sugar.Errorf("Error updating XP for user: %v", err)
		return
	}

	// only check for rewards on level up
	oldLvl := gc.CalculateLevel(uc.XP)
	newLvl := gc.CalculateLevel(newXP)

	if oldLvl >= newLvl {
		return
	}

	reward := bot.getReward(gc.ID, newLvl)
	if reward == nil {
		return
	}

	if !reward.RoleReward.IsValid() {
		return
	}

	// don't announce/log roles the user already has
	for _, r := range vc.Member.RoleIDs {
		if r == reward.RoleReward {
			return
		}
	}

	err = s.AddRole(gc.ID, vc.UserID, reward.RoleReward, api.AddRoleData{
		AuditLogReason: api.AuditLogReason(fmt.Sprintf("Level reward for reaching level %v", newLvl)),
	})
	if err != nil {
		bot.Sugar.Errorf("Error adding role to user: %v", err)
		return
	}

	if gc.RewardLog.IsValid() {
		e := discord.Embed{
			Title:       "Level reward given",
			Description: fmt.Sprintf("%v reached level `%v`.", vc.UserID.Mention(), newLvl),
			Fields: []discord.EmbedField{
				{
					Name:  "Reward given",
					Value: reward.RoleReward.Mention(),
				},
				{
					Name:  "Message",
					Value: "*N/A, voice XP in* " + vc.ChannelID.Mention(),
				},
			},
			Color: bcr.ColourBlurple,
		}

		_, err = s.SendEmbeds(gc.RewardLog, e)
		if err != nil {
			return err
		}
	}

	if gc.LevelMessages != NoMessages && gc.RewardText != "" {
		var msgsDisabled bool
		bot.DB.Pool.QueryRow(context.Background(), "select disable_levelup_messages from user_config where user_id = $1", vc.UserID).Scan(&msgsDisabled)

		if !msgsDisabled {
			txt := strings.NewReplacer("{lvl}", fmt.Sprint(newLvl)).Replace(gc.RewardText)

			ch, err := s.CreatePrivateChannel(vc.UserID)
			if err == nil {
				_, err = s.SendMessage(ch.ID, txt)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
