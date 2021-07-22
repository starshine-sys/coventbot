package bot

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/gateway/shard"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/starshine-sys/bcr"
)

// MessageCreate is run on a message create event
func (bot *Bot) MessageCreate(m *gateway.MessageCreateEvent) {
	// set the bot user if not done already
	if bot.Router.Bot == nil {
		err := bot.Router.SetBotUser(m.GuildID)
		if err != nil {
			bot.Sugar.Fatal(err)
		}
		bot.Router.Prefixes = append(bot.Router.Prefixes, fmt.Sprintf("<@%v>", bot.Router.Bot.ID), fmt.Sprintf("<@!%v>", bot.Router.Bot.ID))
	}

	bot.Counters.Mu.Lock()
	bot.Counters.Messages++
	if strings.Contains(m.Content, fmt.Sprintf("<@%v>", bot.Router.Bot.ID)) || strings.Contains(m.Content, fmt.Sprintf("<@!%v>", bot.Router.Bot.ID)) {
		bot.Counters.Mentions++
	}
	bot.Counters.Mu.Unlock()

	// if the author is a bot, return
	if m.Author.Bot {
		return
	}

	// if the message does not start with any of the bot's prefixes (including mentions), return
	if !bot.Router.MatchPrefix(m.Message) {
		// try DMing the user
		if !m.GuildID.IsValid() && bot.Config.DMs.Open {
			// if the user is blocked, return
			for _, u := range bot.Config.DMs.BlockedUsers {
				if u == m.Author.ID {
					return
				}
			}

			// if there's no DM webhook set, return
			if !bot.Config.DMs.Webhook.ID.IsValid() {
				return
			}

			data := webhook.ExecuteData{
				Username:  bot.Router.Bot.Username,
				AvatarURL: bot.Router.Bot.AvatarURL(),

				Content: "> Received a DM!",

				Embeds: append([]discord.Embed{
					{
						Author: &discord.EmbedAuthor{
							Icon: m.Author.AvatarURL(),
							Name: m.Author.Username + "#" + m.Author.Discriminator,
						},
						Description: m.Author.ID.String(),
						Color:       bot.Router.EmbedColor,
					},
					{
						Author: &discord.EmbedAuthor{
							Icon: m.Author.AvatarURL(),
							Name: m.Author.Username + "#" + m.Author.Discriminator,
						},
						Description: m.Content,
						Color:       bot.Router.EmbedColor,
						Footer: &discord.EmbedFooter{
							Text: m.ID.String(),
						},
						Timestamp: m.Timestamp,
					},
				}, m.Embeds...),
			}

			c := webhook.New(bot.Config.DMs.Webhook.ID, bot.Config.DMs.Webhook.Token)

			c.Execute(data)

			bot.Router.ShardManager.Shard(0).(shard.ShardState).Shard.(*state.State).React(m.ChannelID, m.ID, "✅")
		}
		return
	}

	// get the context
	ctx, err := bot.Router.NewContext(m)
	if err != nil {
		bot.Sugar.Errorf("Error getting context: %v", err)
		return
	}

	err = bot.Router.Execute(ctx)
	if err != nil {
		bot.Sugar.Errorf("Error executing commands: %v", err)
		return
	}

	// handle tags as commands
	err = bot.handleTagCommand(ctx)
	if err != nil {
		bot.Sugar.Errorf("Error sending message: %v", err)
		return
	}
}

func (bot *Bot) handleTagCommand(ctx *bcr.Context) (err error) {
	name := ctx.Command + " " + ctx.RawArgs
	if ctx.RawArgs == "" {
		name = ctx.Command
	}

	t, err := bot.DB.GetTag(ctx.Message.GuildID, name)
	if err != nil {
		return nil
	}

	tr := false

	if len(ctx.Message.Mentions) > 0 {
		tr = true
	}

	data := api.SendMessageData{
		Content: t.Response,
		AllowedMentions: &api.AllowedMentions{
			Parse:       []api.AllowedMentionType{},
			RepliedUser: option.Bool(&tr),
		},
	}

	if ctx.Message.Reference != nil {
		data.Reference = &discord.MessageReference{
			MessageID: ctx.Message.Reference.MessageID,
		}
	}

	_, err = ctx.State.SendMessageComplex(ctx.Message.ChannelID, data)
	return err
}

func (bot *Bot) interactionCreate(ic *gateway.InteractionCreateEvent) {
	bot.Sugar.Infof("Received interaction create event")

	if ic.Type != gateway.CommandInteraction {
		return
	}

	bot.Sugar.Infof("Received command event")

	ctx, err := bot.Router.NewSlashContext(ic)
	if err != nil {
		bot.Sugar.Errorf("Couldn't create slash context: %v", err)
		return
	}

	bot.Sugar.Infof("Created slash context with command %v", ctx.CommandName)

	err = bot.Router.ExecuteSlash(ctx)
	if err != nil {
		bot.Sugar.Errorf("Couldn't create slash context: %v", err)
	}
}
