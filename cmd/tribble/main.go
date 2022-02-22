package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/gateway/shard"
	"github.com/diamondburned/arikawa/v3/state"
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
	"github.com/starshine-sys/tribble/customcommands"
	"github.com/starshine-sys/tribble/db"
	"github.com/starshine-sys/tribble/etc"
	"github.com/starshine-sys/tribble/gatekeeper"
	"github.com/starshine-sys/tribble/keyroles"
	"github.com/starshine-sys/tribble/levels"
	"github.com/starshine-sys/tribble/names"
	"github.com/starshine-sys/tribble/quotes"
	"github.com/starshine-sys/tribble/reactroles"
	"github.com/starshine-sys/tribble/starboard"
	"github.com/starshine-sys/tribble/termora"
	"github.com/starshine-sys/tribble/tickets"
	"github.com/starshine-sys/tribble/todos"
)

const intents = bcr.RequiredIntents | gateway.IntentGuildMembers | gateway.IntentGuildVoiceStates | gateway.IntentGuildPresences

// all of the bot's modules
var modules = []func(*bot.Bot) (string, []*bcr.Command){
	static.Init,
	moderation.Init,
	levels.Init,
	reactroles.Init,
	reminders.Init,
	roles.Init,
	tags.Init,
	config.Init,
	todos.Init,
	chanmirror.Init,
	notes.Init,
	starboard.Init,
	gatekeeper.Init,
	approval.Init,
	names.Init,
	admin.Init,
	tickets.Init,
	keyroles.Init,
	quotes.Init,
	termora.Init,
	customcommands.Init,
}

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
	nsfn := state.NewShardFunc(func(_ *shard.Manager, s *state.State) {
		s.AddIntents(intents)
	})

	mgr, err := shard.NewIdentifiedManager(gateway.IdentifyData{
		Token: "Bot " + c.Token,
		Properties: gateway.IdentifyProperties{
			Browser: "Discord iOS",
		},
		LargeThreshold: 250,
	}, nsfn)
	if err != nil {
		sugar.Fatal("Error creating shard manager: %v", err)
	}

	owners := make([]string, 0)
	for _, o := range c.Owners {
		owners = append(owners, o.String())
	}

	r := bcr.New(mgr, owners, c.Prefixes)
	bcrbot := bcrbot.NewWithRouter(r)

	bot := bot.New(bcrbot, sugar, db, c)
	// set default embed colour
	bot.Router.EmbedColor = etc.ColourBlurple

	// set the bot's prefix and owners
	bot.Prefix(c.Prefixes...)
	bot.Owner(c.Owners...)

	// add modules
	bot.Add(modules...)

	s, _ := bot.Router.StateFromGuildID(0)

	botUser, _ := s.Me()
	sugar.Infof("User: %v#%v (%v)", botUser.Username, botUser.Discriminator, botUser.ID)
	bot.Router.Bot = botUser

	// connect to discord
	if err := bot.Start(context.Background()); err != nil {
		sugar.Fatal("Failed to connect:", err)
	}

	// Defer this to make sure that things are always cleanly shutdown even in the event of a crash
	defer func() {
		db.Pool.Close()
		sugar.Info("Closed database connection.")
		bot.Router.ShardManager.Close()
		sugar.Info("Disconnected from Discord.")
	}()

	sugar.Info("Connected to Discord. Press Ctrl-C or send an interrupt signal to stop.")

	// start event scheduler
	go bot.Scheduler.Start()

	// sync slash commands
	err = bot.Router.SyncCommands(c.SyncCommandsIn...)
	if err == nil {
		sugar.Infof("Synced slash commands!")
		if len(c.SyncCommandsIn) != 0 {
			sugar.Infof("Synced in %v", c.SyncCommandsIn)
		}
	} else {
		sugar.Errorf("Couldn't sync slash commands: %v", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	sugar.Infof("Interrupt signal received. Shutting down...")
}
