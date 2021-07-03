package gatekeeper

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/starshine-sys/tribble/types"
)

// PendingUser ...
type PendingUser struct {
	UserID   discord.UserID
	ServerID discord.GuildID

	Key     uuid.UUID
	Pending bool
}

type gatekeeperData struct {
	UUID   string
	Config *types.BotConfig
}

var gatekeeperTmpl = template.Must(template.ParseFiles("../../gatekeeper/tmpl.html"))

// GatekeeperGET ...
func (bot *Bot) GatekeeperGET(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := uuid.Parse(ps.ByName("uuid"))
	if err != nil {
		fmt.Fprintf(w, "Invalid UUID (%v) provided!\n", ps.ByName("uuid"))
		return
	}

	u, err := bot.userByUUID(id)
	if err != nil {
		fmt.Fprintln(w, "Internal server error")
		bot.Sugar.Errorf("Error getting UUID from database: %v", err)
		return
	}
	if !u.Pending {
		fmt.Fprintf(w, "User %v is not pending!\n", u.UserID)
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
	return
}

// VerifyPOST ...
func (bot *Bot) VerifyPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintln(w, "Internal server error")
		bot.Sugar.Errorf("Error parsing form: %v", err)
		return
	}

	if r.FormValue("uuid") == "" || r.FormValue("h-captcha-response") == "" {
		fmt.Fprintln(w, "One or more required values was empty!")
		return
	}

	id, err := uuid.Parse(r.FormValue("uuid"))
	if err != nil {
		fmt.Fprintln(w, "Invalid UUID provided!")
		return
	}

	u, err := bot.userByUUID(id)
	if err != nil {
		fmt.Fprintln(w, "Internal server error")
		bot.Sugar.Errorf("Error getting UUID from database: %v", err)
		return
	}
	if !u.Pending {
		fmt.Fprintf(w, "User %v is not pending!\n", u.UserID)
		return
	}

	resp := bot.HCaptcha.SiteVerify(r)
	if !resp.Success {
		fmt.Fprintln(w, "Verification failed.")
	}

	s, err := bot.serverSettings(u.ServerID)
	if err != nil {
		fmt.Fprintln(w, "Internal server error")
		bot.Sugar.Errorf("Error getting server settings: %v", err)
	}

	if !s.MemberRole.IsValid() {
		// we shouldn't get here
		fmt.Fprintln(w, "Done! You may now close this window.")
		bot.Sugar.Error("User passed verification, but no member role was set.")
		return
	}

	state, _ := bot.Router.StateFromGuildID(u.ServerID)

	err = state.AddRole(u.ServerID, u.UserID, s.MemberRole)
	if err != nil {
		fmt.Fprintln(w, "There was an error adding your member role. Please contact a server administrator for help.")
		bot.Sugar.Errorf("Error adding role for %v in %v: %v", u.UserID, u.ServerID, err)
	}

	if s.WelcomeChannel.IsValid() && s.WelcomeMessage != "" {
		msg := strings.NewReplacer("{mention}", u.UserID.Mention()).Replace(s.WelcomeMessage)
		_, err = state.SendMessage(s.WelcomeChannel, msg)
		if err != nil {
			bot.Sugar.Errorf("Error sending welcome message: %v", err)
		}
	}

	err = bot.completeCaptcha(u.ServerID, u.UserID)
	if err != nil {
		bot.Sugar.Errorf("Error setting pending status for %v: %v", u.UserID, err)
	}

	fmt.Fprintln(w, "Done! You may now close this window.")
}
