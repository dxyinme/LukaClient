package UserOperator

import (
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/dxyinme/LukaClient/IpcMsg"
	"github.com/dxyinme/LukaClient/window"
	"log"
)

func CreateWindow(windowType string) {
	log.Printf("create window : %s", windowType)
	if windowType == IpcMsg.TypeWindowGroupWindow {
		GroupWindowCreate()
	} else {
		log.Println("No such window to create")
	}
}

func GroupWindowCreate() {
	GroupMutex.Lock()
	defer GroupMutex.Unlock()
	if !GroupWindowAlive {
		GroupWindowAlive = true
		var (
			err error
		)
		if GroupWindow, err = Astilectron.NewWindow(
			window.GroupWindowHtml, window.GroupWindowOptions);
			err != nil {
			log.Fatal(fmt.Errorf("UserOperator.GroupWindowCreate: new window failed: %w", err))
		}
		GroupWindow.OnMessage(RecvIpcMessage)
		GroupWindow.On(astilectron.EventNameWindowEventClosed, func(e astilectron.Event) (deleteListener bool) {
			log.Println("Close group window")
			GroupWindowAlive = false
			deleteListener = false
			return
		})
		if err = GroupWindow.Create(); err != nil {
			log.Fatal(fmt.Errorf("UserOperator.GroupWindowCreate: creating window failed: %w", err))
		}
	} else {
		return
	}
}