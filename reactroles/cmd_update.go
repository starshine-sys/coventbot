package reactroles

import (
	"context"
	"errors"
	"regexp"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

var emojiMatch = regexp.MustCompile("<(?P<animated>a)?:(?P<name>\\w+):(?P<emoteID>\\d{15,})>")
var idRe = regexp.MustCompile(`^\d{15,}$`)

func (bot *Bot) update(ctx *bcr.Context) (err error) {
	msg, err := ctx.ParseMessage(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("Couldn't parse a message from your input.")
		return
	}

	if msg.GuildID != ctx.Message.GuildID {
		_, err = ctx.Send("The message you gave isn't in this server.")
		return
	}

	roles, err := bot.parseRoles(ctx, ctx.Args[1:])
	if err != nil {
		if err == errNoPairs {
			_, err = ctx.Send("You must give emoji-role *pairs*.")
			return
		}
		_, err = ctx.Send("Couldn't parse one or more roles.")
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), `insert into react_roles
	(server_id, channel_id, message_id) values ($1, $2, $3) on conflict (message_id) do nothing`, msg.GuildID, msg.ChannelID, msg.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	for _, r := range roles {
		_, err = bot.DB.Pool.Exec(context.Background(), `insert into react_role_entries
		(message_id, emote, role_id) values ($1, $2, $3) on conflict (message_id, emote) do update
		set role_id = $3`, msg.ID, r.Emote, r.Role.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		emoji := discord.APIEmoji(r.Emote)
		if r.Custom {
			emoji = discord.APIEmoji("emoji:" + r.Emote)
		}

		err = ctx.State.React(msg.ChannelID, msg.ID, emoji)
		if err != nil {
			ctx.Send("I couldn't react to the message.")
			return
		}
	}

	_, err = ctx.Sendf("Success! Added %v reaction roles to the given message.", len(roles))
	return
}

type role struct {
	Emote  string
	Custom bool
	Role   discord.Role
}

var errNoPairs = errors.New("no pairs")

func (bot *Bot) parseRoles(ctx *bcr.Context, args []string) (roles []role, err error) {
	if len(args)%2 != 0 {
		return nil, errNoPairs
	}

	for i := 0; i < len(args); i = i + 2 {
		r := role{
			Emote: args[i],
		}

		if !idRe.MatchString(args[i]) {
			if emojiMatch.MatchString(args[i]) {
				r.Emote = emojiMatch.FindStringSubmatch(args[i])[3]
				r.Custom = true
			}
		} else {
			r.Custom = true
		}

		role, err := ctx.ParseRole(args[i+1])
		if err != nil {
			return nil, err
		}

		r.Role = *role

		roles = append(roles, r)
	}

	return roles, nil
}
