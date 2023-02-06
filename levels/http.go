// SPDX-License-Identifier: AGPL-3.0-only
package levels

import (
	"html/template"
	"net/http"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/dustin/go-humanize"
	"github.com/go-chi/chi/v5"

	_ "embed"
)

type webEntry struct {
	Name  string
	Rank  int
	Level int64
	XP    string
}

//go:embed leaderboard.html
var leaderboardHtml string

var leaderboardTmpl = template.Must(template.New("").Parse(leaderboardHtml))

func (bot *Bot) webLeaderboard(w http.ResponseWriter, r *http.Request) {
	id, err := discord.ParseSnowflake(chi.URLParam(r, "id"))
	if err != nil {
		bot.ShowText(w, "Leaderboard", "Error", "Could not parse %q as a server ID.", chi.URLParam(r, "id"))

		return
	}

	s, _ := bot.Router.StateFromGuildID(discord.GuildID(id))
	g, err := s.Guild(discord.GuildID(id))
	if err != nil {
		bot.ShowText(w, "Leaderboard", "Error", "Could not get this server. This is most likely a bug!")
		bot.Sugar.Errorf("getting guild config: %v", err)

		return
	}

	sc, err := bot.getGuildConfig(discord.GuildID(id))
	if err != nil {
		bot.ShowText(w, "Leaderboard", "Error", "Could not get this server's configuration. This is a bug!")
		bot.Sugar.Errorf("getting guild config: %v", err)

		return
	}

	noRanks, _ := bot.DB.GuildBoolGet(sc.ID, "levels:disable_ranks")
	if noRanks || sc.LeaderboardModOnly || !sc.LevelsEnabled {
		bot.ShowText(w, "Leaderboard", "Leaderboard", "The leaderboard is disabled for this server.")

		return
	}

	lb, err := bot.getLeaderboard(sc.ID, true)
	if err != nil {
		bot.ShowText(w, "Leaderboard", "Error", "Could not get this server's leaderboard.")
		bot.Sugar.Errorf("getting guild leaderboard: %v", err)

		return
	}

	levels := make([]webEntry, 0, len(lb))
	gm := bot.Members(sc.ID)
	for _, l := range lb {
		for _, m := range gm {
			if m.User.ID == l.UserID {
				name := m.User.Username
				if m.Nick != "" {
					name = m.Nick
				}

				levels = append(levels, webEntry{
					Name:  name,
					XP:    humanize.Comma(l.XP),
					Level: sc.CalculateLevel(l.XP),
				})
				break
			}
		}
	}

	for i := range levels {
		levels[i].Rank = i + 1
	}

	if len(levels) == 0 {
		bot.ShowText(w, "Leaderboard", "Leaderboard", "There doesn't seem to be anyone on the leaderboard...")

		return
	}

	err = leaderboardTmpl.Execute(w, struct {
		Guild       string
		Leaderboard []webEntry
	}{Guild: g.Name, Leaderboard: levels})
	if err != nil {
		bot.Sugar.Errorf("executing leaderboard template: %v", err)
	}
}
