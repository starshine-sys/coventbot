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
	"github.com/starshine-sys/coventbot/bot"
	"github.com/starshine-sys/coventbot/commands/static"
	"github.com/starshine-sys/coventbot/commands/tags"
	"github.com/starshine-sys/coventbot/db"
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
	config := getConfig(sugar)

	// open the database
	db, err := db.New(config.DatabaseURL, sugar, config)
	if err != nil {
		sugar.Fatalf("Error connecting to database: %v", err)
	}
	sugar.Info("Connected to database.")

	// create a new bot
	r, err := bcr.NewWithIntents(config.Token, config.Owners, config.Prefixes, bcr.RequiredIntents|gateway.IntentGuildMembers)
	if err != nil {
		sugar.Fatal("Error creating bot:", err)
	}
	bcrbot := bcrbot.NewWithRouter(r)

	bot := bot.New(bcrbot, sugar, db, config)
	// set default embed colour
	bot.Router.EmbedColor = 0x4C7DAA

	// set the bot's prefix and owners
	bot.Prefix(config.Prefixes...)
	bot.Owner(config.Owners...)

	// add the command handler
	bot.Router.Session.AddHandler(bot.Router.MessageCreate)

	// add basic commands
	bot.Add(static.Init)
	// add tag commands
	bot.Add(tags.Init)

	// connect to discord
	if err := bot.Start(); err != nil {
		sugar.Fatal("Failed to connect:", err)
	}

	// Defer this to make sure that things are always cleanly shutdown even in the event of a crash
	defer func() {
		db.Pool.Close()
		sugar.Info("Closed database connection.")
		bot.Router.Session.Close()
		sugar.Info("Disconnected from Discord.")
	}()

	sugar.Info("Connected to Discord. Press Ctrl-C or send an interrupt signal to stop.")

	botUser, _ := bot.Router.Session.Me()
	sugar.Infof("User: %v#%v (%v)", botUser.Username, botUser.Discriminator, botUser.ID)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	sugar.Infof("Interrupt signal received. Shutting down...")
}
