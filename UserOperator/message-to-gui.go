package UserOperator

import (
	"context"
	"encoding/json"
	"github.com/asticode/go-astilectron"
	"github.com/dxyinme/LukaClient/IpcMsg"
	"github.com/dxyinme/LukaComm/Assigneer"
	CynicUClient "github.com/dxyinme/LukaComm/CynicU/Client"
	"github.com/dxyinme/LukaComm/chatMsg"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

func DoSend(w *astilectron.Window, msg *IpcMsg.IpcMsg) {
	if w == nil {
		return
	}
	err := w.SendMessage(msg)
	if err != nil {
		log.Println(err)
	}
}

func Login(msg IpcMsg.IpcMsg) *IpcMsg.IpcMsg {

	// get the keeper host to location
	loginMsg := msg.Msg.(IpcMsg.Login)
	conn,err := grpc.Dial(*AssignHost, grpc.WithInsecure())
	loginClient := Assigneer.NewAssigneerClient(conn)
	resp, err := loginClient.SwitchKeeper(context.Background(), &Assigneer.SwitchKeeperReq{
		Uid: loginMsg.Name,
	})
	if err != nil {
		log.Println(err)
		return &IpcMsg.IpcMsg{
			Type: 		IpcMsg.TypeErr,
			Msg: 		err.Error(),
		}
	}
	KeeperHost = resp.Host

	log.Printf("target KeeperHost is %s", KeeperHost)

	// connect to keeper
	client = &CynicUClient.Client{}
	err = client.Initial(KeeperHost, time.Second * 3)
	if err != nil {
		return &IpcMsg.IpcMsg{
			Type:        IpcMsg.TypeErr,
			Msg:         err.Error(),
		}
	}
	go SyncMessage(msg.Msg.(IpcMsg.Login))
	nowLogin := msg.Msg.(IpcMsg.Login)
	NowLoginUser = &nowLogin
	LoginWg.Done()
	return &IpcMsg.IpcMsg{
		Type:        IpcMsg.TypeNothing,
		Msg:         "nothing",
	}
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
			if pack == nil || pack.MsgList == nil || len(pack.MsgList) == 0 {
				timeLazy *= 2
				if timeLazy > MaxMessageUpdateTime {
					timeLazy = MaxMessageUpdateTime
				}
			} else {
				timeLazy = MinMessageUpdateTime
				for _,v := range pack.MsgList {
					log.Println(v)
					DoSend(ChatWindow, &IpcMsg.IpcMsg{
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

func saveToFile(v *IpcMsg.Video) {
	fileObj, err := os.Create(v.Avid + ".webm")
	if err != nil {
		log.Println(err)
		return
	}
	defer fileObj.Close()
	//log.Println(v.Media)
	n, err := fileObj.Write(IpcMsg.ArrayBufferToByteArray(&(v.Media)))
	log.Printf("webm length: %d",n)
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
		DoSend(LoginWindow, Login(msg))
		break
	case IpcMsg.TypeMessage:
		SendMessage(msg)
		break
	case IpcMsg.TypeVideo:
		log.Println("video receive ok!!!!")
		video := msg.Msg.(IpcMsg.Video)
		saveToFile(&video)
		break
	}
	return nil
}