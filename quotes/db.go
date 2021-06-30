package quotes

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"emperror.dev/errors"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/pkgo"
)

// Quote is a single quote
type Quote struct {
	ID  int64
	HID string

	ServerID  discord.GuildID
	ChannelID discord.ChannelID
	MessageID discord.MessageID
	UserID    discord.UserID
	AddedBy   discord.UserID
	Content   string
	Proxied   bool

	Added time.Time
}

// Embed creates an embed for the quote
func (q Quote) Embed(s *pkgo.Session) discord.Embed {
	e := discord.Embed{
		Footer: &discord.EmbedFooter{
			Text: "ID: " + q.HID,
		},
		Timestamp: discord.NewTimestamp(q.Added),
		Color:     bcr.ColourBlurple,
	}

	mention := "<@!" + q.UserID.String() + ">"
	if q.Proxied {
		msg, err := s.Message(pkgo.Snowflake(q.MessageID))
		if err == nil {
			mention = msg.Member.Name
			if msg.System.Tag != "" {
				mention += " " + msg.System.Tag
			}
		}
	}

	author := fmt.Sprintf("- %v [(Jump)](https://discord.com/channels/%v/%v/%v)", mention, q.ServerID, q.ChannelID, q.MessageID)

	content := q.Content
	if len(content) > 2020-len(author) {
		content = content[:2020-len(author)]
	}

	e.Description = content + "\n" + author

	return e
}

var errNotExists = errors.New("quote with given ID does not exist")

func (bot *Bot) quotesEnabled(guildID discord.GuildID) (b bool) {
	bot.DB.Pool.QueryRow(context.Background(), "select quotes_enabled from servers where id = $1", guildID).Scan(&b)
	return
}

func (bot *Bot) suppressMessages(guildID discord.GuildID) (b bool) {
	bot.DB.Pool.QueryRow(context.Background(), "select quote_suppress_messages from servers where id = $1", guildID).Scan(&b)
	return
}

func (bot *Bot) insertQuote(q Quote) (Quote, error) {
	out, err := Encrypt([]byte(q.Content), bot.AESKey)
	if err != nil {
		return q, err
	}
	q.Content = hex.EncodeToString(out)

	err = pgxscan.Get(context.Background(), bot.DB.Pool, &q, `insert into quotes
(hid, server_id, channel_id, message_id, user_id, added_by, content, proxied)
values (find_free_quote_hid($1), $1, $2, $3, $4, $5, $6, $7) returning *`, q.ServerID, q.ChannelID, q.MessageID, q.UserID, q.AddedBy, q.Content, q.Proxied)
	return q, err
}

func (bot *Bot) getQuote(hid string, guildID discord.GuildID) (q Quote, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &q, "select * from quotes where hid ilike $1 and server_id = $2", hid, guildID)
	if err != nil {
		if errors.Cause(err) == pgx.ErrNoRows {
			return q, errNotExists
		}
		return q, err
	}

	b, err := hex.DecodeString(q.Content)
	if err != nil {
		return q, err
	}

	out, err := Decrypt(b, bot.AESKey)
	if err != nil {
		return q, err
	}

	q.Content = string(out)
	return
}

func (bot *Bot) quoteMessage(id discord.MessageID) (q Quote, err error) {
	err = pgxscan.Get(context.Background(), bot.DB.Pool, &q, "select * from quotes where message_id = $1", id)
	if err != nil {
		return q, err
	}

	b, err := hex.DecodeString(q.Content)
	if err != nil {
		return q, err
	}

	out, err := Decrypt(b, bot.AESKey)
	if err != nil {
		return q, err
	}

	q.Content = string(out)
	return
}

func (bot *Bot) serverQuote(guildID discord.GuildID) (q *Quote, err error) {
	ids := []string{}
	err = bot.DB.Pool.QueryRow(context.Background(), "select array(select hid from quotes where server_id = $1)", guildID).Scan(&ids)
	if err != nil {
		return
	}

	if len(ids) == 1 {
		quote, err := bot.getQuote(ids[0], guildID)
		return &quote, err
	}

	if len(ids) == 0 {
		return nil, pgx.ErrNoRows
	}

	n := rand.Intn(len(ids))

	quote, err := bot.getQuote(ids[n], guildID)
	return &quote, err
}

func (bot *Bot) userQuote(guildID discord.GuildID, userID discord.UserID) (q *Quote, err error) {
	ids := []string{}
	err = bot.DB.Pool.QueryRow(context.Background(), "select array(select hid from quotes where server_id = $1 and user_id = $2)", guildID, userID).Scan(&ids)
	if err != nil {
		return
	}

	if len(ids) == 1 {
		quote, err := bot.getQuote(ids[0], guildID)
		return &quote, err
	}

	if len(ids) == 0 {
		return nil, pgx.ErrNoRows
	}

	n := rand.Intn(len(ids))

	quote, err := bot.getQuote(ids[n], guildID)
	return &quote, err
}

func (bot *Bot) delQuote(guildID discord.GuildID, hid string) (err error) {
	_, err = bot.DB.Pool.Exec(context.Background(), "delete from quotes where server_id = $1 and hid = $2", guildID, hid)
	return
}

func (bot *Bot) quotes(guildID discord.GuildID) (quotes []Quote, err error) {
	err = pgxscan.Select(context.Background(), bot.DB.Pool, &quotes, "select id, hid, server_id, channel_id, message_id, user_id, added_by, added, proxied from quotes where server_id = $1 order by added", guildID)
	return
}
