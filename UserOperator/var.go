package UserOperator

import (
	"flag"
	"github.com/asticode/go-astilectron"
	CynicUClient "github.com/dxyinme/LukaComm/CynicU/Client"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)


var (
	// common use
	client *CynicUClient.Client
	CloseSign chan bool
	mu sync.Mutex
	isClosed bool
	KeeperHost = flag.String("KeeperHost", "127.0.0.1:10137", "keeper host")
	MaxMessageUpdateTime = 10 * time.Second
	MinMessageUpdateTime = 1 * time.Second

	// special for web
	conn *websocket.Conn

	// special for gui
	ChatWindow *astilectron.Window
)

func closeConnect() {
	mu.Lock()
	defer mu.Unlock()
	if !isClosed {
		isClosed = true
		close(CloseSign)
	}
}