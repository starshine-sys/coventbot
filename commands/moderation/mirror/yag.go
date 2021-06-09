package mirror

import (
	"regexp"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/starshine-sys/tribble/commands/moderation/modlog"
)

// regex for getting the moderator
var yagModRegexp = regexp.MustCompile(`\(ID (\d{16,})\)`)

// regex for getting the duration
var yagDurationRegexp = regexp.MustCompile(`Duration: ([\w\s]+)`)

// regexes for figuring out what the action type is + capturing ID and reason
var yagWarnRegexp = regexp.MustCompile(`\*\*⚠Warned [^#]+\*\*#\d{4} \*\(ID (\d{16,})\)\*\n📄\*\*Reason:\*\* (.*)`)
var yagMuteRegexp = regexp.MustCompile(`\*\*🔇Muted [^#]+\*\*#\d{4} \*\(ID (\d{16,})\)\*\n📄\*\*Reason:\*\* (.*) \(\[`)
var yagUnmuteRegexp = regexp.MustCompile(`\*\*🔊Unmuted [^#]+\*\*#\d{4} \*\(ID (\d{16,})\)\*\n📄\*\*Reason:\*\* (.*)`)

func (bot *Bot) processYAG(m *gateway.MessageCreateEvent) {
	// no embeds so we don't care
	if len(m.Embeds) == 0 {
		return
	}

	entry := bot.parseYAGEmbed(m.GuildID, m.Embeds[0])
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

func (bot *Bot) parseYAGEmbed(g discord.GuildID, e discord.Embed) *modlog.Entry {
	if yagWarnRegexp.MatchString(e.Description) {
		return bot.warnYAGEmbed(g, e)
	}

	if yagMuteRegexp.MatchString(e.Description) {
		return bot.muteYAGEmbed(g, e)
	}

	if yagUnmuteRegexp.MatchString(e.Description) {
		return bot.unmuteYAGEmbed(g, e)
	}

	return nil
}

func (bot *Bot) warnYAGEmbed(g discord.GuildID, e discord.Embed) *modlog.Entry {
	groups := yagWarnRegexp.FindStringSubmatch(e.Description)
	if len(groups) < 3 {
		return nil
	}

	user, _ := discord.ParseSnowflake(groups[1])
	reason := groups[2]

	groups = yagModRegexp.FindStringSubmatch(e.Author.Name)
	mod, _ := discord.ParseSnowflake(groups[1])

	return &modlog.Entry{
		ServerID:   g,
		UserID:     discord.UserID(user),
		ModID:      discord.UserID(mod),
		ActionType: "warn",
		Reason:     reason,
	}
}

func (bot *Bot) muteYAGEmbed(g discord.GuildID, e discord.Embed) *modlog.Entry {
	groups := yagMuteRegexp.FindStringSubmatch(e.Description)
	if len(groups) < 3 {
		return nil
	}

	user, _ := discord.ParseSnowflake(groups[1])
	reason := groups[2]

	groups = yagModRegexp.FindStringSubmatch(e.Author.Name)
	mod, _ := discord.ParseSnowflake(groups[1])

	if e.Footer != nil {
		if yagDurationRegexp.MatchString(e.Footer.Text) {
			groups = yagDurationRegexp.FindStringSubmatch(e.Footer.Text)
			if groups[1] != "" && groups[1] != "permanent" {
				reason += " (duration: " + groups[1] + ")"
			}
		}
	}

	return &modlog.Entry{
		ServerID:   g,
		UserID:     discord.UserID(user),
		ModID:      discord.UserID(mod),
		ActionType: "mute",
		Reason:     reason,
	}
}

func (bot *Bot) unmuteYAGEmbed(g discord.GuildID, e discord.Embed) *modlog.Entry {
	groups := yagUnmuteRegexp.FindStringSubmatch(e.Description)
	if len(groups) < 3 {
		return nil
	}

	user, _ := discord.ParseSnowflake(groups[1])
	reason := groups[2]

	groups = yagModRegexp.FindStringSubmatch(e.Author.Name)
	mod, _ := discord.ParseSnowflake(groups[1])

	return &modlog.Entry{
		ServerID:   g,
		UserID:     discord.UserID(user),
		ModID:      discord.UserID(mod),
		ActionType: "unmute",
		Reason:     reason,
	}
}
