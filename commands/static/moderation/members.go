package moderation

import (
	"fmt"
	"regexp"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) members(ctx *bcr.Context) (err error) {
	// Variables for flags
	var (
		flagRoles, flagAllRoles, flagNotRoles string

		nameContains, nameRegex string

		count, humans, bots bool
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

	fs.BoolVarP(&count, "count", "c", false, "Shows a member count instead of a member list")
	fs.BoolVarP(&humans, "humans", "h", false, "Shows only humans")
	fs.BoolVarP(&bots, "bots", "b", false, "Shows only bots")

	fs.Parse(ctx.Args)
	// send help if needed
	if b.Len() != 0 {
		b.Reset()
		usage := fmt.Sprintf("%vmembers ", ctx.Prefix)
		fs.VisitAll(func(f *flag.Flag) {
			usage += fmt.Sprintf("[-%v %v] ", f.Shorthand, f.Value.Type())

			b.WriteString(
				fmt.Sprintf("`--%v` (`-%v`): %v\n\n", f.Name, f.Shorthand, f.Usage),
			)
		})

		_, err = ctx.Send("", &discord.Embed{
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

	gm, err := bot.State.Session.MembersAfter(ctx.Message.GuildID, 0, 0)
	if err != nil {
		_, err = ctx.Sendf("Error: %v", err)
		return
	}

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
			_, err = ctx.Send("There was an error parsing the given regex.", nil)
			return err
		}

		gm = filterMembers(gm, func(m discord.Member) bool {
			return r.MatchString(m.User.Username + "#" + m.User.Discriminator)
		})
	}

	// send count if that flag is set
	if count {
		_, err = ctx.Send("", &discord.Embed{
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

	for i, m := range gm {
		b.WriteString(fmt.Sprintf("%v. %v#%v (%v)", i+1, m.User.Username, m.User.Discriminator, m.User.ID))
		if m.Nick != "" {
			b.WriteString(fmt.Sprintf("\n%v  - %v", strings.Repeat(" ", len(fmt.Sprint(i+1))), m.Nick))
		}
		b.WriteString("\n\n")
	}
	if len(gm) == 0 {
		b.WriteString("No results.")
	}

	// if length is longer than 2k, we need to send a text file instead
	if b.Len() >= 2000 {
		_, err = ctx.NewMessage().AddFile("members.txt", strings.NewReader(b.String())).Send()
		if err == bcr.ErrBotMissingPermissions {
			_, err = ctx.Send(":x: I can't attach files in this channel.", nil)
			return
		}
		return
	}

	_, err = ctx.Send("", &discord.Embed{
		Title:       "Members in Query",
		Description: "```py\n" + b.String() + "\n```",
		Color:       ctx.Router.EmbedColor,
		Fields: []discord.EmbedField{{
			Name:  "Query",
			Value: "```" + ctx.RawArgs + "```",
		}},
	})
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
