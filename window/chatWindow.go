package window

import (
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

var ChatWindowOptions = &astilectron.WindowOptions{
	Center: astikit.BoolPtr(true),
	Height: astikit.IntPtr(700),
	Width:  astikit.IntPtr(700),
}


var ChatWindowHtml = "ClientExample/chat-electron.html"