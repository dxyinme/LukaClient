package UserOperator

import (
	"flag"
	"github.com/asticode/go-astilectron"
	"github.com/dxyinme/LukaClient/IpcMsg"
	CynicUClient "github.com/dxyinme/LukaComm/CynicU/Client"
	"github.com/dxyinme/LukaComm/chatMsg"
	"sync"
	"time"
)


var (
	// common use
	client *CynicUClient.Client
	CloseSign chan bool
	mu sync.Mutex
	isClosed bool

	KeeperHost string

	MaxMessageUpdateTime = 10 * time.Second
	MinMessageUpdateTime = 1 * time.Second


	// special for gui

	nowChatLock sync.Mutex
	nowChatTarget string
	nowChatType	chatMsg.MsgType

	Astilectron		*astilectron.Astilectron
	MainWindow		*astilectron.Window
	LoginWindow 	*astilectron.Window
	VideoWinodw		*astilectron.Window
	LoginWg 		sync.WaitGroup
	NowLoginUser 	*IpcMsg.Login

	AssignHost		= flag.String("AssignHost", "127.0.0.1:10197", "assigneer server host")

	preMsgLoad 		= "SaveTmp/"

)

func closeConnect() {
	mu.Lock()
	defer mu.Unlock()
	if !isClosed {
		isClosed = true
		close(CloseSign)
	}
}