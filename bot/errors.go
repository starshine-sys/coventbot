package bot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/httputil"
	"github.com/getsentry/sentry-go"
	"github.com/starshine-sys/bcr"
)

// Report reports a issue to Sentry, if it's enabled
func (bot *Bot) Report(ctx *bcr.Context, e error) (err error) {
	bot.Sugar.Errorf("Error in %v (%v), guild %v: %v", ctx.Channel.ID, ctx.Channel.Name, ctx.Message.GuildID, e)

	if bot.Hub == nil {
		return
	}

	hub := bot.Hub.Clone()

	// add the user's ID
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{ID: ctx.Author.ID.String()})
	})

	// add some more info
	hub.AddBreadcrumb(&sentry.Breadcrumb{
		Category: "cmd",
		Data: map[string]interface{}{
			"user":    ctx.Author.ID,
			"channel": ctx.Channel.ID,
			"guild":   ctx.Message.GuildID,
			"command": ctx.Command,
		},
		Level:     sentry.LevelError,
		Timestamp: time.Now().UTC(),
	}, nil)

	m := ctx.NewMessage()

	var id *sentry.EventID
	if IsOurProblem(e) {
		id = hub.CaptureException(e)
	}

	m = m.Embed(&discord.Embed{
		Title:       "Internal error occurred",
		Description: "An internal error has occurred. If this issue persists, please contact the bot developer in the support server (linked in the help command).",
		Color:       0xE74C3C,

		Timestamp: discord.NowTimestamp(),
	})

	if id != nil {
		m = m.Content(fmt.Sprintf("Error code: ``%v``", bcr.EscapeBackticks(string(*id))))
		m.Data.Embed.Description = "An internal error has occurred. If this issue persists, please contact the bot developer in the support server (linked in the help command) with the error code above."
		m.Data.Embed.Footer = &discord.EmbedFooter{
			Text: string(*id),
		}
	}

	_, err = m.Send()
	return
}

// IsOurProblem checks if an error is "our problem", as in, should be in the logs and reported to Sentry.
// Will be expanded eventually once we get more insight into what type of errors we get.
func IsOurProblem(e error) bool {
	switch e.(type) {
	case *strconv.NumError:
		// this is because the user inputted an invalid number for string conversion
		// we should handle this in the command itself instead but we're lazy, and this shouldn't come up in normal usage, only with admin commands
		return false
	case *httputil.HTTPError:
		// usually caused by a message being deleted while we're still doing stuff with it (so if someone selects an option in the search results before the bot is done adding reactions)
		return false
	}

	// ignore some specific errors
	switch e {
	case bcr.ErrBotMissingPermissions:
		return false
	case bcr.ErrorNotEnoughArgs, bcr.ErrorTooManyArgs, bcr.ErrInvalidMention, bcr.ErrChannelNotFound, bcr.ErrMemberNotFound, bcr.ErrUserNotFound, bcr.ErrRoleNotFound:
		// we're not sure if these are ever returned, but ignore them anyway
		return false
	}

	return true
}
