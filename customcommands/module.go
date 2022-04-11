package customcommands

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
	"github.com/starshine-sys/tribble/customcommands/cc"
	"gitlab.com/1f320/x/parameters"
)

type Bot struct {
	*bot.Bot

	Client *http.Client
}

// Init ...
func Init(b *bot.Bot) (s string, list []*bcr.Command) {
	s = "Custom commands"

	bot := &Bot{
		Bot:    b,
		Client: &http.Client{},
	}

	bot.Scheduler.AddType(&cc.ScheduledCC{})

	bot.Router.AddHandler(bot.messageCreate)

	bot.Router.AddCommand(&bcr.Command{
		Name:              "cc",
		Summary:           "Show or create a custom command",
		Usage:             "[name]",
		Command:           bot.showOrAdd,
		CustomPermissions: bot.ManagerRole,
	})

	return
}

func (bot *Bot) messageCreate(m *gateway.MessageCreateEvent) {
	if !m.GuildID.IsValid() || m.Author.Bot {
		return
	}

	allowCC := false
	for _, id := range bot.Config.AllowCCs {
		if id == m.GuildID {
			allowCC = true
			break
		}
	}
	if !allowCC {
		return
	}

	if !bot.Router.MatchPrefix(m.Message) {
		return
	}

	ctx, err := bot.Router.NewContext(m)
	if err != nil {
		bot.Sugar.Errorf("Error getting context: %v", err)
		return
	}

	i := bot.CheckPrefix(ctx.Message)
	if i == -1 {
		return
	}
	content := ctx.Message.Content[i:]
	content = strings.TrimSpace(content)

	params := parameters.NewParameters(content, false)
	cmdName := strings.ToLower(params.Pop())
	if cmdName == "" {
		return
	}

	cmd, err := bot.DB.CustomCommand(ctx.Message.GuildID, cmdName)
	if err != nil {
		return
	}

	t := time.Now()
	s := cc.NewState(bot.Bot, ctx, params)

	cctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	err = s.Do(cctx, cmd.ID, cmd.Source)
	if err != nil {
		if s.FilterErrors(err) {
			return
		}

		_, err = ctx.State.SendMessageComplex(ctx.Message.ChannelID, api.SendMessageData{
			Content:         fmt.Sprintf("An error occurred while executing the custom command ``%v/%v``:\n```lua\n%v\n```", cmd.ID, bcr.EscapeBackticks(cmd.Name), err),
			AllowedMentions: &api.AllowedMentions{},
		})
		if err != nil {
			bot.Sugar.Errorf("error sending message: %v", err)
		}
	}

	bot.Sugar.Debugf("Executed custom command %v/%v in %v", cmd.ID, cmd.Name, time.Since(t))
}
