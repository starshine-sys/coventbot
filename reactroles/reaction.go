package reactroles

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

func (bot *Bot) reactionAdd(ev *gateway.MessageReactionAddEvent) {
	if !ev.GuildID.IsValid() || ev.UserID == bot.Router.Bot.ID {
		return
	}

	entries, err := bot.getEntries(ev.MessageID)
	if err != nil {
		bot.Sugar.Errorf("Error: %v", err)
		return
	}

	var toAdd []discord.RoleID
	for _, e := range entries {
		if ev.Emoji.ID.String() == e.Emote || ev.Emoji.Name == e.Emote {
			toAdd = append(toAdd, e.RoleID)
			break
		}
	}

	for _, r := range toAdd {
		err = bot.State.AddRole(ev.GuildID, ev.UserID, r)
		if err != nil {
			bot.Sugar.Errorf("Error adding role %v to user: %v", r, err)
		}
	}
}

func (bot *Bot) reactionRemove(ev *gateway.MessageReactionRemoveEvent) {
	if !ev.GuildID.IsValid() || ev.UserID == bot.Router.Bot.ID {
		return
	}

	entries, err := bot.getEntries(ev.MessageID)
	if err != nil {
		bot.Sugar.Errorf("Error: %v", err)
		return
	}

	var toRemove []discord.RoleID
	for _, e := range entries {
		if ev.Emoji.ID.String() == e.Emote || ev.Emoji.Name == e.Emote {
			toRemove = append(toRemove, e.RoleID)
			break
		}
	}

	for _, r := range toRemove {
		err = bot.State.RemoveRole(ev.GuildID, ev.UserID, r)
		if err != nil {
			bot.Sugar.Errorf("Error removing role %v from user: %v", r, err)
		}
	}
}
