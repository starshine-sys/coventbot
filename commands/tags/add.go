package tags

import (
	"strings"

	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/db"
)

func (bot *Bot) addTag(ctx *bcr.Context) (err error) {
	args := strings.Split(ctx.RawArgs, "\n")
	if len(args) < 2 {
		_, err = ctx.Send("Not enough arguments given: need at least 2, separated by a newline.")
		return
	}

	t := db.Tag{
		Name:     strings.TrimSpace(args[0]),
		Response: strings.TrimSpace(strings.Join(args[1:], "\n")),
	}

	t, err = bot.DB.AddTag(ctx, t)
	if err != nil {
		_, err = ctx.Send("An error occurred while saving your tag. Are you sure the name is unique?")
		return
	}

	_, err = ctx.Sendf("✅ Added tag `%v`. (ID: %s)", strings.ToLower(t.Name), t.ID)
	return
}
