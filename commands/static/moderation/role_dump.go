package moderation

import (
	"fmt"
	"sort"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) roleDump(ctx *bcr.Context) (err error) {
	var (
		msgs []string
		b    strings.Builder
	)

	g, err := bot.State.Guild(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Sendf("An error occurred while getting the server info: %v", err)
		return
	}
	gm, err := bot.State.MembersAfter(ctx.Message.GuildID, 0, 0)
	if err != nil {
		_, err = ctx.Sendf("Couldn't get all the server's members: %v", err)
		return
	}

	sort.Slice(g.Roles, func(i, j int) bool {
		return g.Roles[i].Position > g.Roles[j].Position
	})

	for i, r := range g.Roles {
		members := 0
		for _, m := range gm {
			for _, id := range m.RoleIDs {
				if r.ID == id {
					members++
				}
			}
		}
		b.WriteString(fmt.Sprintf(" %v\n=== Mentionable: %v | Hoisted: %v | Members: %v", r.Name, r.Mentionable, r.Hoist, members))
		if b.Len() > 1900 {
			msgs = append(msgs, b.String())
			b.Reset()
		}
		// do this thrice, once for each type of perms
		m := strings.Join(bcr.PermStringsFor(bcr.MajorPerms, r.Permissions), ", ")
		if b.Len()+len(m) > 1900 {
			msgs = append(msgs, b.String())
			b.Reset()
		}
		if len(m) > 0 {
			b.WriteString("\n- " + m)
		}

		m = strings.Join(bcr.PermStringsFor(bcr.NotablePerms, r.Permissions), ", ")
		if b.Len()+len(m) > 1900 {
			msgs = append(msgs, b.String())
			b.Reset()
		}
		if len(m) > 0 {
			b.WriteString("\n+ " + m)
		}

		m = strings.Join(bcr.PermStringsFor(bcr.MinorPerms, r.Permissions), ", ")
		if b.Len()+len(m) > 1900 {
			msgs = append(msgs, b.String())
			b.Reset()
		}
		if len(m) > 0 {
			b.WriteString("\n= " + m)
		}

		if i != len(g.Roles)-1 {
			b.WriteString("\n\n\n=================\n")
		}
		if b.Len() > 1900 {
			msgs = append(msgs, b.String())
			b.Reset()
		}
	}
	if b.Len() > 0 {
		msgs = append(msgs, b.String())
	}

	for _, m := range msgs {
		ctx.Sendf("```diff\n%v\n```", m)
	}
	return
}

var (
	minorPerms = map[discord.Permissions]string{
		discord.PermissionEmbedLinks:         "Embed Links",
		discord.PermissionAttachFiles:        "Attach Files",
		discord.PermissionStream:             "Stream",
		discord.PermissionViewChannel:        "View Channel",
		discord.PermissionSendMessages:       "Send Messages",
		discord.PermissionAddReactions:       "Add Reactions",
		discord.PermissionConnect:            "Connect",
		discord.PermissionSpeak:              "Speak",
		discord.PermissionUseVAD:             "Use VAD",
		1 << 31:                              "Use Slash Commands",
		discord.PermissionReadMessageHistory: "Read Message History",
		discord.PermissionUseExternalEmojis:  "Use External Emojis",
		discord.PermissionChangeNickname:     "Change Nickname",
	}

	notablePerms = map[discord.Permissions]string{
		discord.PermissionCreateInstantInvite: "Create Instant Invite",
		discord.PermissionViewAuditLog:        "View Audit Log",
		discord.PermissionPrioritySpeaker:     "Priority Speaker",
		discord.PermissionSendTTSMessages:     "Send TTS Messages",
	}

	majorPerms = map[discord.Permissions]string{
		discord.PermissionManageNicknames: "Manage Nicknames",
		discord.PermissionManageRoles:     "Manage Roles",
		discord.PermissionManageWebhooks:  "Manage Webhooks",
		discord.PermissionManageEmojis:    "Manage Emojis",
		discord.PermissionMentionEveryone: "Mention Everyone",
		discord.PermissionMuteMembers:     "Mute Members",
		discord.PermissionDeafenMembers:   "Deafen Members",
		discord.PermissionMoveMembers:     "Move Members",
		discord.PermissionManageMessages:  "Manage Messages",
		discord.PermissionKickMembers:     "Kick Members",
		discord.PermissionBanMembers:      "Ban Members",
		discord.PermissionAdministrator:   "Administrator",
		discord.PermissionManageChannels:  "Manage Channels",
		discord.PermissionManageGuild:     "Manage Server",
	}
)
