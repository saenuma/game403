package main

import (
	"embed"
)

//go:embed Roboto-Light.ttf
var DefaultFont []byte

//go:embed pngs
var PNGs embed.FS

//go:embed audio.mp3
var AudioBytes []byte
