package bot

import (
	"fmt"
	"strings"
	"time"

	"1f320.xyz/x/parameters"
	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/gateway/shard"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/customcommands/cc"
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
			_ = bot.Router.ShardManager.Shard(0).(shard.ShardState).Shard.(*state.State).React(m.ChannelID, m.ID, "✅")
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
		if err == errHadTag {
			return
		}
		bot.Sugar.Errorf("Error sending message: %v", err)
		return
	}

	bot.customCommands(ctx)
}

func (bot *Bot) customCommands(ctx *bcr.Context) {
	i := bot.CheckPrefix(ctx.Message)
	if i == -1 {
		return
	}
	content := ctx.Message.Content[i:]
	content = strings.TrimSpace(content)

	params := parameters.NewParameters(content, false)
	cmdName := strings.ToLower(params.Pop())
	if cmdName == "" {
		return
	}

	cmd, err := bot.DB.CustomCommand(ctx.Message.GuildID, cmdName)
	if err != nil {
		return
	}

	t := time.Now()
	s := cc.NewState(ctx, params)

	err = s.Do(cmd.Source, time.Minute)
	if err != nil {
		_, err = ctx.State.SendMessageComplex(ctx.Message.ChannelID, api.SendMessageData{
			Content:         fmt.Sprintf("An error occurred while executing the custom command ``%v/%v``:\n```lua\n%v\n```", cmd.ID, bcr.EscapeBackticks(cmd.Name), err),
			AllowedMentions: &api.AllowedMentions{},
		})
		if err != nil {
			bot.Sugar.Errorf("error sending message: %v", err)
		}
	}

	bot.Sugar.Debugf("Executed custom command %v/%v in %v", cmd.ID, cmd.Name, time.Since(t))
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
	if err != nil {
		return err
	}
	return errHadTag
}

const errHadTag = errors.Sentinel("message had tag")

func (bot *Bot) interactionCreate(ic *gateway.InteractionCreateEvent) {
	if ic.Type != discord.CommandInteraction {
		return
	}

	ctx, err := bot.Router.NewSlashContext(ic)
	if err != nil {
		bot.Sugar.Errorf("Couldn't create slash context: %v", err)
		return
	}

	err = bot.Router.ExecuteSlash(ctx)
	if err != nil {
		bot.Sugar.Errorf("Couldn't execute command: %v", err)
	}
}
