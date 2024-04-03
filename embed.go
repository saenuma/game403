package main

import (
	"embed"
	_ "embed"
)

//go:embed Roboto-Light.ttf
var DefaultFont []byte

//go:embed pngs
var PNGs embed.FS
