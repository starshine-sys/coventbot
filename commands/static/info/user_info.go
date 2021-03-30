package info

import (
	"fmt"
	"sort"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/commands/static/roles"
	"github.com/starshine-sys/tribble/etc"
)

func (bot *Bot) memberInfo(ctx *bcr.Context) (err error) {
	m := ctx.Member
	m.User = ctx.Author

	if len(ctx.Args) > 0 {
		m, err = ctx.ParseMember(ctx.RawArgs)
		if err != nil {
			return bot.userInfo(ctx)
		}
	}

	// get guild
	g, err := ctx.State.Guild(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Sendf("Internal error occurred:\n```%v```", err)
		return
	}

	// filter the roles to only the ones the user has
	var rls etc.Roles
	for _, gr := range g.Roles {
		for _, ur := range m.RoleIDs {
			if gr.ID == ur {
				rls = append(rls, gr)
			}
		}
	}
	sort.Sort(rls)

	// get global info
	// can't do this with a single loop because the loop for colour has to break the moment it's found one
	var (
		userRoles   []string
		userPerms   discord.Permissions
		highestRole = "No roles"
	)
	for _, r := range rls {
		userPerms |= r.Permissions
		userRoles = append(userRoles, r.Mention())
	}
	if len(rls) > 0 {
		highestRole = rls[0].Name
	}

	var perms []string
	if g.OwnerID == m.User.ID {
		perms = append(perms, "Server Owner")
		userPerms.Add(discord.PermissionAll)
	}
	perms = append(perms, roles.PermStringsFor(roles.MajorPerms, userPerms)...)

	permString := strings.Join(perms, ", ")
	if len(permString) > 1000 {
		permString = permString[:1000] + "..."
	}
	var b strings.Builder
	for i, r := range rls {
		if b.Len() > 900 {
			b.WriteString(fmt.Sprintf("\n```Too many roles to list (showing %v/%v)```", i, len(rls)))
			break
		}
		b.WriteString(r.Mention())
		if i != len(rls)-1 {
			b.WriteString(", ")
		}
	}
	if b.Len() == 0 {
		b.WriteString("No roles.")
	}

	colour := discord.MemberColor(*g, *m)
	if colour == 0 {
		colour = ctx.Router.EmbedColor
	}

	e := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: m.User.Username + "#" + m.User.Discriminator,
			Icon: m.User.AvatarURL(),
		},
		Thumbnail: &discord.EmbedThumbnail{
			URL: m.User.AvatarURL(),
		},
		Description: m.User.ID.String(),
		Color:       colour,

		Fields: []discord.EmbedField{
			{
				Name:  "User information for",
				Value: m.Mention(),
			},
			{
				Name:   "Avatar",
				Value:  fmt.Sprintf("[Link](%v?size=1024)", m.User.AvatarURL()),
				Inline: true,
			},
			{
				Name:   "Highest role",
				Value:  highestRole,
				Inline: true,
			},
			{
				Name: "Created at",
				Value: fmt.Sprintf("%v\n(%v)",
					m.User.ID.Time().UTC().Format("Jan _2 2006, 15:04:05 MST"),
					etc.HumanizeTime(etc.DurationPrecisionMinutes, m.User.ID.Time().UTC()),
				),
			},
			{
				Name:   "Username",
				Value:  m.User.Username + "#" + m.User.Discriminator,
				Inline: true,
			},
			{
				Name:   "Nickname",
				Value:  fmt.Sprintf("%v", If(m.Nick != "", m.Nick, m.User.Username)),
				Inline: true,
			},
			{
				Name: "Joined at",
				Value: fmt.Sprintf("%v\n(%v)\n%v days after the server was created",
					m.Joined.Time().UTC().Format("Jan _2 2006, 15:04:05 MST"),
					etc.HumanizeTime(etc.DurationPrecisionMinutes, m.Joined.Time().UTC()),
					int(
						m.Joined.Time().Sub(ctx.Message.GuildID.Time()).Hours()/24,
					),
				),
			},
			{
				Name:  fmt.Sprintf("Roles (%v)", len(rls)),
				Value: b.String(),
			},
			{
				Name:  "Key permissions",
				Value: permString,
			},
		},

		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v", m.User.ID),
		},
		Timestamp: discord.NowTimestamp(),
	}

	_, err = ctx.Send("", &e)
	return
}

func (bot *Bot) userInfo(ctx *bcr.Context) (err error) {
	u, err := ctx.ParseUser(ctx.RawArgs)
	if err != nil {
		_, err = ctx.Send("User not found.", nil)
		return
	}

	e := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: u.Username + "#" + u.Discriminator,
			Icon: u.AvatarURL(),
		},
		Thumbnail: &discord.EmbedThumbnail{
			URL: u.AvatarURL(),
		},
		Description: u.ID.String(),
		Color:       ctx.Router.EmbedColor,

		Fields: []discord.EmbedField{
			{
				Name:  "User information for",
				Value: u.Mention(),
			},
			{
				Name:   "Avatar",
				Value:  fmt.Sprintf("[Link](%v?size=1024)", u.AvatarURL()),
				Inline: true,
			},
			{
				Name:   "Username",
				Value:  u.Username + "#" + u.Discriminator,
				Inline: true,
			},
			{
				Name: "Created at",
				Value: fmt.Sprintf("%v\n(%v)",
					u.ID.Time().UTC().Format("Jan _2 2006, 15:04:05 MST"),
					etc.HumanizeTime(etc.DurationPrecisionMinutes, u.ID.Time().UTC()),
				),
			},
		},

		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v", u.ID),
		},
		Timestamp: discord.NowTimestamp(),
	}

	_, err = ctx.Send("", &e)
	return
}

// If tries to emulate a ternary operation as well as possible
func If(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}

// admin isn't in here as we do that one manually, to make sure it shows up at the very beginning
var majorPerms = map[discord.Permissions]string{
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
	discord.PermissionManageChannels:  "Manage Channels",
	discord.PermissionManageGuild:     "Manage Server",
}
