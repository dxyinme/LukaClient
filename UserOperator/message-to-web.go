package UserOperator

import (
	"encoding/json"
	CynicUClient "github.com/dxyinme/LukaComm/CynicU/Client"
	"github.com/dxyinme/LukaComm/chatMsg"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)


func readLoop(uid string) {
	TimeLazy := MinMessageUpdateTime
	for {
		glog.Info(TimeLazy)
		select {
		case <-time.After(TimeLazy):
			{
				pack, err := client.Pull(&chatMsg.PullReq{
					From: uid,
				})
				if err != nil {
					glog.Errorln(err)
				}
				// flow control
				if err != nil || len(pack.MsgList) == 0 {
					TimeLazy = TimeLazy * 2
					if TimeLazy > MaxMessageUpdateTime {
						TimeLazy = MaxMessageUpdateTime
					}
				} else {
					TimeLazy = MinMessageUpdateTime
				}

				if pack.MsgList != nil {
					for i := 0; i < len(pack.MsgList); i++ {
						err = conn.WriteJSON(pack.MsgList[i])
						if err != nil {
							glog.Infof("%v read error", pack.MsgList[i])
						}
					}
				}
			}
		case <-CloseSign:
			glog.Infof("read fail: connect is close")
			goto ERROR
		}
	}
ERROR:
	closeConnect()
}

func sendLoop(uid string) {
	var (
		msg = &chatMsg.Msg{}
		data []byte
		err error
	)
	for {
		if _,data,err = conn.ReadMessage(); err != nil {
			glog.Error(err)
			goto ERROR
		}
		err = json.Unmarshal(data, msg)
		glog.Infof("receive : %v", msg)
		if err != nil {
			glog.Error(err)
			goto ERROR
		}
		msg.From = uid
		msg.MsgContentType = chatMsg.MsgContentType_Text
		msg.MsgType = chatMsg.MsgType_Single
		err = client.SendTo(msg)
		if err != nil {
			glog.Error(err)
		}
	}
ERROR:
	closeConnect()
}

func serve() error {
	select {
	case <-CloseSign:
		break
	}
	return nil
}

func closeConnect() {
	mu.Lock()
	defer mu.Unlock()
	if !isClosed {
		isClosed = true
		close(CloseSign)
	}
}

// 登录处理，我们将会把他升级成websocket
// 一个机器只允许有一个同时登录用户
func Connect(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	uid := r.Form.Get("uid")
	if err != nil || uid == "" {
		_,_ = w.Write([]byte("CanNot make connect"))
		return
	}
	glog.Infof("User %s , login Success" , uid)
	client = &CynicUClient.Client{}
	err = client.Initial(*KeeperHost, time.Second * 3)
	defer client.Close()
	if err != nil {
		glog.Error(err)
	}
	upgrade := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// 支持跨域
			return true
		},
	}
	if conn, err = upgrade.Upgrade(w,r,nil); err != nil {
		glog.Error(err)
		return
	}
	isClosed = false
	CloseSign = make(chan bool, 1)
	go readLoop(uid)
	go sendLoop(uid)
	if err = serve(); err != nil {
		glog.Errorf("User %s Disconnected , because of %v", uid, err)
	}
}
