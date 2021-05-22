package UserOperator

import (
	"flag"
	"github.com/asticode/go-astilectron"
	"github.com/dxyinme/LukaClient/IpcMsg"
	CynicUClient "github.com/dxyinme/LukaComm/CynicU/Client"
	"github.com/dxyinme/LukaComm/CynicU/SendMsg"
	"github.com/dxyinme/LukaComm/chatMsg"
	"sync"
	"time"
)


var (
	// common use
	client *CynicUClient.Client
	udpClient *SendMsg.Client
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
	MainWindow		*astilectron.Window = nil
	LoginWindow 	*astilectron.Window = nil
	VideoWindow		*astilectron.Window = nil

	GroupWindow		*astilectron.Window = nil
	GroupWindowAlive bool = false
	GroupMutex		sync.Mutex

	LoginWg 		sync.WaitGroup
	NowLoginUser 	*IpcMsg.Login

	AssignHost		= flag.String("AssignHost", "127.0.0.1:10197", "assigneer server host")

	AuthHost 		= flag.String("AuthHost", "127.0.0.1:20020", "auth server host")

	FileServerHost 	= flag.String("FileServer", "127.0.0.1:18081", "file server host")

	enableUdp		= flag.Bool("enableUdp", false, "enable udp sender")

	// DB save.
	preMsgLoad 		= "SaveTmp/"

	priKeyLoad 		= "privateKey/"

)

func closeConnect() {
	mu.Lock()
	defer mu.Unlock()
	if !isClosed {
		isClosed = true
		close(CloseSign)
	}
}