package bot

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway/shard"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/handler"
	"github.com/getsentry/sentry-go"
	"github.com/starshine-sys/bcr"
	bcrbot "github.com/starshine-sys/bcr/bot"
	"github.com/starshine-sys/pkgo"
	"github.com/starshine-sys/tribble/db"
	"github.com/starshine-sys/tribble/types"
	"go.uber.org/zap"
)

// Bot is the main bot struct
type Bot struct {
	*bcrbot.Bot

	Config *types.BotConfig
	Sugar  *zap.SugaredLogger
	DB     *db.DB

	GuildLogWebhook *webhook.Client

	Counters struct {
		Mu sync.Mutex

		Mentions int
		Messages int
	}

	Hub *sentry.Hub

	PK *pkgo.Session

	members   map[memberKey]member
	membersMu sync.RWMutex

	HelperRole bcr.CustomPerms
	ModRole    bcr.CustomPerms
	AdminRole  bcr.CustomPerms
}

type member struct {
	discord.Member
	GuildID discord.GuildID
}

type memberKey struct {
	GuildID discord.GuildID
	UserID  discord.UserID
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
		PK:     pkgo.New(""),

		members: map[memberKey]member{},
	}

	b.HelperRole = &HelperRole{b}
	b.ModRole = &ModRole{b}
	b.AdminRole = &AdminRole{b}

	// create a Sentry config
	if config.SentryURL != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: config.SentryURL,
		})
		if err != nil {
			sugar.Fatalf("sentry.Init: %s", err)
		}
		sugar.Infof("Initialised Sentry")
		// defer this to flush buffered events
		defer sentry.Flush(2 * time.Second)
	}
	b.Hub = sentry.CurrentHub()
	if config.SentryURL == "" {
		b.Hub = nil
	}

	if config.GuildLogWebhook.ID.IsValid() {
		b.GuildLogWebhook = webhook.New(config.GuildLogWebhook.ID, config.GuildLogWebhook.Token)
	}

	// set the prefix checker
	b.Router.Prefixer = b.CheckPrefix

	b.Router.ShardManager.ForEach(func(s shard.Shard) {
		state := s.(*state.State)

		// add guild create handler
		state.AddHandler(b.GuildCreate)

		// add guild remove handler
		state.PreHandler = handler.New()
		state.PreHandler.Synchronous = true
		state.PreHandler.AddHandler(b.guildDelete)

		// add message create handler
		state.AddHandler(b.MessageCreate)

		// add member update handler (this isn't handled by default apparently?)
		state.AddHandler(b.guildMemberUpdate)

		// add cache handlers
		state.AddHandler(b.requestGuildMembers)
		state.AddHandler(b.guildMemberChunk)
		state.AddHandler(b.memberUpdateEvent)
		state.AddHandler(b.memberAddEvent)
		state.AddHandler(b.memberRemoveEvent)
	})

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

// PagedEmbed ...
func (bot *Bot) PagedEmbed(ctx *bcr.Context, embeds []discord.Embed, timeout time.Duration) (msg *discord.Message, err error) {
	var reactions bool
	bot.DB.Pool.QueryRow(context.Background(), "select reaction_pages from user_config where user_id = $1", ctx.Author.ID).Scan(&reactions)

	if reactions {
		msg, _, err = ctx.PagedEmbedTimeout(embeds, false, timeout)
		return
	}

	msg, _, err = ctx.ButtonPages(embeds, timeout)
	return
}
