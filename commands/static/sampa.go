package static

import (
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/starshine-sys/bcr"
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

	msg, err := ctx.Send("", &embed)
	if err != nil {
		return err
	}

	ctx.AddReactionHandlerWithTimeout(msg.ID, ctx.Author.ID, "❌", true, false, 2*time.Hour, func(ctx *bcr.Context) {
		err := bot.State.DeleteMessage(msg.ChannelID, msg.ID)
		if err != nil {
			bot.Sugar.Errorf("Error deleting message: %v", err)
		}
	})
	return nil
}
