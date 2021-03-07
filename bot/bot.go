package bot

import (
	"sort"

	"github.com/starshine-sys/coventbot/db"
	"github.com/starshine-sys/coventbot/types"
	"github.com/starshine-sys/bcr"
	bcrbot "github.com/starshine-sys/bcr/bot"
	"go.uber.org/zap"
)

// Bot is the main bot struct
type Bot struct {
	*bcrbot.Bot

	Config *types.BotConfig
	Sugar  *zap.SugaredLogger
	DB     *db.DB
}

// Module is a single module/category of commands
type Module interface {
	String() string
	Commands() []*bcr.Command
}

// New creates a new instance of Bot
func New(
	bot *bcrbot.Bot,
	sugar *zap.SugaredLogger,
	db *db.DB,
	config *types.BotConfig) *Bot {
	b := &Bot{
		Bot:    bot,
		Sugar:  sugar,
		DB:     db,
		Config: config,
	}
	return b
}

// Add adds a module to the bot
func (bot *Bot) Add(f func(*Bot) (string, []*bcr.Command)) {
	m, c := f(bot)

	// sort the list of commands
	sort.Sort(bcr.Commands(c))

	// add the module
	bot.Modules = append(bot.Modules, &botModule{
		name:     m,
		commands: c,
	})
}

type botModule struct {
	name     string
	commands bcr.Commands
}

// String returns the module's name
func (b botModule) String() string {
	return b.name
}

// Commands returns a list of commands
func (b *botModule) Commands() []*bcr.Command {
	return b.commands
}
