package pklog

import (
	"context"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

// Bot ...
type Bot struct {
	*bot.Bot

	WebhookCache *ttlcache.Cache
}

// Init ...
func Init(bot *bot.Bot) (s string, list []*bcr.Command) {
	s = "PluralKit logging"

	b := &Bot{Bot: bot}

	b.WebhookCache = ttlcache.NewCache()
	b.WebhookCache.SetCacheSizeLimit(2000)
	b.WebhookCache.SetTTL(1 * time.Hour)

	b.State.AddHandler(b.pkMessageCreate)
	b.State.AddHandler(b.messageDelete)

	c := b.Router.AddCommand(&bcr.Command{
		Name:    "pk-log",
		Aliases: []string{"pklog"},
		Summary: "Set the PluralKit logging channel.",
		Args:    bcr.MinArgs(1),

		Permissions: discord.PermissionManageGuild,
		Command:     b.setChannel,
	})

	c.AddSubcommand(&bcr.Command{
		Name:    "clear-cache",
		Aliases: []string{"clearcache"},
		Summary: "Clear the webhook cache for this server.",

		Permissions: discord.PermissionManageGuild,
		Command:     b.resetCache,
	})

	// start clean messages loop
	go b.cleanMessages()

	return s, append(list, c)
}

func (bot *Bot) cleanMessages() {
	for {
		c, err := bot.DB.Pool.Exec(context.Background(), "delete from pk_messages where msg_id < $1", discord.NewSnowflake(time.Now().UTC().Add(-720*time.Hour)))
		if err != nil {
			time.Sleep(1 * time.Minute)
			continue
		}

		if n := c.RowsAffected(); n != 0 {
			bot.Sugar.Debugf("Deleted %v messages older than 30 days.", n)
		}

		time.Sleep(1 * time.Minute)
	}
}
