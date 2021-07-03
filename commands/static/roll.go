package static

import (
	"math/rand"
	"regexp"
	"strconv"

	"github.com/starshine-sys/bcr"
)

var dice = regexp.MustCompile(`(\d*) ?[dD](\d+)`)

func (bot *Bot) roll(ctx *bcr.Context) (err error) {
	num, die := 1, 20
	var total int

	if dice.MatchString(ctx.RawArgs) {
		m := dice.FindStringSubmatch(ctx.RawArgs)
		numStr, dieStr := m[1], m[2]

		if numStr == "" {
			num = 1
		} else {
			num, err = strconv.Atoi(numStr)
			if err != nil {
				_, err = ctx.Send("Couldn't parse your input as standard dice notation.")
				return
			}
		}
		die, err = strconv.Atoi(dieStr)
		if err != nil {
			_, err = ctx.Send("Couldn't parse your input as standard dice notation.")
			return
		}

		if num > 100 {
			_, err = ctx.Send("Too many dice! (Maximum 100)")
			return
		}
		if die > 1000 {
			_, err = ctx.Send("Too large a die! (Maximum 1000)")
			return
		}
		if num <= 0 || die <= 0 {
			_, err = ctx.Send("Both the number of dice and the dice to use must be larger than 0.")
			return
		}

		var results []int
		for i := 0; i < num; i++ {
			r := rand.Intn(die) + 1

			total += r
			results = append(results, r)
		}

		_, err = ctx.Sendf("%vd%v: **%v** (%v)", num, die, total, results)
		return
	}

	if len(ctx.Args) > 0 {
		die, err = strconv.Atoi(ctx.RawArgs)
		if err != nil {
			_, err = ctx.Send("Couldn't parse your input as a number.")
			return
		}
	}

	_, err = ctx.Sendf("d%v: %v", die, rand.Intn(die)+1)
	return
}
