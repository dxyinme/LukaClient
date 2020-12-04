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
	var err error
	// Set logger
	logger := log.New(log.Writer(), log.Prefix(), log.Flags())
	// Create astilectron
	UserOperator.Astilectron, err = astilectron.New(logger, astilectron.Options{
		AppName:           "Luka",
		BaseDirectoryPath: "example",
	})
	if err != nil {
		logger.Fatal(fmt.Errorf("main: creating astilectron failed: %w", err))
	}
	UserOperator.Astilectron.On(astilectron.EventNameAppCmdQuit, func(e astilectron.Event) (deleteListener bool) {
		logger.Fatal(fmt.Errorf("main: app closed"))
		return
	})
	defer UserOperator.Astilectron.Close()

	// Handle signals
	UserOperator.Astilectron.HandleSignals()

	// Start
	if err = UserOperator.Astilectron.Start(); err != nil {
		logger.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
	}

	// login window
	var loginWindow *astilectron.Window

	if loginWindow, err = UserOperator.Astilectron.NewWindow(
		window.LoginWindowHtml, window.LoginWindowOptions);
	err != nil {
		logger.Fatal(fmt.Errorf("main: new login window failed: %w", err))
	}

	UserOperator.LoginWindow = loginWindow

	loginWindow.OnMessage(UserOperator.RecvIpcMessage)

	UserOperator.LoginWg.Add(1)
	if err = loginWindow.Create(); err != nil {
		logger.Fatal(fmt.Errorf("main: creating login window failed: %w", err))
	}

	// main window
	var mainWindow *astilectron.Window
	if mainWindow, err = UserOperator.Astilectron.NewWindow(
		window.MainWindowHtml, window.MainWindowOptions);
	err != nil {
		logger.Fatal(fmt.Errorf("main: new window failed: %w", err))
	}

	UserOperator.MainWindow = mainWindow

	mainWindow.OnMessage(UserOperator.RecvIpcMessage)

	UserOperator.LoginWg.Wait()

	// Create windows
	if err = mainWindow.Create(); err != nil {
		logger.Fatal(fmt.Errorf("main: creating window failed: %w", err))
	}

	// login success , close loginWindow
	err = UserOperator.LoginWindow.Close()
	if err != nil {
		logger.Fatal(fmt.Errorf("main: close Login window failed: %w", err))
	}
	// login success
	UserOperator.DoSend(UserOperator.MainWindow, &IpcMsg.IpcMsg{
		Type:        IpcMsg.TypeLoginFinished,
		Msg:         UserOperator.NowLoginUser,
	})

	err = UserOperator.MainWindow.OpenDevTools()
	if err != nil {
		logger.Fatal(fmt.Errorf("main: main window open devTools failed: %w", err))
	}
	// Blocking pattern
	UserOperator.Astilectron.Wait()
}