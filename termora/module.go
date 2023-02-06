// SPDX-License-Identifier: AGPL-3.0-only
package termora

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"emperror.dev/errors"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/bot"
)

type Bot struct {
	*bot.Bot

	HTTP *http.Client
}

func Init(b *bot.Bot) (s string, list []*bcr.Command) {
	bot := &Bot{
		Bot:  b,
		HTTP: &http.Client{},
	}

	bot.Router.AddHandler(bot.messageCreate)

	return
}

type Term struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Aliases []string `json:"aliases"`
}

const baseURL = "https://api.termora.org/v1"

func (bot *Bot) Terms() (ts []Term, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/list", nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	req.Header.Add("User-Agent", fmt.Sprintf("Tribble; %v-%v/%v", bot.Router.Bot.Username, bot.Router.Bot.Discriminator, bot.Router.Bot.ID))

	resp, err := bot.HTTP.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "doing request")
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ts)
	return ts, errors.Wrap(err, "decoding request")
}
