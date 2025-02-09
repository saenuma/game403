package internal

import (
	g143 "github.com/bankole7782/graphics143"
)

const (
	FPS      = 10
	FontSize = 30

	BlessBtn = 101
	CurseBtn = 102
)

var ObjCoords map[int]g143.RectSpecs
var CurrentScene int
var SceneLimit int
