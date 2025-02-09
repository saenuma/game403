package internal

import (
	"bytes"
	"fmt"
	"image"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kovidgoyal/imaging"
)

func DrawScene(window *glfw.Window, scene int) {
	wWidth, wHeight := window.GetSize()

	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	fontPath := GetDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, FontSize)
	if err != nil {
		panic(err)
	}

	// intro text

	msg := fmt.Sprintf("Scene %d", scene)
	msgW, _ := ggCtx.MeasureString(msg)
	msgX := (wWidth - int(msgW)) / 2
	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(msg, float64(msgX), 10+FontSize)

	rawPNG, err := PNGs.ReadFile(fmt.Sprintf("pngs/%d.png", scene))
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(bytes.NewReader(rawPNG))
	if err != nil {
		panic(err)
	}
	imgW, imgH := int(1000*0.9), int(700*0.9)
	if img.Bounds().Dx() > imgW {
		img = imaging.Fit(img, imgW, imgH, imaging.Lanczos)
	} else {
		img = imaging.Fill(img, imgW, imgH, imaging.Center, imaging.Lanczos)
	}
	imgX := (wWidth - img.Bounds().Dx()) / 2
	ggCtx.DrawImage(img, imgX, 60)

	// buttons
	msgR1, msgR2 := "BLESS", "CURSE"
	msgR1W, _ := ggCtx.MeasureString(msgR1)
	msgR2W, _ := ggCtx.MeasureString(msgR2)

	buttonsY := wHeight - 80
	ggCtx.SetHexColor("#8B5A87")

	giftImg, _, _ := image.Decode(bytes.NewReader(GiftBytes))
	giftImg = imaging.Fit(giftImg, 50, 50, imaging.Lanczos)
	ggCtx.DrawImage(giftImg, 200-60, buttonsY)
	ggCtx.Fill()

	ggCtx.DrawRoundedRectangle(200, float64(buttonsY), msgR1W+40, 50, 10)
	ggCtx.Fill()
	bbRS := g143.NewRect(200-60, buttonsY, int(msgR1W)+40+60, 50)
	ObjCoords[BlessBtn] = bbRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(msgR1, 200+20, float64(buttonsY)+FontSize+5)

	swordImg, _, _ := image.Decode(bytes.NewReader(SwordBytes))
	swordImg = imaging.Fit(swordImg, 50, 50, imaging.Lanczos)
	ggCtx.DrawImage(swordImg, 900-60, buttonsY)
	ggCtx.Fill()

	ggCtx.SetHexColor("#85836E")
	ggCtx.DrawRoundedRectangle(900, float64(buttonsY), msgR2W+40, 50, 10)
	ggCtx.Fill()
	cbRS := g143.NewRect(900-60, buttonsY, int(msgR2W)+40, 50)
	ObjCoords[CurseBtn] = cbRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(msgR2, 900+20, float64(buttonsY)+FontSize+5)

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

}

func MouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	// wWidth, wHeight := window.GetSize()

	// var widgetRS g143.Rect
	var widgetCode int

	for code, RS := range ObjCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	CurrentScene += 1
	if CurrentScene > SceneLimit {
		CurrentScene = 1
	}

	switch widgetCode {
	case BlessBtn:
		DrawScene(window, CurrentScene)

	case CurseBtn:
		DrawScene(window, CurrentScene)
	}
}
