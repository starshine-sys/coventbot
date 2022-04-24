package bot

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/v3/gateway"
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

	// don't try commands if it can't *possibly* be one
	if !bot.Router.MatchPrefix(m.Message) {
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
