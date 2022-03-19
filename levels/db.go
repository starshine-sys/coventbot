package levels

import (
	"context"
	"math/rand"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/georgysavva/scany/pgxscan"
)

// LevelMessages ...
type LevelMessages string

// ...
const (
	AllDM          LevelMessages = "ALL_DM"
	RewardsDM      LevelMessages = "REWARDS_DM"
	AllChannel     LevelMessages = "ALL_CHANNEL"
	RewardsChannel LevelMessages = "REWARDS_CHANNEL"
	NoMessages     LevelMessages = "NONE"
)

// Server ...
type Server struct {
	ID discord.GuildID

	BlockedChannels   []uint64
	BlockedRoles      []uint64
	BlockedCategories []uint64

	RewardLog   discord.ChannelID
	NolevelsLog discord.ChannelID

	BetweenXP  time.Duration
	RewardText string
	Background string

	LevelsEnabled      bool
	LeaderboardModOnly bool
	ShowNextReward     bool
	Voice              bool

	LevelMessages LevelMessages
	LevelChannel  discord.ChannelID
}

// Levels ...
type Levels struct {
	ServerID discord.GuildID
	UserID   discord.UserID

	XP int64

	Colour     discord.Color
	Background string

	NextTime time.Time
}

// Reward ...
type Reward struct {
	ServerID   discord.GuildID
	Lvl        int64
	RoleReward discord.RoleID
}

func (bot *Bot) getGuildConfig(guildID discord.GuildID) (s Server, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &s, "insert into server_levels (id) values ($1) on conflict (id) do update set id = $1 returning *", guildID)
	return s, err
}

func (bot *Bot) getUser(guildID discord.GuildID, userID discord.UserID) (l Levels, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &l, "insert into levels (server_id, user_id) values ($1, $2) on conflict (server_id, user_id) do update set server_id = $1 returning *", guildID, userID)
	return l, err
}

func (bot *Bot) incrementXP(guildID discord.GuildID, userID discord.UserID) (newXP int64, err error) {
	xp := 15 + int64(rand.Intn(12))

	err = bot.DB.Pool.QueryRow(context.Background(), "update levels set xp = xp + $4, next_time = $3 where server_id = $1 and user_id = $2 returning xp", guildID, userID, time.Now().UTC(), xp).Scan(&newXP)
	return
}

func (bot *Bot) getReward(guildID discord.GuildID, lvl int64) *Reward {
	r := Reward{}

	var exists bool
	err := bot.DB.Pool.QueryRow(context.Background(), "select exists(select * from level_rewards where server_id = $1 and lvl = $2)", guildID, lvl).Scan(&exists)
	if err != nil {
		bot.Sugar.Errorf("Error checking if guild exists in db: %v", err)
	}
	if !exists {
		return nil
	}

	err = pgxscan.Get(context.Background(), bot.DB.Pool, &r, "select * from level_rewards where server_id = $1 and lvl = $2", guildID, lvl)
	if err != nil {
		bot.Sugar.Errorf("Error getting reward: %v", err)
		return nil
	}

	return &r
}

func (bot *Bot) getNextReward(guildID discord.GuildID, lvl int64) *Reward {
	r := Reward{}

	var exists bool
	err := bot.DB.Pool.QueryRow(context.Background(), "select exists(select * from level_rewards where server_id = $1 and lvl > $2)", guildID, lvl).Scan(&exists)
	if err != nil {
		bot.Sugar.Errorf("Error checking if guild exists in db: %v", err)
	}
	if !exists {
		return nil
	}

	err = pgxscan.Get(context.Background(), bot.DB.Pool, &r, "select * from level_rewards where server_id = $1 and lvl > $2 order by lvl asc limit 1", guildID, lvl)
	if err != nil {
		bot.Sugar.Errorf("Error getting reward: %v", err)
		return nil
	}

	return &r
}

func (bot *Bot) addReward(guildID discord.GuildID, lvl int64, roleID discord.RoleID) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "insert into level_rewards (server_id, lvl, role_reward) values ($1, $2, $3) on conflict (server_id, lvl) do update set role_reward = $3", guildID, lvl, roleID)
	return
}

func (bot *Bot) getAllRewards(guildID discord.GuildID) (rwds []Reward, err error) {
	err = pgxscan.Select(context.Background(), bot.DB.Pool, &rwds, "select * from level_rewards where server_id = $1 order by lvl asc", guildID)
	return
}

func (bot *Bot) getLeaderboard(guildID discord.GuildID, full bool) (lb []Levels, err error) {
	err = pgxscan.Select(context.Background(), bot.DB.Pool, &lb, "select * from levels where server_id = $1 order by xp desc, user_id asc", guildID)
	if err != nil || full {
		return
	}

	gm := bot.Members(guildID)

	filtered := []Levels{}

	for _, l := range lb {
		for _, m := range gm {
			if m.User.ID == l.UserID {
				filtered = append(filtered, l)
				break
			}
		}
	}

	return filtered, nil
}

func (bot *Bot) setNolevels(guildID discord.GuildID, userID discord.UserID, expires bool, expiry time.Time) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), `insert into nolevels
	(server_id, user_id, expires, expiry) values ($1, $2, $3, $4)
	on conflict (server_id, user_id) do update
	set expires = $3, expiry = $4`, guildID, userID, expires, expiry)
	return
}

// Nolevels ...
type Nolevels struct {
	ServerID discord.GuildID
	UserID   discord.UserID
	Expires  bool
	Expiry   time.Time

	LogChannel discord.ChannelID
}

func (bot *Bot) guildNolevels(guildID discord.GuildID) (list []Nolevels, err error) {
	err = pgxscan.Select(context.Background(), bot.DB.Pool, &list, "select * from nolevels where server_id = $1 order by user_id", guildID)
	return
}

func (bot *Bot) isBlacklisted(guildID discord.GuildID, userID discord.UserID) (blacklisted bool) {
	err := bot.DB.Pool.QueryRow(context.Background(), "select exists(select user_id from nolevels where server_id = $1 and user_id = $2)", guildID, userID).Scan(&blacklisted)
	if err != nil {
		bot.Sugar.Errorf("Error checking if user is blacklisted: %v", err)
	}
	return blacklisted
}
