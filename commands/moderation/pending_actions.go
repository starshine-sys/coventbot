package moderation

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/georgysavva/scany/pgxscan"
)

// Pending actions
const (
	PendingUnban   = "unban"
	PendingUnmute  = "unmute"
	PendingUnpause = "unpause"
)

// PendingAction is a pending action
type PendingAction struct {
	ID      int64
	GuildID discord.GuildID
	UserID  discord.UserID
	Expires time.Time

	Type   string
	Log    bool
	Reason string
}

func (bot *Bot) insertPendingAction(guildID discord.GuildID, userID discord.UserID, expires time.Time, t string, log bool, reason string) (pa PendingAction, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &pa, "insert into pending_actions (guild_id, user_id, expires, type, log, reason) values ($1, $2, $3, $4, $5, $6) returning *", guildID, userID, expires.UTC(), t, log, reason)
	return
}

func (bot *Bot) pendingActions(limit int) (actions []PendingAction, err error) {
	err = pgxscan.Select(context.Background(), bot.DB.Pool, &actions, "select * from pending_actions where expires < $1 limit $2", time.Now().UTC(), limit)
	return
}

func (bot *Bot) deleteAction(id int64) (n int64, err error) {
	ct, err := bot.DB.Pool.Exec(context.Background(), "delete from pending_actions where id = $1", id)
	return ct.RowsAffected(), err
}

func (bot *Bot) actionsFor(guildID discord.GuildID, userID discord.UserID) (actions []PendingAction, err error) {
	err = pgxscan.Select(context.Background(), bot.DB.Pool, &actions, "select * from pending_actions where guild_id = $1 and user_id = $2", guildID, userID)
	return
}

func (bot *Bot) doPendingActions(s *state.State) {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case <-sc:
			break
		default:
		}

		actions, err := bot.pendingActions(5)
		if err != nil {
			bot.Sugar.Errorf("Error fetching pending actions: %v", err)
			time.Sleep(time.Minute)
			continue
		}

		for _, pa := range actions {
			switch pa.Type {
			case PendingUnban:
				err = s.Unban(pa.GuildID, pa.UserID, api.AuditLogReason(pa.Reason))
				if err != nil {
					bot.Sugar.Errorf("Error unbanning user: %v", err)
					bot.deleteAction(pa.ID)
					continue
				}

				if pa.Log {
					err = bot.ModLog.Unban(s, pa.GuildID, pa.UserID, bot.Router.Bot.ID, pa.Reason)
					if err != nil {
						bot.Sugar.Errorf("Error sending log message: %v", err)
					}
				}
			case PendingUnmute, PendingUnpause:
				roles, err := bot.muteRoles(pa.GuildID)
				if err != nil {
					bot.Sugar.Errorf("Error getting mute roles for %v: %v", pa.GuildID, err)
					bot.deleteAction(pa.ID)
					continue
				}

				if pa.Type == PendingUnmute {
					if !roles.MuteRole.IsValid() {
						bot.deleteAction(pa.ID)
						continue
					}

					err = s.RemoveRole(pa.GuildID, pa.UserID, roles.MuteRole, api.AuditLogReason(pa.Reason))
					if err != nil {
						bot.Sugar.Errorf("Error unmuting user: %v", err)
						bot.deleteAction(pa.ID)
						continue
					}

					if pa.Log {
						err = bot.ModLog.Unmute(s, pa.GuildID, pa.UserID, bot.Router.Bot.ID, pa.Reason)
						if err != nil {
							bot.Sugar.Errorf("Error sending log message: %v", err)
						}
					}
				} else {
					if !roles.PauseRole.IsValid() {
						bot.deleteAction(pa.ID)
						continue
					}

					err = s.RemoveRole(pa.GuildID, pa.UserID, roles.PauseRole, api.AuditLogReason(pa.Reason))
					if err != nil {
						bot.Sugar.Errorf("Error unpausing user: %v", err)
						bot.deleteAction(pa.ID)
						continue
					}

					if pa.Log {
						err = bot.ModLog.Unpause(s, pa.GuildID, pa.UserID, bot.Router.Bot.ID, pa.Reason)
						if err != nil {
							bot.Sugar.Errorf("Error sending log message: %v", err)
						}
					}
				}
			}

			bot.deleteAction(pa.ID)
		}

		time.Sleep(time.Second)
	}
}

func (bot *Bot) muteOnJoin(ev *gateway.GuildMemberAddEvent) {
	actions, err := bot.actionsFor(ev.GuildID, ev.User.ID)
	if err != nil {
		bot.Sugar.Errorf("Error getting pending actions for %v: %v", ev.GuildID, err)
		return
	}

	roles, err := bot.muteRoles(ev.GuildID)
	if err != nil {
		bot.Sugar.Errorf("Error getting mute roles for %v: %v", ev.GuildID, err)
		return
	}

	s, _ := bot.Router.StateFromGuildID(ev.GuildID)

	for _, pa := range actions {
		// cancel the unban as they're obviously already unbanned
		if pa.Type == PendingUnban {
			bot.deleteAction(pa.ID)
		}

		switch pa.Type {
		case PendingUnmute:
			if !roles.MuteRole.IsValid() {
				continue
			}

			s.AddRole(ev.GuildID, ev.User.ID, roles.MuteRole, api.AddRoleData{
				AuditLogReason: "Member was previously muted",
			})
		case PendingUnpause:
			if !roles.PauseRole.IsValid() {
				continue
			}

			s.AddRole(ev.GuildID, ev.User.ID, roles.PauseRole, api.AddRoleData{
				AuditLogReason: "Member was previously paused",
			})
		}
	}
}
