// SPDX-License-Identifier: AGPL-3.0-only
package info

import (
	"fmt"
	"sort"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	bcr1 "github.com/starshine-sys/bcr"
	"github.com/starshine-sys/bcr/v2"
	"github.com/starshine-sys/tribble/common"
)

func (bot *Bot) slashMemberInfo(ctx *bcr.CommandContext) (err error) {
	m := ctx.Member

	sf, err := ctx.Option("user").SnowflakeValue()
	if err == nil && sf.IsValid() {
		m, err = ctx.State.Member(ctx.Guild.ID, discord.UserID(sf))
		if err != nil {
			return bot.slashUserInfo(ctx)
		}
	}

	// filter the roles to only the ones the user has
	var rls bcr1.Roles
	for _, gr := range ctx.Guild.Roles {
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
		userPerms   discord.Permissions
		highestRole = "No roles"
	)
	for _, r := range rls {
		userPerms |= r.Permissions
	}
	if len(rls) > 0 {
		highestRole = rls[0].Name
	}

	var perms []string
	if ctx.Guild.OwnerID == m.User.ID {
		perms = append(perms, "Server Owner")
		userPerms = userPerms.Add(discord.PermissionAll)
	}
	perms = append(perms, bcr.PermStringsFor(bcr.MajorPerms, userPerms)...)

	permString := strings.Join(perms, ", ")
	if len(permString) > 1000 {
		permString = permString[:1000] + "..."
	} else if permString == "" {
		permString = "No special permissions"
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

	colour, _ := discord.MemberColor(*ctx.Guild, *m)
	if colour == 0 {
		colour = m.User.Accent
		if colour == 0 {
			colour = bot.EmbedColour
		}
	}

	avatarURL := m.User.AvatarURL()
	if m.Avatar != "" {
		avatarURL = m.AvatarURL(ctx.Guild.ID)
	}

	e := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: m.User.Tag(),
			Icon: m.User.AvatarURL(),
		},
		Thumbnail: &discord.EmbedThumbnail{
			URL: avatarURL,
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
				Name:   "Username",
				Value:  m.User.Tag(),
				Inline: true,
			},
		},

		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v", m.User.ID),
		},
		Timestamp: discord.NowTimestamp(),
	}

	if m.User.DisplayName != "" {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:   "Display name",
			Value:  m.User.DisplayName,
			Inline: true,
		})
	}

	e.Fields = append(e.Fields, discord.EmbedField{
		Name: "Created at",
		Value: fmt.Sprintf("<t:%v:D> <t:%v:T>\n(%v)",
			m.User.ID.Time().Unix(), m.User.ID.Time().Unix(),
			common.FormatTime(m.User.ID.Time().UTC()),
		),
	})

	if m.Avatar != "" {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:   "Server avatar",
			Value:  fmt.Sprintf("[Link](%v?size=1024)", m.AvatarURL(ctx.Guild.ID)),
			Inline: true,
		})
	}

	e.Fields = append(e.Fields, []discord.EmbedField{
		{
			Name:   "Nickname",
			Value:  fmt.Sprintf("%v", FirstNonZero(m.Nick, m.User.Username)),
			Inline: true,
		},
		{
			Name:   "Highest role",
			Value:  highestRole,
			Inline: true,
		},
		{
			Name: "Joined at",
			Value: fmt.Sprintf("<t:%v:D> <t:%v:T>\n(%v)\n%v days after the server was created",
				m.Joined.Time().Unix(), m.Joined.Time().Unix(),
				common.FormatTime(m.Joined.Time().UTC()),
				int(
					m.Joined.Time().Sub(ctx.Guild.ID.Time()).Hours()/24,
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
	}...)

	if u, err := ctx.State.User(m.User.ID); err == nil {
		if u.Banner != "" {
			e.Image = &discord.EmbedImage{
				URL: u.BannerURL() + "?size=1024",
			}
		}
	}

	return ctx.Reply("", e)
}

func (bot *Bot) slashUserInfo(ctx *bcr.CommandContext) (err error) {
	sf, err := ctx.Option("user").SnowflakeValue()
	if err != nil {
		return ctx.ReplyEphemeral("User not found.")
	}

	u, err := ctx.State.User(discord.UserID(sf))
	if err != nil {
		return ctx.ReplyEphemeral("User not found.")
	}

	e := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: u.Tag(),
			Icon: u.AvatarURL(),
		},
		Thumbnail: &discord.EmbedThumbnail{
			URL: u.AvatarURL(),
		},
		Description: u.ID.String(),
		Color:       bot.EmbedColour,

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
				Value:  u.Tag(),
				Inline: true,
			},
		},

		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("ID: %v", u.ID),
		},
		Timestamp: discord.NowTimestamp(),
	}

	if u.DisplayName != "" {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:   "Display name",
			Value:  u.DisplayName,
			Inline: true,
		})
	}

	e.Fields = append(e.Fields, discord.EmbedField{
		Name: "Created at",
		Value: fmt.Sprintf("<t:%v:D> <t:%v:T>\n(%v)",
			u.ID.Time().Unix(), u.ID.Time().Unix(),
			common.FormatTime(u.ID.Time().UTC()),
		),
	})

	if u.Accent != 0 {
		e.Color = u.Accent
	}

	if u.Banner != "" {
		e.Image = &discord.EmbedImage{
			URL: u.BannerURL() + "?size=1024",
		}
	}

	return ctx.Reply("", e)
}
