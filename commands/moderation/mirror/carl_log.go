package mirror

import (
	"errors"
	"regexp"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/starshine-sys/tribble/commands/moderation/modlog"
)

var (
	carlDurationRegexp = regexp.MustCompile(`\*\*Duration:\*\* (.*)`)
	carlReasonRegexp   = regexp.MustCompile(`\*\*Reason:\*\* ([\s\S]*)\n\*\*Responsible moderator:\*\*`)
	carlModRegexp      = regexp.MustCompile(`\*\*Responsible moderator:\*\* (.*)`)
)

func (bot *Bot) processCarlLog(m *gateway.MessageCreateEvent) {
	entry := bot.parseCarlEmbed(m.GuildID, m.Embeds[0])
	if entry == nil {
		return
	}

	entry, err := bot.ModLog.InsertEntry(entry.ServerID, entry.UserID, entry.ModID, m.ID.Time().UTC(), entry.ActionType, entry.Reason)
	if err != nil {
		bot.Sugar.Errorf("Error adding mod log entry: %v", err)
	} else {
		bot.Sugar.Debugf("Added mod log entry %v for %v", entry.ID, entry.ServerID)
	}
}

func (bot *Bot) parseCarlEmbed(id discord.GuildID, e discord.Embed) *modlog.Entry {
	if e.Footer == nil {
		return nil
	}

	if !strings.HasPrefix(e.Footer.Text, "ID: ") {
		return nil
	}

	entry := &modlog.Entry{
		ServerID: id,
	}

	userID, err := discord.ParseSnowflake(strings.TrimSpace(strings.TrimPrefix(e.Footer.Text, "ID: ")))
	if err != nil {
		bot.Sugar.Errorf("Error parsing user ID from \"%v\": %v", e.Footer.Text, err)
	}

	entry.UserID = discord.UserID(userID)

	groups := carlModRegexp.FindStringSubmatch(e.Description)
	if len(groups) < 2 {
		bot.Sugar.Infof("Couldn't match moderator from (probably) mod log embed")
		return nil
	}

	bot.modIDMapMu.Lock()
	if v, ok := bot.modIDMap[groups[1]]; ok {
		entry.ModID = v
	} else {
		mod, err := bot.ParseMember(id, groups[1])
		if err != nil {
			bot.Sugar.Infof("Couldn't find member named %v", groups[1])
			bot.modIDMapMu.Unlock()
			return nil
		}
		entry.ModID = mod.User.ID
		bot.modIDMap[groups[1]] = mod.User.ID
	}
	bot.modIDMapMu.Unlock()

	groups = carlReasonRegexp.FindStringSubmatch(e.Description)
	if len(groups) < 2 {
		entry.Reason = "No reason given"
	} else {
		if strings.Contains(groups[1], "No reason given,") {
			entry.Reason = "No reason given"
		} else {
			entry.Reason = groups[1]
		}
	}

	groups = carlDurationRegexp.FindStringSubmatch(e.Description)
	if len(groups) >= 2 {
		entry.Reason += " (duration: " + groups[1] + ")"
	}

	switch {
	case strings.Contains(e.Title, "warn"):
		entry.ActionType = "warn"
	case strings.Contains(e.Title, "unmute"):
		entry.ActionType = "unmute"
	case strings.Contains(e.Title, "mute"):
		entry.ActionType = "mute"
	case strings.Contains(e.Title, "ban"):
		entry.ActionType = "ban"
	case strings.Contains(e.Title, "kick"):
		entry.ActionType = "kick"
	default:
		// unsupported type
		return nil
	}

	return entry
}

// ParseMember parses a member name, because Carl doesn't show IDs for some reason
func (bot *Bot) ParseMember(guildID discord.GuildID, s string) (c *discord.Member, err error) {
	members, err := bot.State.Members(guildID)
	if err != nil {
		return nil, err
	}

	for _, m := range members {
		if m.User.Username+"#"+m.User.Discriminator == s {
			return &m, nil
		}
	}

	return nil, errors.New("member not found")
}
