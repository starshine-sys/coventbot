package main

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/starshine-sys/bcr"
	bcrbot "github.com/starshine-sys/bcr/bot"
	"github.com/starshine-sys/tribble/approval"
	"github.com/starshine-sys/tribble/bot"
	"github.com/starshine-sys/tribble/commands/admin"
	"github.com/starshine-sys/tribble/commands/config"
	"github.com/starshine-sys/tribble/commands/moderation"
	"github.com/starshine-sys/tribble/commands/reminders"
	"github.com/starshine-sys/tribble/commands/static"
	"github.com/starshine-sys/tribble/commands/tags"
	"github.com/starshine-sys/tribble/db"
	"github.com/starshine-sys/tribble/etc"
	"github.com/starshine-sys/tribble/gatekeeper"
	"github.com/starshine-sys/tribble/names"
	"github.com/starshine-sys/tribble/pklog"
	"github.com/starshine-sys/tribble/starboard"
	"github.com/starshine-sys/tribble/tickets"
	"go.uber.org/zap"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// set up a logger
	zap, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	sugar := zap.Sugar()

	// get the config
	c := getConfig(sugar)

	// open the database
	db, err := db.New(c.DatabaseURL, sugar, c)
	if err != nil {
		sugar.Fatalf("Error connecting to database: %v", err)
	}
	sugar.Info("Connected to database.")

	// create a new bot
	r, err := bcr.NewWithIntents(c.Token, c.Owners, c.Prefixes, bcr.RequiredIntents|gateway.IntentGuildMembers)
	if err != nil {
		sugar.Fatal("Error creating bot:", err)
	}
	bcrbot := bcrbot.NewWithRouter(r)

	bot := bot.New(bcrbot, sugar, db, c)
	// set default embed colour
	bot.Router.EmbedColor = etc.ColourBlurple

	// set the bot's prefix and owners
	bot.Prefix(c.Prefixes...)
	bot.Owner(c.Owners...)

	// add basic commands
	bot.Add(static.Init)
	// add mod commands
	bot.Add(moderation.Init)
	// add reminder commands
	bot.Add(reminders.Init)
	// add tag commands
	bot.Add(tags.Init)
	// add ticket commands
	bot.Add(tickets.Init)
	// add config commands
	bot.Add(config.Init)
	// add starboard
	starboard.Init(bot)
	// add gatekeeper
	bot.Add(gatekeeper.Init)
	// add approval
	bot.Add(approval.Init)
	// add logging
	bot.Add(names.Init)
	// add PK logging
	bot.Add(pklog.Init)
	// add admin commands
	bot.Add(admin.Init)

	// connect to discord
	if err := bot.Start(); err != nil {
		sugar.Fatal("Failed to connect:", err)
	}

	// Defer this to make sure that things are always cleanly shutdown even in the event of a crash
	defer func() {
		db.Pool.Close()
		sugar.Info("Closed database connection.")
		bot.Router.State.Close()
		bot.Router.State.Gateway.Close()
		sugar.Info("Disconnected from Discord.")
	}()

	sugar.Info("Connected to Discord. Press Ctrl-C or send an interrupt signal to stop.")

	botUser, _ := bot.Router.State.Me()
	sugar.Infof("User: %v#%v (%v)", botUser.Username, botUser.Discriminator, botUser.ID)
	bot.Router.Bot = botUser

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	sugar.Infof("Interrupt signal received. Shutting down...")
}
