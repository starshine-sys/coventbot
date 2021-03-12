package info

import (
	"fmt"
	"sort"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
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
	g, err := ctx.Session.Guild(ctx.Message.GuildID)
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
		roles       []string
		userPerms   discord.Permissions
		highestRole = "No roles"
	)
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
		Color:       discord.MemberColor(*g, *m),

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
