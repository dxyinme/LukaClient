package UserOperator

import (
	"encoding/json"
	"github.com/asticode/go-astilectron"
	"github.com/dxyinme/LukaClient/IpcMsg"
	CynicUClient "github.com/dxyinme/LukaComm/CynicU/Client"
	"github.com/dxyinme/LukaComm/chatMsg"
	"log"
	"time"
)

func DoSend(msg *IpcMsg.IpcMsg) {
	err := ChatWindow.SendMessage(msg)
	if err != nil {
		log.Println(err)
	}
}

func Login(msg IpcMsg.IpcMsg) *IpcMsg.IpcMsg {
	client = &CynicUClient.Client{}
	err := client.Initial(*KeeperHost, time.Second * 3)
	if err != nil {
		log.Println(err)
	}
	go  SyncMessage(msg.Msg.(IpcMsg.Login))
	return &msg
}

func SendMessage(msg IpcMsg.IpcMsg) *IpcMsg.IpcMsg {
	log.Println((msg.Msg.(chatMsg.Msg)))
	tmp := (msg.Msg.(chatMsg.Msg))
	err := client.SendTo(&tmp)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func SyncMessage(login IpcMsg.Login) {
	var (
		err error
		pack *chatMsg.MsgPack
		timeLazy time.Duration = MinMessageUpdateTime
	)
	for {
		select {
		case <-time.After(timeLazy):
			pack,err = client.Pull(&chatMsg.PullReq{
				From: login.Name,
			})
			if err != nil {
				log.Println(err)
			}
			if pack.MsgList == nil || len(pack.MsgList) == 0 {
				timeLazy *= 2
				if timeLazy > MaxMessageUpdateTime {
					timeLazy = MaxMessageUpdateTime
				}
			} else {
				timeLazy = MinMessageUpdateTime
				for _,v := range pack.MsgList {
					log.Println(v)
					DoSend(&IpcMsg.IpcMsg{
						Type:        IpcMsg.TypeMessage,
						ContextByte: nil,
						Msg:         v,
					})
				}
			}
		case <-CloseSign:
			log.Println("SyncMessage close , the connect is closed")
			return
		}
	}
}

func RecvIpcMessage(m *astilectron.EventMessage) interface{} {
	var (
		msgs 	string
		msg 	IpcMsg.IpcMsg
		err 	error
	)

	err = m.Unmarshal(&msgs)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = json.Unmarshal([]byte(msgs) ,&msg)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = msg.Marshalify()
	if err != nil {
		log.Println(err)
		return nil
	}
	switch msg.Type {
	case IpcMsg.TypeLogin:
		DoSend(Login(msg))
		break
	case IpcMsg.TypeMessage:
		SendMessage(msg)
		break
	}
	return nil
}