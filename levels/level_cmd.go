// SPDX-License-Identifier: AGPL-3.0-only
package levels

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"strconv"
	"strings"

	// to decode JPG backgrounds
	_ "image/jpeg"

	"github.com/AndreKR/multiface"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"github.com/disintegration/imaging"
	"github.com/dustin/go-humanize"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/starshine-sys/bcr"
	"github.com/starshine-sys/tribble/etc"
	"golang.org/x/image/font"
)

//go:embed templates
var imageData embed.FS

var normalFont font.Face

var montserrat, noto, emoji *truetype.Font

func mustParse(path string) *truetype.Font {
	b, err := imageData.ReadFile(path)
	if err != nil {
		panic(err)
	}

	f, err := truetype.Parse(b)
	if err != nil {
		panic(err)
	}

	return f
}

const defaultBoldSize = 60

func init() {
	// montserrat for most latin letters
	montserrat = mustParse("templates/Montserrat-Medium.ttf")
	// noto as fallback for other characters
	noto = mustParse("templates/NotoSans-Medium.ttf")
	// emoji fallback
	emoji = mustParse("templates/NotoEmoji-Regular.ttf")

	normalFont = truetype.NewFace(
		mustParse("templates/Montserrat-Regular.ttf"),
		&truetype.Options{
			Size: 40,
		},
	)
}

func boldFontSize(size float64) *multiface.Face {
	mf := &multiface.Face{}

	mf.AddTruetypeFace(truetype.NewFace(montserrat, &truetype.Options{
		Size: size,
	}), montserrat)

	// add noto font
	mf.AddTruetypeFace(truetype.NewFace(noto, &truetype.Options{
		Size: size,
	}), noto)

	// add noto emoji
	mf.AddTruetypeFace(truetype.NewFace(emoji, &truetype.Options{
		Size: size,
	}), emoji)

	return mf
}

const (
	width          = 1200
	height         = 400
	progressBarLen = width - 450
)

func (bot *Bot) level(ctx *bcr.Context) (err error) {
	sc, err := bot.getGuildConfig(ctx.Message.GuildID)
	if err != nil {
		return bot.Report(ctx, err)
	}
	if !sc.LevelsEnabled {
		_, err = ctx.Send("Levels are not enabled on this server.")
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
			_, err = ctx.Send("User not found.")
			return
		}
	}

	uc, err := bot.getUser(ctx.Message.GuildID, u.ID)
	if err != nil {
		return bot.Report(ctx, err)
	}

	lvl := sc.CalculateLevel(uc.XP)
	xpForNext := sc.CalculateExp(lvl + 1)
	xpForPrev := sc.CalculateExp(lvl)

	// get leaderboard (for rank)
	// filter the leaderboard to match the `leaderboard` command
	var rank int
	noRanks, err := bot.DB.GuildBoolGet(ctx.Message.GuildID, "levels:disable_ranks")
	if err != nil {
		return bot.Report(ctx, err)
	}
	if !noRanks {
		lb, err := bot.getLeaderboard(ctx.Message.GuildID, false)
		if err == nil {
			for i, uc := range lb {
				if uc.UserID == u.ID {
					rank = i + 1
					break
				}
			}
		}
	}

	// get user colour + avatar URL
	clr := uc.Colour
	avatarURL := u.AvatarURLWithType(discord.PNGImage) + "?size=256"
	username := u.Username
	if ctx.Guild != nil {
		m, err := bot.Member(ctx.Guild.ID, u.ID)
		if err == nil {
			if clr == 0 {
				clr, _ = discord.MemberColor(*ctx.Guild, m)
			}
			if m.Avatar != "" {
				avatarURL = m.AvatarURLWithType(discord.PNGImage, ctx.Message.GuildID) + "?size=256"
			}
			if m.Nick != "" {
				username = m.Nick
			}
		}
	}

	if embed {
		return ctx.SendX("", bot.generateEmbed(username, avatarURL, clr, rank, lvl, uc.XP, xpForNext, xpForPrev, sc))
	}

	// background image
	background := ""
	if sc.Background != "" {
		background = sc.Background
	}
	if uc.Background != "" {
		background = uc.Background
	}

	r, err := bot.generateImage(username, avatarURL, background, clr, rank, lvl, uc.XP, xpForNext, xpForPrev)
	if err != nil {
		bot.Sugar.Errorf("Error generating level card: %v", err)
		return ctx.SendX("", bot.generateEmbed(username, avatarURL, clr, rank, lvl, uc.XP, xpForNext, xpForPrev, sc))
	}

	return ctx.SendFiles("", sendpart.File{
		Name:   "level_card.png",
		Reader: r,
	})
}

