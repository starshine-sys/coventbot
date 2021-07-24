package reactroles

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) new(ctx *bcr.Context) (err error) {
	ch, err := ctx.ParseChannel(ctx.Args[0])
	if err != nil {
		_, err = ctx.Send("Couldn't parse a channel from your input.")
		return
	}

	name := ctx.Args[1]

	if ch.GuildID != ctx.Channel.GuildID || (ch.Type != discord.GuildText && ch.Type != discord.GuildNews) {
		_, err = ctx.Send("The channel you gave isn't in this server.")
		return
	}

	roles, err := bot.parseRoles(ctx, ctx.Args[2:])
	if err != nil {
		if err == errNoPairs {
			_, err = ctx.Send("You must give emoji-role *pairs*.")
			return
		}
		_, err = ctx.Send("Couldn't parse one or more roles.")
		return
	}

	e := discord.Embed{
		Color: 0x7ED321,
		Author: &discord.EmbedAuthor{
			Name: name,
			Icon: "https://cdn.discordapp.com/emojis/757537919794937936.png",
		},
	}

	description, _ := ctx.Flags.GetString("description")
	if description != "" {
		e.Description = description + "\n\n"
	}

	showAsMention, _ := ctx.Flags.GetBool("mention")

	for _, r := range roles {
		emoji := r.Emote
		if r.Custom {
			emoji = "<:emoji:" + r.Emote + ">"
		}

		if showAsMention {
			e.Description += fmt.Sprintf("%v %v\n", emoji, r.Role.Mention())
		} else {
			e.Description += fmt.Sprintf("%v %v\n", emoji, r.Role.Name)
		}
	}

	msg, err := ctx.State.SendEmbeds(ch.ID, e)
	if err != nil {
		_, err = ctx.Send("I couldn't send a message in the target channel.")
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), `insert into react_roles
	(server_id, channel_id, message_id, title, description, mention) values ($1, $2, $3, $4, $5, $6) on conflict (message_id) do nothing`, msg.GuildID, msg.ChannelID, msg.ID, name, description, showAsMention)
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

	_, err = ctx.Sendf("Success! Added %v reaction roles to the given message.\nNote that my highest role must be above the roles you added to this message for the reactions to work.", len(roles))
	return
}
