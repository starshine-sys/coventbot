package static

import (
	"fmt"
	"math/rand"

	"github.com/starshine-sys/bcr"
)

func (bot *Bot) bubbleSlash(ctx bcr.Contexter) (err error) {

	size := ctx.GetIntFlag("size")
	prepop := ctx.GetBoolFlag("prepop")
	ephemeral := ctx.GetBoolFlag("ephemeral")

	if size == 0 {
		size = 10
	}

	if size > 13 {
		return ctx.SendEphemeral(
			fmt.Sprintf("A size of %v is too large, maximum 13.", size),
		)
	} else if size < 1 {
		return ctx.SendEphemeral(
			fmt.Sprintf("A size of %v is too small, minimum 1.", size),
		)
	}

	var out string
	for i := size; i > 0; i-- {
		for j := size; j > 0; j-- {
			if prepop {
				if j != 1 && j != size && i != 1 && i != size {
					if rand.Intn(6) == 5 {
						out += "pop"
					} else {
						out += "||pop||"
					}
				} else {
					out += "||pop||"
				}
			} else {
				out += "||pop||"
			}
		}
		out += "\n"
	}
	if _, ok := ctx.(*bcr.SlashContext); ok && ephemeral {
		return ctx.SendEphemeral(out)
	}
	return ctx.SendX(out)
}
