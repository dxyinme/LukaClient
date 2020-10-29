package window

import (
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

var LoginWindowOptions = &astilectron.WindowOptions{
	Center: astikit.BoolPtr(true),
	Height: astikit.IntPtr(300),
	Width: astikit.IntPtr(500),
}

var LoginWindowHtml = "ClientExample/login-electron.html"
