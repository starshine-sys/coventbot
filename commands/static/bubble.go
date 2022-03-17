package static

import (
	"fmt"
	"math/rand"

	"github.com/starshine-sys/bcr"
	bcr2 "github.com/starshine-sys/bcr/v2"
)

func (bot *Bot) bubbleSlash(ctx *bcr2.CommandContext) error {
	size, err := ctx.Option("size").IntValue()
	if err != nil || size == 0 {
		size = 10
	}

	prepop, _ := ctx.Option("prepop").BoolValue()
	ephemeral, _ := ctx.Option("ephemeral").BoolValue()

	if size > 13 {
		return ctx.ReplyEphemeral(
			fmt.Sprintf("A size of %v is too large, maximum 13.", size),
		)
	} else if size < 1 {
		return ctx.ReplyEphemeral(
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

	if ephemeral {
		return ctx.ReplyEphemeral(out)
	}
	return ctx.Reply(out)
}

func (bot *Bot) bubble(ctx *bcr.Context) (err error) {

	size, _ := ctx.Flags.GetInt("size")
	prepop, _ := ctx.Flags.GetBool("prepop")

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
	return ctx.SendX(out)
}
