package main

import (
	"flag"
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/dxyinme/LukaClient/IpcMsg"
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

	// login window
	var loginWindow *astilectron.Window

	if loginWindow, err = a.NewWindow(window.LoginWindowHtml, window.LoginWindowOptions); err != nil {
		logger.Fatal(fmt.Errorf("main: new login window failed: %w", err))
	}

	UserOperator.LoginWindow = loginWindow

	loginWindow.OnMessage(UserOperator.RecvIpcMessage)

	UserOperator.LoginWg.Add(1)
	if err = loginWindow.Create(); err != nil {
		logger.Fatal(fmt.Errorf("main: creating login window failed: %w", err))
	}

	// chat window
	var chatWindow *astilectron.Window
	if chatWindow, err = a.NewWindow(window.ChatWindowHtml, window.ChatWindowOptions); err != nil {
		logger.Fatal(fmt.Errorf("main: new window failed: %w", err))
	}

	UserOperator.ChatWindow = chatWindow

	chatWindow.OnMessage(UserOperator.RecvIpcMessage)

	UserOperator.LoginWg.Wait()

	// Create windows
	if err = chatWindow.Create(); err != nil {
		logger.Fatal(fmt.Errorf("main: creating window failed: %w", err))
	}

	// login success , close loginWindow
	err = UserOperator.LoginWindow.Close()
	if err != nil {
		logger.Fatal(fmt.Errorf("main: close Login window failed: %w", err))
	}
	// login success
	UserOperator.DoSend(UserOperator.ChatWindow, &IpcMsg.IpcMsg{
		Type:        IpcMsg.TypeLoginFinished,
		Msg:         UserOperator.NowLoginUser,
	})

	// Blocking pattern
	a.Wait()
}