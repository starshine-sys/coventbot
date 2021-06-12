package levels

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"strings"

	// to decode JPG backgrounds
	_ "image/jpeg"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/disintegration/imaging"
	"github.com/dustin/go-humanize"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/starshine-sys/bcr"
)

// oh gods i hate this
// basically we iterate over this 4 times to crop the avatar into a circle
// oh gods this is so hacky
var blankPixels = []int{96, 96, 96, 96, 85, 85, 85, 85, 74, 74, 74, 74, 68, 68, 68, 62, 62, 62, 62, 55, 55, 55, 55, 50, 50, 50, 50, 45, 45, 45, 45, 39, 39, 39, 39, 39, 39, 33, 33, 33, 33, 33, 33, 28, 28, 28, 28, 28, 24, 24, 24, 24, 24, 24, 20, 20, 20, 20, 20, 20, 16, 16, 16, 16, 16, 16, 16, 12, 12, 12, 12, 12, 12, 12, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 4, 4, 4, 4, 4, 4, 4, 4}

//go:embed templates
var imageData embed.FS

func (bot *Bot) level(ctx *bcr.Context) (err error) {
	sc, err := bot.getGuildConfig(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}
	if !sc.LevelsEnabled {
		_, err = ctx.Send("Levels are not enabled on this server.", nil)
		return
	}

	var embed bool
	if len(ctx.Args) > 0 && strings.EqualFold(ctx.Args[0], "embed") {
		embed = true
		ctx.RawArgs = strings.TrimSpace(strings.TrimPrefix(ctx.RawArgs, ctx.Args[0]))
	}

	u := &ctx.Author
	if ctx.RawArgs != "" {
		u, err = ctx.ParseUser(ctx.RawArgs)
		if err != nil {
			_, err = ctx.Send("User not found.", nil)
			return
		}
	}

	uc, err := bot.getUser(ctx.Message.GuildID, u.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	lvl := currentLevel(uc.XP)
	xpForNext := expForNextLevel(lvl)
	xpForPrev := expForNextLevel(lvl - 1)

	// get leaderboard (for rank)
	// filter the leaderboard to match the `leaderboard` command
	var rank int
	lb, err := bot.getLeaderboard(ctx.Message.GuildID, false)
	if err == nil {
		for i, uc := range lb {
			if uc.UserID == u.ID {
				rank = i + 1
				break
			}
		}
	}

	// get user colour
	clr := uc.Colour
	if clr == 0 {
		clr, err = ctx.State.MemberColor(ctx.Message.GuildID, u.ID)
		if err != nil || clr == 0 {
			clr = bcr.ColourBlurple
		}
	}

	if embed {
		return bot.lvlEmbed(ctx, u, sc, uc, lvl, xpForNext, xpForPrev, rank, clr)
	}

	img := gg.NewContext(1200, 400)

	// background image
	if uc.Background != "" || sc.Background != "" {
		url := uc.Background
		if url == "" {
			url = sc.Background
		}

		resp, err := http.Get(url)
		if err == nil {
			defer resp.Body.Close()

			bg, _, err := image.Decode(resp.Body)
			if err == nil {
				bg = imaging.Resize(bg, 1200, 0, imaging.NearestNeighbor)
				img.DrawImageAnchored(bg, 0, 0, 0, 0)
			}
		}
	}

	// background
	img.SetHexColor("#00000088")
	img.DrawRoundedRectangle(50, 50, 1100, 300, 20)
	img.Fill()

	resp, err := http.Get(u.AvatarURLWithType(discord.PNGImage) + "?size=256")
	if err != nil {
		return bot.lvlEmbed(ctx, u, sc, uc, lvl, xpForNext, xpForPrev, rank, clr)
	}
	defer resp.Body.Close()

	pfp, err := png.Decode(resp.Body)
	if err != nil {
		return bot.lvlEmbed(ctx, u, sc, uc, lvl, xpForNext, xpForPrev, rank, clr)
	}

	pfp = imaging.Resize(pfp, 256, 256, imaging.NearestNeighbor)

	pfpImg := gg.NewContextForImage(pfp)
	pfpImg.SetColor(color.RGBA{0, 0, 0, 0})

	for y := 0; y < len(blankPixels); y++ {
		for x := 0; x < blankPixels[y]; x++ {
			pfpImg.SetPixel(x, y)
			pfpImg.SetPixel(256-x, 256-y)
			pfpImg.SetPixel(x, 256-y)
			pfpImg.SetPixel(256-x, y)
		}
	}

	img.DrawImageAnchored(pfpImg.Image(), 200, 200, 0.5, 0.5)

	// set font
	var f, bf *truetype.Font
	{
		fontBytes, err := imageData.ReadFile("templates/Montserrat-Regular.ttf")
		if err != nil {
			return bot.lvlEmbed(ctx, u, sc, uc, lvl, xpForNext, xpForPrev, rank, clr)
		}

		f, err = truetype.Parse(fontBytes)
		if err != nil {
			return bot.lvlEmbed(ctx, u, sc, uc, lvl, xpForNext, xpForPrev, rank, clr)
		}

		bfb, err := imageData.ReadFile("templates/Montserrat-Medium.ttf")
		if err != nil {
			return bot.lvlEmbed(ctx, u, sc, uc, lvl, xpForNext, xpForPrev, rank, clr)
		}

		bf, err = truetype.Parse(bfb)
		if err != nil {
			return bot.lvlEmbed(ctx, u, sc, uc, lvl, xpForNext, xpForPrev, rank, clr)
		}
	}

	img.SetHexColor(fmt.Sprintf("#%06x", clr))
	img.DrawCircle(200, 200, 130)

	r, g, b := clr.RGB()

	img.SetStrokeStyle(gg.NewSolidPattern(color.NRGBA{r, g, b, 0xFF}))
	img.SetLineWidth(5)
	img.Stroke()

	progress := uc.XP - xpForPrev
	needed := xpForNext - xpForPrev

	p := float64(progress) / float64(needed)

	end := 750 * p

	img.DrawRectangle(350, 275, end, 50)
	img.Fill()

	img.SetHexColor("#686868")
	img.DrawRectangle(350+end, 275, 750-end, 50)
	img.Fill()

	img.SetHexColor(fmt.Sprintf("#%06xCC", clr))

	img.SetColor(color.NRGBA{0xB5, 0xB5, 0xB5, 0xCC})

	img.DrawRectangle(350, 180, 750, 3)
	img.Fill()

	img.SetStrokeStyle(gg.NewSolidPattern(color.NRGBA{0xB5, 0xB5, 0xB5, 0xFF}))

	img.DrawRoundedRectangle(350, 275, 750, 50, 5)
	img.SetLineWidth(2)
	img.Stroke()

	img.SetHexColor("#ffffff")

	img.SetFontFace(truetype.NewFace(bf, &truetype.Options{
		Size: 60,
	}))

	name := ""

	for i, r := range u.Username {
		if i > 16 {
			name += string(r) + "..."
			break
		}

		name += string(r)
	}

	img.DrawStringAnchored(name, 350, 120, 0, 0.5)

	// rank/xp
	img.SetFontFace(truetype.NewFace(f, &truetype.Options{
		Size: 40,
	}))

	img.DrawStringAnchored(fmt.Sprintf("Rank #%v", rank), 1100, 120, 1, 0.5)

	img.DrawStringAnchored(fmt.Sprintf("Level %v", lvl), 1100, 200, 1, 1)

	img.DrawStringAnchored(fmt.Sprintf("%v%%", int64(p*100)), 725, 295, 0.5, 0.5)

	progressStr := fmt.Sprintf("%v/%v XP", humanize.Comma(progress), humanize.Comma(needed))

	img.DrawStringAnchored(progressStr, 350, 200, 0, 1)

	buf := new(bytes.Buffer)

	err = img.EncodePNG(buf)
	if err != nil {
		return bot.lvlEmbed(ctx, u, sc, uc, lvl, xpForNext, xpForPrev, rank, clr)
	}

	_, err = ctx.NewMessage().AddFile("level_card.png", buf).Send()
	return
}

func (bot *Bot) lvlEmbed(ctx *bcr.Context, u *discord.User, sc Server, uc Levels, lvl, xpForNext, xpForPrev int64, rank int, clr discord.Color) (err error) {
	e := discord.Embed{
		Thumbnail: &discord.EmbedThumbnail{
			URL: u.AvatarURLWithType(discord.PNGImage),
		},
		Title: fmt.Sprintf("Level %v - Rank #%v", lvl, rank),
		Description: fmt.Sprintf(
			"%v/%v XP", humanize.Comma(uc.XP), humanize.Comma(xpForNext),
		),
		Color: clr,
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("%v#%v", u.Username, u.Discriminator),
		},
		Timestamp: discord.NowTimestamp(),
	}

	{
		progress := uc.XP - xpForPrev
		needed := xpForNext - xpForPrev

		p := float64(progress) / float64(needed)

		percent := int64(p * 100)

		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Progress to next level",
			Value: fmt.Sprintf("%v%% (%v/%v XP)", percent, humanize.Comma(progress), humanize.Comma(needed)),
		})
	}

	// get next reward
	reward := bot.getNextReward(ctx.Message.GuildID, lvl)
	if reward != nil && sc.ShowNextReward {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Next reward",
			Value: fmt.Sprintf("%v\nat level %v", reward.RoleReward.Mention(), reward.Lvl),
		})
	} else if sc.ShowNextReward {
		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Next reward",
			Value: "No more rewards",
		})
	}

	_, err = ctx.Send("", &e)
	return
}