func (bot *Bot) generateImage(
	name, avatarURL, backgroundURL string, clr discord.Color,
	rank int, lvl, xp, xpForNext, xpForPrev int64,
) (
	r io.Reader, err error,
) {
	img := gg.NewContext(width, height)

	// background
	if backgroundURL != "" {
		resp, err := http.Get(backgroundURL)
		if err == nil {
			defer resp.Body.Close()

			bg, _, err := image.Decode(resp.Body)
			if err == nil {
				bg = imaging.Resize(bg, 1200, 0, imaging.NearestNeighbor)
				img.DrawImageAnchored(bg, 0, 0, 0, 0)
			}
		}
	}

	img.SetHexColor("#00000088")
	img.DrawRoundedRectangle(50, 50, width-100, height-100, 20)
	img.Fill()

	// fetch avatar
	resp, err := http.Get(avatarURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// decode avatar
	pfp, err := png.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	// use average of avatar if the user has no colour
	if clr == 0 {
		r, g, b, _ := etc.AverageColour(pfp)

		clr = discord.Color(r)<<16 + discord.Color(g)<<8 + discord.Color(b)
	}

	// resize pfp to fit + crop to circle (shoddily)
	pfp = imaging.Resize(pfp, 256, 256, imaging.NearestNeighbor)

	pfpImg := gg.NewContext(256, 256)
	pfpImg.DrawCircle(128, 128, 128)
	pfpImg.Clip()
	pfpImg.DrawImage(pfp, 0, 0)

	// draw pfp to context
	img.SetHexColor(clr.String())
	img.DrawCircle(200, 200, 130)
	img.FillPreserve()

	img.DrawImageAnchored(pfpImg.Image(), 200, 200, 0.5, 0.5)

	img.SetLineWidth(5)
	img.Stroke()

	progress := xp - xpForPrev
	needed := xpForNext - xpForPrev

	p := float64(progress) / float64(needed)

	end := progressBarLen * p

	img.DrawRectangle(350, 275, end, 50)
	img.Fill()

	img.SetHexColor("#686868")
	img.DrawRectangle(350+end, 275, progressBarLen-end, 50)
	img.Fill()

	img.SetHexColor(clr.String())

	img.SetColor(color.NRGBA{0xB5, 0xB5, 0xB5, 0xCC})

	img.DrawRectangle(350, 180, progressBarLen, 3)
	img.Fill()

	img.SetStrokeStyle(gg.NewSolidPattern(color.NRGBA{0xB5, 0xB5, 0xB5, 0xFF}))

	img.DrawRoundedRectangle(350, 275, progressBarLen, 50, 5)
	img.SetLineWidth(2)
	img.Stroke()

	img.SetHexColor("#ffffff")

	currentSize := float64(defaultBoldSize)
	img.SetFontFace(boldFontSize(currentSize))

	targetLen := (width - 100) - 350
	if rank != 0 {
		targetLen -= 200
	}

	for currentSize > 1 {
		w, h := img.MeasureString(name)

		if w < float64(targetLen) {
			bot.Sugar.Debugf("name %q fits in %v height", name, h)

			break
		}

		currentSize -= 1
		img.SetFontFace(boldFontSize(currentSize))
	}

	// name
	img.DrawStringAnchored(name, 350, 120, 0, 0.5)

	// rank/xp
	img.SetFontFace(normalFont)

	if rank != 0 {
		img.DrawStringAnchored(fmt.Sprintf("Rank #%v", rank), width-100, 120, 1, 0.5)
	}

	img.DrawStringAnchored(fmt.Sprintf("%v%%", int64(p*100)), 350+(progressBarLen/2), 295, 0.5, 0.5)

	progressStr := fmt.Sprintf("%v/%v", HumanizeInt64(progress), HumanizeInt64(needed))

	img.DrawStringAnchored(fmt.Sprintf("Level %v", lvl), 350, 210, 0, 1)
	img.DrawStringAnchored(progressStr, width-100, 210, 1, 1)

	buf := new(bytes.Buffer)

	err = img.EncodePNG(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (bot *Bot) generateEmbed(
	name, avatarURL string, clr discord.Color,
	rank int, lvl, xp, xpForNext, xpForPrev int64,
	sc Server,
) discord.Embed {

	e := discord.Embed{
		Color:       clr,
		Title:       fmt.Sprintf("Level %v - Rank #%v", lvl, rank),
		Description: fmt.Sprintf("%v/%v XP", humanize.Comma(xp), humanize.Comma(sc.CalculateExp(lvl+1))),

		Thumbnail: &discord.EmbedThumbnail{
			URL: avatarURL,
		},
		Author: &discord.EmbedAuthor{
			Name: name,
		},
	}

	if rank == 0 {
		e.Title = "Level " + strconv.FormatInt(lvl, 10)
	}

	{
		progress := xp - xpForPrev
		needed := xpForNext - xpForPrev

		p := float64(progress) / float64(needed)

		percent := int64(p * 100)

		e.Fields = append(e.Fields, discord.EmbedField{
			Name:  "Progress to next level",
			Value: fmt.Sprintf("%v%% (%v/%v)", percent, HumanizeInt64(progress), HumanizeInt64(needed)),
		})
	}

	reward := bot.getNextReward(sc.ID, lvl)
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

	return e
}
