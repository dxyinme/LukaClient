package main

import (
	"flag"
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/dxyinme/LukaClient/UserOperator"
	"github.com/dxyinme/LukaClient/window"
	"log"
)

func main() {
	flag.Parse()
	// Set logger
	logger := log.New(log.Writer(), log.Prefix(), log.Flags())
	// Create astilectron
	a, err := astilectron.New(logger, astilectron.Options{
		AppName:           "Luka-video",
		BaseDirectoryPath: "example-video",
	})
	if err != nil {
		logger.Fatal(fmt.Errorf("main: creating astilectron failed: %w", err))
	}
	defer a.Close()

	// Handle signals
	a.HandleSignals()

	// Start
	if err = a.Start(); err != nil {
		logger.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
	}

	var videoWindow *astilectron.Window

	if videoWindow, err = a.NewWindow(window.VideoWindowHtml, window.VideoWindowOptions); err != nil {
		logger.Fatal(fmt.Errorf("main: new video window failed: %w", err))
	}

	UserOperator.VideoWinodw = videoWindow

	videoWindow.OnMessage(UserOperator.RecvIpcMessage)

	// Create windows
	if err = videoWindow.Create(); err != nil {
		logger.Fatal(fmt.Errorf("main: creating window failed: %w", err))
	}

	a.Wait()
}