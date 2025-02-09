package main

import (
	"runtime"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/game403/internal"
)

func main() {
	runtime.LockOSThread()

	internal.ObjCoords = make(map[int]g143.Rect)

	window := g143.NewWindow(1200, 800, "Game403: a game about rewards", false)
	internal.CurrentScene = 1
	internal.DrawScene(window, 1)

	dirEs, _ := internal.PNGs.ReadDir("pngs")
	internal.SceneLimit = len(dirEs)

	go func() {
		playAudio()
	}()

	// respond to the mouse
	window.SetMouseButtonCallback(internal.MouseBtnCallback)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		time.Sleep(time.Second/time.Duration(internal.FPS) - time.Since(t))
	}

}
