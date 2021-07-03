package todos

import (
	"fmt"
	"strconv"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/starshine-sys/bcr"
)

func (bot *Bot) cmdComplete(ctx *bcr.Context) (err error) {
	id, err := strconv.ParseInt(ctx.RawArgs, 0, 0)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "You didn't give a valid ID.")
	}

	t, err := bot.getTodo(id, ctx.Author.ID)
	if err != nil {
		_, err = ctx.Replyc(bcr.ColourRed, "Couldn't find a todo with that ID.")
		return
	}

	if t.Complete {
		_, err = ctx.Replyc(bcr.ColourRed, "That todo is already completed!")
		return
	}

	err = bot.complete(ctx.State, t)
	if err != nil {
		bot.Sugar.Errorf("Error completing todo ID %v: %v", t.ID, err)
		return bot.Report(ctx, err)
	}

	_, err = ctx.Reply("Completed todo #%v!", t.ID)
	return
}

func (bot *Bot) reactionAdd(ev *gateway.MessageReactionAddEvent) {
	if ev.Emoji.APIString() != "✅" {
		return
	}

	t, err := bot.getTodoMessage(ev.MessageID, ev.UserID)
	if err != nil {
		return
	}

	if t.Complete {
		return
	}

	s, _ := bot.Router.StateFromGuildID(ev.GuildID)

	err = bot.complete(s, t)
	if err != nil {
		bot.Sugar.Errorf("Error completing todo ID %v: %v", t.ID, err)
	}
}

func (bot *Bot) complete(s *state.State, t *Todo) (err error) {
	msg, err := s.Message(t.ChannelID, t.MID)
	if err != nil {
		return err
	}

	var e discord.Embed
	if len(msg.Embeds) > 0 {
		e = msg.Embeds[0]

		e.Title = "Todo completed"
		e.Color = bcr.ColourGreen

		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Completed at",
			Value: time.Now().UTC().Format("2006-01-02 15:04") + " UTC",
		})
	} else {
		jumpLink := "https://discord.com/channels/"
		if !t.OrigServerID.IsValid() {
			jumpLink += "@me/"
		} else {
			jumpLink += t.OrigServerID.String() + "/"
		}
		jumpLink += fmt.Sprintf("%v/%v", t.OrigChannelID, t.OrigMID)

		e = discord.Embed{
			Title:       "Todo",
			Color:       bcr.ColourBlurple,
			Description: t.Description,

			Fields: []discord.EmbedField{{
				Name:  "Source",
				Value: fmt.Sprintf("[Jump!](%v)", jumpLink),
			}},

			Timestamp: discord.NewTimestamp(t.Created),
		}
	}

	err = bot.markComplete(t.ID)
	if err != nil {
		bot.Sugar.Error(err)
		return
	}

	_, err = s.EditEmbeds(t.ChannelID, t.MID, e)
	return
}
