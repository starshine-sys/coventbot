package gatekeeper

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/starshine-sys/tribble/common"
)

// PendingUser ...
type PendingUser struct {
	UserID  discord.UserID
	GuildID discord.GuildID

	Key     uuid.UUID
	Pending bool
}

type gatekeeperData struct {
	UUID   string
	Config *common.BotConfig
}

var gatekeeperTmpl = template.Must(template.New("").Parse(gatekeeperHtml))

// GatekeeperGET ...
func (bot *Bot) GatekeeperGET(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		bot.ShowText(w, "Gatekeeper error", "", "Invalid UUID (%v) provided.", chi.URLParam(r, "uuid"))
		return
	}

	u, err := bot.userByUUID(id)
	if err != nil {
		bot.ShowText(w, "Internal server error", "", "An internal server error has occurred.")

		bot.Sugar.Errorf("Error getting UUID from database: %v", err)
		return
	}
	if !u.Pending {
		bot.ShowText(w, "Gatekeeper error", "", "You have already passed the gateway.", u.UserID)
		return
	}

	d := gatekeeperData{
		UUID:   id.String(),
		Config: bot.Config,
	}

	err = gatekeeperTmpl.Execute(w, d)
	if err != nil {
		fmt.Fprintf(w, "Internal server error.\n")
		bot.Sugar.Errorf("Error rending template: %v", err)
		return
	}
}

// VerifyPOST ...
func (bot *Bot) VerifyPOST(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintln(w, "Internal server error")
		bot.Sugar.Errorf("Error parsing form: %v", err)
		return
	}

	if r.FormValue("uuid") == "" || r.FormValue("h-captcha-response") == "" {
		bot.ShowText(w, "Gatekeeper error", "", "One or more required values was empty.")
		return
	}

	id, err := uuid.Parse(r.FormValue("uuid"))
	if err != nil {
		bot.ShowText(w, "Gatekeeper error", "", "Invalid UUID provided.")
		return
	}

	u, err := bot.userByUUID(id)
	if err != nil {
		bot.ShowText(w, "Gatekeeper error", "", "Internal server error.")

		bot.Sugar.Errorf("Error getting UUID from database: %v", err)
		return
	}
	if !u.Pending {
		bot.ShowText(w, "Gatekeeper error", "", "User %v is not pending.", u.UserID)
		return
	}

	resp := bot.HCaptcha.SiteVerify(r)
	if !resp.Success {
		bot.ShowText(w, "Verification failed", "", "Verification failed.")
		fmt.Fprintln(w, "Verification failed.")
	}

	s, err := bot.serverSettings(u.GuildID)
	if err != nil {
		bot.ShowText(w, "Gatekeeper error", "", "Internal server error.")
		bot.Sugar.Errorf("Error getting server settings: %v", err)
	}

	if !s.MemberRole.IsValid() {
		// we shouldn't get here
		bot.ShowText(w, "Gatekeeper", "", "Done! You may now close this window.")
		bot.Sugar.Error("User passed verification, but no member role was set.")
		return
	}

	state, _ := bot.Router.StateFromGuildID(u.GuildID)

	err = state.AddRole(u.GuildID, u.UserID, s.MemberRole, api.AddRoleData{
		AuditLogReason: "Gatekeeper: add member role",
	})
	if err != nil {
		bot.ShowText(w, "Gatekeeper error", "", "There was an error adding your member role. Please contact a server administrator for help.")
		bot.Sugar.Errorf("Error adding role for %v in %v: %v", u.UserID, u.GuildID, err)
	}

	if s.WelcomeChannel.IsValid() && s.WelcomeMessage != "" {
		msg := strings.NewReplacer("{mention}", u.UserID.Mention()).Replace(s.WelcomeMessage)
		_, err = state.SendMessage(s.WelcomeChannel, msg)
		if err != nil {
			bot.Sugar.Errorf("Error sending welcome message: %v", err)
		}
	}

	err = bot.completeCaptcha(u.GuildID, u.UserID)
	if err != nil {
		bot.Sugar.Errorf("Error setting pending status for %v: %v", u.UserID, err)
	}

	bot.ShowText(w, "Gatekeeper", "", "Done! You may now close this window.")
}
