package bot

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	bcr2 "github.com/starshine-sys/bcr/v2"
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

			_ = c.Execute(data)
			_ = bot.Router.ShardManager.Shard(0).(*state.State).React(m.ChannelID, m.ID, "✅")
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
}

func (bot *Bot) interactionCreate(ic *gateway.InteractionCreateEvent) {
	err := bot.Interactions.Execute(ic)
	if err == bcr2.ErrUnknownCommand {
		return
	} else if err != nil {
		bot.Sugar.Error("error in bcr v2 handler:", err)
	}
}
