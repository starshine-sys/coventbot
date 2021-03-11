package starboard

import (
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

// Bot ...
type Bot struct {
	*bot.Bot

	mu map[discord.MessageID]*sync.Mutex
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Starboard"

	b := &Bot{
		Bot: bot,
		mu:  make(map[discord.MessageID]*sync.Mutex),
	}

	b.State.AddHandler(b.MessageReactionAdd)
	b.State.AddHandler(b.MessageReactionDelete)
	b.State.AddHandler(b.MessageReactionRemoveEmoji)
	return
}
