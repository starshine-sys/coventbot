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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/starshine-sys/tribble/approval"
	"github.com/starshine-sys/tribble/bot"
	"github.com/starshine-sys/tribble/chanmirror"
	"github.com/starshine-sys/tribble/commands/admin"
	"github.com/starshine-sys/tribble/commands/config"
	"github.com/starshine-sys/tribble/commands/moderation"
	"github.com/starshine-sys/tribble/commands/moderation/notes"
	"github.com/starshine-sys/tribble/commands/reminders"
	"github.com/starshine-sys/tribble/commands/roles"
	"github.com/starshine-sys/tribble/commands/static"
	"github.com/starshine-sys/tribble/commands/tags"
	"github.com/starshine-sys/tribble/db"
	"github.com/starshine-sys/tribble/etc"
	"github.com/starshine-sys/tribble/gatekeeper"
	"github.com/starshine-sys/tribble/keyroles"
	"github.com/starshine-sys/tribble/levels"
	"github.com/starshine-sys/tribble/names"
	"github.com/starshine-sys/tribble/pklog"
	"github.com/starshine-sys/tribble/reactroles"
	"github.com/starshine-sys/tribble/starboard"
	"github.com/starshine-sys/tribble/tickets"
	"github.com/starshine-sys/tribble/todos"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// get the config
	c := getConfig()

	// set up a logger
	zcfg := zap.NewProductionConfig()
	zcfg.Encoding = "console"
	zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zcfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zcfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	zcfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	if c.DebugLogging {
		zcfg.Level.SetLevel(zapcore.DebugLevel)
	} else {
		zcfg.Level.SetLevel(zapcore.InfoLevel)
	}

	zap, err := zcfg.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		panic(err)
	}
	sugar := zap.Sugar()

	if c.DebugLogging {
		sugar.Warn("Debug logging enabled. Set \"debug_logging\" in config.yaml to \"false\" to disable.")
	}

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
	// add todos
	bot.Add(todos.Init)
	// add role commands
	bot.Add(roles.Init)
	// add config commands
	bot.Add(config.Init)

	chanmirror.Init(bot)
	notes.Init(bot)
	starboard.Init(bot)
	gatekeeper.Init(bot)
	approval.Init(bot)
	names.Init(bot)
	pklog.Init(bot)
	admin.Init(bot)
	levels.Init(bot)
	reactroles.Init(bot)
	reminders.Init(bot)
	tags.Init(bot)
	tickets.Init(bot)
	keyroles.Init(bot)

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
