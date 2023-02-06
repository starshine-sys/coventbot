// SPDX-License-Identifier: AGPL-3.0-only
package static

import (
	"context"
	"strings"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/starshine-sys/bcr"
	bcr2 "github.com/starshine-sys/bcr/v2"
)

var re = strings.NewReplacer(
	`b_<`, `ɓ`,
	"d`", `ɖ`,
	`d_<`, `ɗ`,
	`g_<`, `ɠ`,
	`h\`, `ɦ`,
	`j\`, `ʝ`,
	"l`", "ɭ",
	`l\`, `ɺ`,
	"n`", "ɳ",
	`p\`, `ɸ`,
	"r`", "ɽ",
	`r\`, "ɹ",
	"r\\`", "ɻ",
	"s`", "ʂ",
	`s\`, "ɕ",
	"t`", "ʈ",
	`v\`, `ʋ`,
	`x\`, `ɧ`,
	"z`", "ʐ",
	`z\`, "ʑ",

	`_A`, `̘`,
	"A", "ɑ",
	`B\`, `ʙ`,
	"B", "β",
	`C`, `ç`,
	`D`, `ð`,
	`E`, `ɛ`,
	`<F>`, `↘︎`,
	`F`, `ɱ`,
	`_G`, `ˠ`,
	`G\_<`, `ʛ`,
	`G\`, `ɢ`,
	`G`, `ɣ`,
	`H\`, `ʜ`,
	`H`, `ɥ`,
	`I\`, `ᵻ`,
	`I`, `ɪ`,
	`J\_<`, `ʄ`,
	`J\`, `ɟ`,
	`J`, `ɲ`,
	`K\`, `ɮ`,
	`K`, `ɬ`,
	`L\`, `ʟ`,
	`L`, `ʎ`,
	`M\`, `ɰ`,
	`M`, `ɯ`,
	`N\`, `ɴ`,
	`_N`, `̼`,
	`N`, `ŋ`,
	`O\`, `ʘ`,
	`_O`, `̹`,
	`O`, `ɔ`,
	`P`, `ʋ`,
	`Q`, `ɒ`,
	`R\`, `ʀ`,
	`<R>`, `↗︎`,
	`R`, `ʁ`,
	`S`, `ʃ`,
	`T`, `θ`,
	`U\`, `ᵿ`,
	`U`, `ʊ`,
	`V`, `ʌ`,
	`W`, `ʍ`,
	`X\`, `ħ`,
	`_X`, `̆`,
	`X`, `χ`,
	`Y`, `ʏ`,
	`Z`, `ʒ`,
	`"`, `ˈ`,
	`%`, `ˌ`,
	`'`, `ʲ`,
	`:\`, `ˑ`,
	`:`, `ː`,

	"@`", `ɚ`,
	`@\`, `ɘ`,
	`@`, `ə`,
	`{`, `æ`,
	`}`, `ʉ`,
	`1`, `ɨ`,
	`2`, `ø`,
	`3\`, `ɞ`,
	`3`, `ɜ`,
	`4`, `ɾ`,
	`5`, `ɫ`,
	`6`, `ɐ`,
	`7`, `ɤ`,
	`8`, `ɵ`,
	`9`, `œ`,
	`&`, `ɶ`,
	`?\`, `ʕ`,
	`?`, `ʔ`,
	`<\`, `ʢ`,
	`>\`, `ʡ`,
	`^`, `ꜛ`,
	`!\`, `ǃ`,
	`!`, `ꜜ`,
	`|\|\`, `ǁ`,
	`||`, `‖`,
	`|\`, `ǀ`,
	`|`, `|`,
	`=\`, `ǂ`,

	// diacritics
	`_"`, `̈`,
	`_+`, `̟`,
	`_-`, `̠`,
	`_0`, `̥`,
	`_=`, `̩`,
	`=`, `̩`,
	`_>`, `ʼ`,
	`_?\`, `ˤ`,
	`_^`, `̯`,
	`_}`, `̚`,
	"`", `˞`,
	`_~`, `̃`,
	`~`, `̃`,
	`_a`, `̺`,
	`_c`, `̜`,
	`_d`, `̪`,
	`_h`, `ʰ`,
	`_j`, `ʲ`,
	`_l`, `ˡ`,
	`_m`, `̻`,
	`_n`, `ⁿ`,
	`_o`, `̞`,
	`_q`, `̙`,
	`_r`, `̝`,
	`_t`, `̤`,
	`_v`, `̬`,
	`_w`, `ʷ`,
	`_x`, `̽`,

	// joiners
	`_`, `͡`,
	`-\`, `͜`,

	// separator
	`-`, ``,
)

func (bot *Bot) sampa(ctx *bcr.Context) (err error) {
	if ctx.RawArgs == "" {
		_, err = ctx.Replyc(bcr.ColourRed, "You must actually give something to convert!")
		return
	}

	embed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Icon: ctx.Author.AvatarURLWithType(discord.PNGImage),
			Name: ctx.DisplayName(),
		},
		Description: re.Replace(ctx.RawArgs),
		Color:       bcr.ColourBlurple,
	}

	msg, err := ctx.Send("", embed)
	if err != nil {
		return err
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "insert into command_responses (message_id, user_id) values ($1, $2)", msg.ID, ctx.Author.ID)
	return
}

func (bot *Bot) sampaSlash(ctx *bcr2.CommandContext) (err error) {
	text := ctx.Option("text").String()

	err = ctx.ReplyComplex(api.InteractionResponseData{
		Content: option.NewNullableString(re.Replace(text)),
		AllowedMentions: &api.AllowedMentions{
			Parse: []api.AllowedMentionType{},
		},
	})
	if err != nil {
		return
	}

	msg, err := ctx.Original()
	if err != nil {
		return err
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "insert into command_responses (message_id, user_id) values ($1, $2)", msg.ID, ctx.User.ID)
	return
}

func (bot *Bot) sampaReaction(ev *gateway.MessageReactionAddEvent) {
	if ev.Emoji.Name != "❌" && ev.Emoji.Name != "🗑️" {
		return
	}

	var matchUser discord.UserID

	err := bot.DB.Pool.QueryRow(context.Background(), "select user_id from command_responses where message_id = $1", ev.MessageID).Scan(&matchUser)
	if err != nil {
		return
	}

	if matchUser != ev.UserID {
		return
	}

	s, _ := bot.Router.StateFromGuildID(ev.GuildID)

	err = s.DeleteMessage(ev.ChannelID, ev.MessageID, "")
	if err != nil {
		bot.Sugar.Errorf("Error deleting message: %v", err)
		return
	}

	_, err = bot.DB.Pool.Exec(context.Background(), "delete from command_responses where message_id = $1", ev.MessageID)
	if err != nil {
		bot.Sugar.Errorf("Error deleting message: %v", err)
	}
}
