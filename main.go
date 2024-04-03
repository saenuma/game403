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
	"github.com/ebitengine/oto/v3"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/hajimehoshi/go-mp3"
)

const (
	fps        = 10
	sceneLimit = 5
	fontSize   = 30

	BlessBtn = 101
	CurseBtn = 102
)

var objCoords map[int]g143.RectSpecs
var currentScene int

func main() {
	runtime.LockOSThread()

	objCoords = make(map[int]g143.RectSpecs)

	window := g143.NewWindow(1200, 800, "Game403: a game about rewards", false)
	currentScene = 1
	drawScene(window, 1)

	go func() {
		// Convert the pure bytes into a reader object that can be used with the mp3 decoder
		fileBytesReader := bytes.NewReader(AudioBytes)

		// Decode file
		decodedMp3, err := mp3.NewDecoder(fileBytesReader)
		if err != nil {
			panic("mp3.NewDecoder failed: " + err.Error())
		}

		// Prepare an Oto context (this will use your default audio device) that will
		// play all our sounds. Its configuration can't be changed later.

		op := &oto.NewContextOptions{}

		// Usually 44100 or 48000. Other values might cause distortions in Oto
		op.SampleRate = 44100

		// Number of channels (aka locations) to play sounds from. Either 1 or 2.
		// 1 is mono sound, and 2 is stereo (most speakers are stereo).
		op.ChannelCount = 2

		// Format of the source. go-mp3's format is signed 16bit integers.
		op.Format = oto.FormatSignedInt16LE

		// Remember that you should **not** create more than one context
		otoCtx, readyChan, err := oto.NewContext(op)
		if err != nil {
			panic("oto.NewContext failed: " + err.Error())
		}
		// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
		<-readyChan

		// Create a new 'player' that will handle our sound. Paused by default.
		player := otoCtx.NewPlayer(decodedMp3)

		// Play starts playing the sound and returns without waiting for it (Play() is async).
		player.Play()

		// We can wait for the sound to finish playing using something like this
		for player.IsPlaying() {
			time.Sleep(time.Millisecond)
		}

		// Now that the sound finished playing, we can restart from the beginning (or go to any location in the sound) using seek
		// newPos, err := player.(io.Seeker).Seek(0, io.SeekStart)
		// if err != nil{
		//     panic("player.Seek failed: " + err.Error())
		// }
		// println("Player is now at position:", newPos)
		// player.Play()
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
	os.WriteFile(fontPath, DefaultFont, 0777)
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

	rawPNG, err := PNGs.ReadFile(fmt.Sprintf("pngs/%d.png", scene))
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
	ggCtx.DrawRoundedRectangle(200, float64(buttonsY), msgR1W+40, 50, 10)
	ggCtx.Fill()
	bbRS := g143.NRectSpecs(200, buttonsY, int(msgR1W+40), 50)
	objCoords[BlessBtn] = bbRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(msgR1, 200+20, float64(buttonsY)+fontSize+5)

	ggCtx.SetHexColor("#85836E")
	ggCtx.DrawRoundedRectangle(900, float64(buttonsY), msgR2W+40, 50, 10)
	ggCtx.Fill()
	cbRS := g143.NRectSpecs(900, buttonsY, int(msgR2W)+40, 50)
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
