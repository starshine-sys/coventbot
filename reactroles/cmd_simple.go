package reactroles

import (
	"context"
	"fmt"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/starshine-sys/bcr"
)

// emote IDs used for the `rr simple` command
var simpleEmotes = []string{
	"845994591954010163",
	"845994619673509898",
	"845994640090857474",
	"845994650947026955",
	"845994681379717190",
	"845994709490991105",
	"845994785743175731",
	"845994802373984286",
	"845994831322021918",
	"845994842760151050",
	"845994877405888522",
	"845994911648055336",
	"845994925953777675",
	"845994938399064084",
	"845994954148937748",
	"845994967503732756",
	"845994979431284776",
	"845994994905251850",
	"845995013238423603",
	"845995025997627432",
}

func (bot *Bot) simple(ctx *bcr.Context) (err error) {
	ch, err := ctx.ParseChannel(ctx.Args[0])
	if err != nil || ch.GuildID != ctx.Message.GuildID || (ch.Type != discord.GuildText && ch.Type != discord.GuildNews) {
		_, err = ctx.Send("Channel not found.")
		return
	}

	name := ctx.Args[1]

	rls, n := ctx.GreedyRoleParser(ctx.Args[2:])
	if len(rls) == 0 {
		_, err = ctx.Send("Couldn't parse any of the given roles.")
		return
	} else if n != -1 {
		_, err = ctx.Send("Note: not all roles could be parsed; I'm only adding the roles I could parse.")
	} else if n > 20 {
		_, err = ctx.Send("You can only have a maximum of 20 reaction roles per message.")
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

	for i, r := range rls {
		if showAsMention {
			e.Description += fmt.Sprintf("<:emoji:%v> %v\n", simpleEmotes[i], r.Mention())
		} else {
			e.Description += fmt.Sprintf("<:emoji:%v> %v\n", simpleEmotes[i], r.Name)
		}
	}

	m, err := ctx.State.SendEmbeds(ch.ID, e)
	if err != nil {
		_, err = ctx.Send("I couldn't send a message in the target channel.")
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), `insert into react_roles
	(server_id, channel_id, message_id, title, description, mention) values ($1, $2, $3, $4, $5, $6) on conflict (message_id) do nothing`, m.GuildID, m.ChannelID, m.ID, name, description, showAsMention)
	if err != nil {
		return bot.Report(ctx, err)
	}

	for i, r := range rls {
		_, err = bot.DB.Pool.Exec(context.Background(), `insert into react_role_entries
		(message_id, emote, role_id) values ($1, $2, $3) on conflict (message_id, emote) do update
		set role_id = $3`, m.ID, simpleEmotes[i], r.ID)
		if err != nil {
			return bot.Report(ctx, err)
		}

		emoji := discord.APIEmoji("emoji:" + simpleEmotes[i])

		err = ctx.State.React(m.ChannelID, m.ID, emoji)
		if err != nil {
			ctx.Send("I couldn't react to the message.")
			return
		}
	}

	_, err = ctx.Sendf("Done! Added %v react roles.\nNote that my highest role must be above the roles you added to this message for the reactions to work.", len(rls))
	return
}
