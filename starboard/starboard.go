package starboard

import (
	"context"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/webhook"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

// Bot ...
type Bot struct {
	*bot.Bot

	mu map[discord.MessageID]*sync.Mutex

	webhooks   map[discord.WebhookID]*webhook.Client
	webhooksMu sync.Mutex
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "Starboard"

	b := &Bot{
		Bot:      bot,
		mu:       make(map[discord.MessageID]*sync.Mutex),
		webhooks: make(map[discord.WebhookID]*webhook.Client),
	}

	b.Router.AddHandler(b.guildCreate)

	b.Router.AddHandler(b.MessageReactionAdd)
	b.Router.AddHandler(b.MessageReactionDelete)
	b.Router.AddHandler(b.MessageReactionRemoveEmoji)
	return
}

func (bot *Bot) guildCreate(ev *gateway.GuildCreateEvent) {
	bot.updateStarboardWebhook(ev.ID)
}

func (bot *Bot) updateStarboardWebhook(id discord.GuildID) {
	time.Sleep(time.Second)
	conf, err := bot.DB.Starboard(id)
	if err != nil {
		bot.Sugar.Errorf("error getting starboard config for %v: %v", id, err)
		return
	}

	if !conf.StarboardChannel.IsValid() {
		return
	}

	_, err = bot.DB.StarboardChannelWebhook(conf.StarboardChannel)
	if err == nil || errors.Cause(err) != pgx.ErrNoRows {
		if err != nil {
			bot.Sugar.Errorf("error getting starboard webhook for %v: %v", id, err)
		}

		return
	}

	// create a webhook
	s, _ := bot.Router.StateFromGuildID(id)
	wh, err := s.CreateWebhook(conf.StarboardChannel, api.CreateWebhookData{
		Name: s.Ready().User.Username + " Starboard",
	})
	if err != nil {
		bot.Sugar.Errorf("error creating starboard webhook for %v: %v", id, err)
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "insert into starboard_webhooks (id, channel_id, token) values ($1, $2, $3)", wh.ID, wh.ChannelID, wh.Token)
	if err != nil {
		bot.Sugar.Errorf("error saving starboard webhook for %v: %v", id, err)
		return
	}

	bot.Sugar.Debugf("created missing starboard webhook for %v (%v)", id, conf.StarboardChannel)
}
