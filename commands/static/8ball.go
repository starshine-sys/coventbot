package static

import (
	"crypto/md5"
	"encoding/binary"
	"math/rand"
	"time"

	"github.com/starshine-sys/bcr"
)

var eightballResponses = [...]string{
	"It is certain",
	"It is decidedly so",
	"Without a doubt",
	"Yes definitely",
	"You may rely on it",
	"As I see it, yes",
	"Most likely",
	"Outlook good",
	"Yes",
	"Signs point to yes",
	"Reply hazy, try again",
	"Ask again later",
	"Better not tell you now",
	"Cannot predict now",
	"Concentrate and ask again",
	"Don't count on it",
	"My reply is no",
	"My sources say no",
	"Outlook not so good",
	"Very doubtful",
}

func (bot *Bot) eightball(ctx *bcr.Context) (err error) {
	md5 := md5.Sum([]byte(ctx.RawArgs))
	hash := int64(binary.BigEndian.Uint64(md5[:]))
	t := time.Now().UnixNano()

	r := rand.New(rand.NewSource(t ^ hash))

	idx := r.Intn(len(eightballResponses))

	name := ctx.Author.Username
	if ctx.Member != nil && ctx.Member.Nick != "" {
		name = ctx.Member.Nick
	}

	return ctx.SendfX("🔮 %v, %v.", eightballResponses[idx], name)
}
