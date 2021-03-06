package window

import (
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

var VideoWindowOptions = &astilectron.WindowOptions{
	Center: astikit.BoolPtr(true),
	Height: astikit.IntPtr(700),
	Width: astikit.IntPtr(1200),
}

var VideoWindowHtml = ""