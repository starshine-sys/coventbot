package moderation

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) members(ctx *bcr.Context) (err error) {
	// Variables for flags
	var (
		flagRoles, flagAllRoles, flagNotRoles string

		nameContains, nameRegex string

		count, humans, bots bool

		format     string
		enkoFormat bool
		code       bool

		file bool
	)

	var b strings.Builder

	fs := flag.NewFlagSet(fmt.Sprintf("%vmembers", ctx.Prefix), flag.ContinueOnError)
	// set output for help
	fs.SetOutput(&b)

	fs.StringVarP(&flagRoles, "roles", "r", "", "Shows members with at least one of these roles (comma-separated)")
	fs.StringVarP(&flagAllRoles, "all-roles", "a", "", "Shows members with all of these roles (comma-separated)")
	fs.StringVarP(&flagNotRoles, "not-roles", "n", "", "Shows members with none of these roles (comma-separated)")
	fs.StringVarP(&nameContains, "name-contains", "C", "", "Shows members whose names contain the given text")
	fs.StringVarP(&nameRegex, "name-regex", "R", "", "Shows members whose names match the given regex (use `(?i)` at the beginning for case-insensitive matching)")

	fs.StringVarP(&format, "format", "f", "%in. %u#%d (%id)", `Format to use for the member list.
Supported options:
- %in: index
- %id: user ID
- %u: username
- %d: discriminator
- %n: nickname
- %m: mention
- %cd: creation date
- %jd: join date`)

	fs.BoolVarP(&count, "count", "c", false, "Shows a member count instead of a member list")
	fs.BoolVarP(&humans, "humans", "h", false, "Shows only humans")
	fs.BoolVarP(&bots, "bots", "b", false, "Shows only bots")
	fs.BoolVarP(&enkoFormat, "enko-format", "E", false, "Use EnkoMojishia's format for the member list.")
	fs.BoolVarP(&code, "code", "B", false, "Output as a code block (only in embeds).")
	fs.BoolVarP(&file, "file", "F", false, "Output to a file instead of a paginated embed.")

	fs.Parse(ctx.Args)
	// send help if needed
	if b.Len() != 0 {
		b.Reset()
		usage := fmt.Sprintf("%vmembers ", ctx.Prefix)
		fs.VisitAll(func(f *flag.Flag) {
			usage += fmt.Sprintf("[-%v %v] ", f.Shorthand, f.Value.Type())

			b.WriteString(
				fmt.Sprintf("`-%v, --%v`: %v\n", f.Shorthand, f.Name, f.Usage),
			)
		})

		_, err = ctx.Send("", discord.Embed{
			Title:       "`MEMBERS`",
			Description: "```" + usage + "```",
			Fields: []discord.EmbedField{{
				Name:  "Flags",
				Value: b.String(),
			}},
			Color: ctx.Router.EmbedColor,
		})
		return
	}

	// we don't need ctx.RawArgs, so just set it to "None" if it's empty
	if ctx.RawArgs == "" {
		ctx.RawArgs = "None"
	}

	if enkoFormat {
		format = "%in. %u#%d (%id)\n​  - %n"
	}

	gm := bot.Members(ctx.Message.GuildID)

	// filter stuff
	if humans {
		gm = filterMembers(gm, func(m discord.Member) bool {
			return !m.User.Bot
		})
	}
	if bots {
		gm = filterMembers(gm, func(m discord.Member) bool {
			return m.User.Bot
		})
	}

	// filter single roles
	if flagRoles != "" {
		roles, _ := ctx.GreedyRoleParser(strings.Split(flagRoles, ","))
		var ids []discord.RoleID
		for _, r := range roles {
			ids = append(ids, r.ID)
		}
		gm = filterMembers(gm, func(m discord.Member) bool {
			for _, r := range ids {
				for _, mr := range m.RoleIDs {
					if r == mr {
						return true
					}
				}
			}
			return false
		})
	}

	if flagNotRoles != "" {
		roles, _ := ctx.GreedyRoleParser(strings.Split(flagNotRoles, ","))
		var ids []discord.RoleID
		for _, r := range roles {
			ids = append(ids, r.ID)
		}
		gm = filterMembers(gm, func(m discord.Member) bool {
			for _, r := range ids {
				for _, mr := range m.RoleIDs {
					if r == mr {
						return false
					}
				}
			}
			return true
		})
	}

	if flagAllRoles != "" {
		roles, _ := ctx.GreedyRoleParser(strings.Split(flagAllRoles, ","))
		var ids []discord.RoleID
		for _, r := range roles {
			ids = append(ids, r.ID)
		}

		gm = filterMembers(gm, func(m discord.Member) bool {
			i := 0
			for _, r := range ids {
				for _, mr := range m.RoleIDs {
					if r == mr {
						i++
					}
				}
			}
			return i >= len(ids)
		})
	}

	// filter names
	if nameContains != "" {
		gm = filterMembers(gm, func(m discord.Member) bool {
			return strings.Contains(strings.ToLower(m.User.Username+"#"+m.User.Discriminator), strings.ToLower(nameContains))
		})
	}

	if nameRegex != "" {
		r, err := regexp.Compile(nameRegex)
		if err != nil {
			_, err = ctx.Send("There was an error parsing the given regex.")
			return err
		}

		gm = filterMembers(gm, func(m discord.Member) bool {
			return r.MatchString(m.User.Username + "#" + m.User.Discriminator)
		})
	}

	// send count if that flag is set
	if count {
		_, err = ctx.Send("", discord.Embed{
			Title:       "Members in Query",
			Description: fmt.Sprintf("Member count is: %v", len(gm)),
			Color:       ctx.Router.EmbedColor,
			Fields: []discord.EmbedField{{
				Name:  "Query",
				Value: "```" + ctx.RawArgs + "```",
			}},
		})
		return
	}

	var members []string

	for i, m := range gm {
		nick := m.Nick
		if nick == "" {
			nick = "no nickname"
		}

		if !enkoFormat {
			members = append(members, strings.NewReplacer(
				"%in", fmt.Sprint(i+1),
				"%id", m.User.ID.String(),
				"%u", m.User.Username,
				"%d", m.User.Discriminator,
				"%n", nick,
				"%m", m.Mention(),
				"%cd", m.User.ID.Time().UTC().Format(time.RFC1123),
				"%jd", m.Joined.Time().UTC().Format(time.RFC1123),
			).Replace(format+"\n"))
		} else {
			s := fmt.Sprintf("%v. %v#%v (%v)\n", i+1, m.User.Username, m.User.Discriminator, m.User.ID)
			if m.Nick != "" {
				s += fmt.Sprintf("  - %v\n", m.Nick)
			}
			members = append(members, s)
		}
	}
	if len(gm) == 0 {
		members = append(members, "No results.")
	}

	if file {
		_, err = ctx.NewMessage().AddFile("members.txt", strings.NewReader(strings.Join(members, ""))).Send()
		if err == bcr.ErrBotMissingPermissions {
			_, err = ctx.Send("I can't attach files in this channel.")
			return
		}
		return
	}

	var embeds []discord.Embed

	{
		var count int
		for _, m := range members {
			count++
			if b.Len()+len(m) > 2000 || count > 25 {
				var desc string
				if enkoFormat || code {
					desc = "```py\n" + b.String() + "\n```"
				} else {
					desc = b.String()
				}

				embeds = append(embeds, discord.Embed{
					Title:       "Members",
					Description: desc,
					Fields: []discord.EmbedField{{
						Name:  "Query",
						Value: "```" + ctx.RawArgs + "```",
					}},
					Color: ctx.Router.EmbedColor,
				})

				b.Reset()
				count = 0
			}

			b.WriteString(m)
		}
	}

	var desc string
	if enkoFormat || code {
		desc = "```py\n" + b.String() + "\n```"
	} else {
		desc = b.String()
	}

	embeds = append(embeds, discord.Embed{
		Title:       "Members",
		Description: desc,
		Fields: []discord.EmbedField{{
			Name:  "Query",
			Value: "```" + ctx.RawArgs + "```",
		}},
		Color: ctx.Router.EmbedColor,
	})

	_, err = bot.PagedEmbed(ctx, embeds, 10*time.Minute)
	return
}

// filterMembers is a helper function for filtering a slice of members
func filterMembers(in []discord.Member, filter func(discord.Member) bool) (out []discord.Member) {
	for _, m := range in {
		if filter(m) {
			out = append(out, m)
		}
	}
	return
}
