package main

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"runtime"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/game403/internal"
)

const (
	fps      = 10
	fontSize = 30

	BlessBtn = 101
	CurseBtn = 102
)

var objCoords map[int]g143.RectSpecs
var currentScene int
var sceneLimit int

func main() {
	runtime.LockOSThread()

	objCoords = make(map[int]g143.RectSpecs)

	window := g143.NewWindow(1200, 800, "Game403: a game about rewards", false)
	currentScene = 1
	drawScene(window, 1)

	dirEs, _ := internal.PNGs.ReadDir("pngs")
	sceneLimit = len(dirEs)

	go func() {
		playAudio()
	}()

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}

}

func getDefaultFontPath() string {
	fontPath := filepath.Join(os.TempDir(), "g403_font.ttf")
	os.WriteFile(fontPath, internal.DefaultFont, 0777)
	return fontPath
}

func drawScene(window *glfw.Window, scene int) {
	wWidth, wHeight := window.GetSize()

	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	fontPath := getDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, fontSize)
	if err != nil {
		panic(err)
	}

	// intro text

	msg := fmt.Sprintf("Scene %d", scene)
	msgW, _ := ggCtx.MeasureString(msg)
	msgX := (wWidth - int(msgW)) / 2
	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(msg, float64(msgX), 10+fontSize)

	rawPNG, err := internal.PNGs.ReadFile(fmt.Sprintf("pngs/%d.png", scene))
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(bytes.NewReader(rawPNG))
	if err != nil {
		panic(err)
	}
	imgW, imgH := int(float64(img.Bounds().Dx())*0.9), int(float64(img.Bounds().Dy())*0.9)
	img = imaging.Fit(img, imgW, imgH, imaging.Lanczos)
	imgX := (wWidth - img.Bounds().Dx()) / 2
	ggCtx.DrawImage(img, imgX, 60)

	// buttons
	msgR1, msgR2 := "BLESS", "CURSE"
	msgR1W, _ := ggCtx.MeasureString(msgR1)
	msgR2W, _ := ggCtx.MeasureString(msgR2)

	buttonsY := wHeight - 80
	ggCtx.SetHexColor("#8B5A87")

	giftImg, _, _ := image.Decode(bytes.NewReader(internal.GiftBytes))
	giftImg = imaging.Fit(giftImg, 50, 50, imaging.Lanczos)
	ggCtx.DrawImage(giftImg, 200-60, buttonsY)
	ggCtx.Fill()

	ggCtx.DrawRoundedRectangle(200, float64(buttonsY), msgR1W+40, 50, 10)
	ggCtx.Fill()
	bbRS := g143.NRectSpecs(200-60, buttonsY, int(msgR1W)+40+60, 50)
	objCoords[BlessBtn] = bbRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(msgR1, 200+20, float64(buttonsY)+fontSize+5)

	swordImg, _, _ := image.Decode(bytes.NewReader(internal.SwordBytes))
	swordImg = imaging.Fit(swordImg, 50, 50, imaging.Lanczos)
	ggCtx.DrawImage(swordImg, 900-60, buttonsY)
	ggCtx.Fill()

	ggCtx.SetHexColor("#85836E")
	ggCtx.DrawRoundedRectangle(900, float64(buttonsY), msgR2W+40, 50, 10)
	ggCtx.Fill()
	cbRS := g143.NRectSpecs(900-60, buttonsY, int(msgR2W)+40, 50)
	objCoords[CurseBtn] = cbRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(msgR2, 900+20, float64(buttonsY)+fontSize+5)

	// send the frame to glfw window
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

}

func mouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	// wWidth, wHeight := window.GetSize()

	// var widgetRS g143.RectSpecs
	var widgetCode int

	for code, RS := range objCoords {
		if g143.InRectSpecs(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	currentScene += 1
	if currentScene > sceneLimit {
		currentScene = 1
	}

	switch widgetCode {
	case BlessBtn:
		drawScene(window, currentScene)

	case CurseBtn:
		drawScene(window, currentScene)
	}
}
