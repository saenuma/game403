package internal

import (
	"embed"
	"os"
	"path/filepath"
)

//go:embed Roboto-Light.ttf
var DefaultFont []byte

//go:embed pngs
var PNGs embed.FS

//go:embed audio.mp3
var AudioBytes []byte

//go:embed gift.png
var GiftBytes []byte

//go:embed sword.png
var SwordBytes []byte

func GetDefaultFontPath() string {
	fontPath := filepath.Join(os.TempDir(), "g403_font.ttf")
	os.WriteFile(fontPath, DefaultFont, 0777)
	return fontPath
}
