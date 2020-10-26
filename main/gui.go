package main

import (
	"flag"
	"fmt"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/dxyinme/LukaClient/IpcMsg"
	"github.com/dxyinme/LukaClient/UserOperator"
	"log"
)

func main() {
	flag.Parse()
	// Set logger
	logger := log.New(log.Writer(), log.Prefix(), log.Flags())
	// Create astilectron
	a, err := astilectron.New(logger, astilectron.Options{
		AppName:           "Luka",
		BaseDirectoryPath: "example",
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

	//// login window
	var loginWindow *astilectron.Window

	if loginWindow, err = a.NewWindow("ClientExample/login-electron.html", &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(300),
		Width: astikit.IntPtr(500),
	}); err != nil {
		logger.Fatal(fmt.Errorf("main: new login window failed: %w", err))
	}

	UserOperator.LoginWindow = loginWindow

	loginWindow.OnMessage(UserOperator.RecvIpcMessage)

	UserOperator.LoginWg.Add(1)
	if err = loginWindow.Create(); err != nil {
		logger.Fatal(fmt.Errorf("main: creating login window failed: %w", err))
	}


	// chat window
	var w *astilectron.Window
	if w, err = a.NewWindow("ClientExample/chat-electron.html", &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(700),
		Width:  astikit.IntPtr(700),
	}); err != nil {
		logger.Fatal(fmt.Errorf("main: new window failed: %w", err))
	}

	UserOperator.ChatWindow = w

	w.OnMessage(UserOperator.RecvIpcMessage)

	UserOperator.LoginWg.Wait()

	// Create windows
	if err = w.Create(); err != nil {
		logger.Fatal(fmt.Errorf("main: creating window failed: %w", err))
	}

	// login success
	UserOperator.DoSend(UserOperator.ChatWindow, &IpcMsg.IpcMsg{
		Type:        IpcMsg.TypeLoginFinished,
		Msg:         UserOperator.NowLoginUser,
	})

	// Blocking pattern
	a.Wait()
}