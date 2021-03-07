package static

import (
	"fmt"
	"sort"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/coventbot/etc"
)

func (bot *Bot) memberInfo(ctx *bcr.Context) (err error) {
	var m *discord.Member

	if len(ctx.Args) > 0 {
		m, err = ctx.ParseMember(ctx.RawArgs)
		if err != nil {
			return bot.userInfo(ctx)
		}
	} else {
		m, err = ctx.ParseMember(ctx.Author.ID.String())
	}
	if err != nil {
		_, err = ctx.Send(":x: User not found.", nil)
		return err
	}

	guildRoles, err := ctx.Session.Roles(ctx.Message.GuildID)
	if err != nil {
		_, err = ctx.Sendf("Internal error occurred:\n```%v```", err)
		return
	}

	// filter the roles to only the ones the user has
	var rls etc.Roles
	for _, gr := range guildRoles {
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
		colour      discord.Color
		roles       []string
		userPerms   discord.Permissions
		highestRole = "No roles"
	)
	for _, r := range rls {
		if r.Color != 0 {
			colour = r.Color
			break
		}
	}
	for _, r := range rls {
		userPerms |= r.Permissions
		roles = append(roles, r.Mention())
	}
	if len(rls) > 0 {
		highestRole = rls[0].Name
	}

	perms := strings.Join(bcr.PermStrings(userPerms), ", ")
	if len(perms) > 1000 {
		perms = perms[:1000] + "..."
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
				Name:  "Roles",
				Value: b.String(),
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
	return
}

// If tries to emulate a ternary operation as well as possible
func If(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}
