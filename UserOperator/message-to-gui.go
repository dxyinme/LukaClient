package UserOperator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/dxyinme/LukaClient/IpcMsg"
	"github.com/dxyinme/LukaClient/db"
	"github.com/dxyinme/LukaComm/Assigneer"
	CynicUClient "github.com/dxyinme/LukaComm/CynicU/Client"
	"github.com/dxyinme/LukaComm/chatMsg"
	"github.com/dxyinme/LukaComm/util"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)


var (
	lastTime int64
	lastTips int64
	msgMutex sync.Mutex
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
	err = db.NewConn(preMsgLoad + loginMsg.Name + ".db")
	if err != nil {
		log.Fatal(fmt.Errorf("login: connect to db error : %w", err))
	}

	_, err = db.ExecCmd(db.CREATE_MSG_TABLE)
	if err != nil {
		log.Fatal(fmt.Errorf("login: db prepare error : %w", err))
	}

	_, err = db.ExecCmd(db.CREATE_USERINFO_TABLE)
	if err != nil {
		log.Fatal(fmt.Errorf("login: db prepare error : %w", err))
	}

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
	log.Println(msg.Msg.(chatMsg.Msg))
	tmp := msg.Msg.(chatMsg.Msg)
	msgMutex.Lock()
	lastTime, lastTips, tmp.MsgId = util.MsgIdGen(NowLoginUser.Name,0, lastTime, lastTips)
	msgMutex.Unlock()

	err := client.SendTo(&tmp)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = db.SaveChatMsg(&tmp)
	if err != nil {
		log.Println(err)
		return nil
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
					DoSend(MainWindow, &IpcMsg.IpcMsg{
						Type:        IpcMsg.TypeMessage,
						ContextByte: nil,
						Msg:         v,
					})
					err = db.SaveChatMsg(v)
					if err != nil {
						log.Println(err)
					}
				}
			}
		case <-CloseSign:
			log.Println("SyncMessage close , the connect is closed")
			return
		}
	}
}

func SetRecvTarget(msg IpcMsg.IpcMsg) {
	msgRequired := msg.Msg.(IpcMsg.MsgRequired)
	var (
		retMsg []*chatMsg.Msg
		err error
	)
	nowChatLock.Lock()
	defer nowChatLock.Unlock()

	nowChatTarget = msgRequired.From
	nowChatType = msgRequired.MsgType

	if msgRequired.MsgType == chatMsg.MsgType_Single {
		retMsg,err = db.LoadSingleChatMsgAll(msgRequired.From, NowLoginUser.Name)
		if err != nil {
			log.Println(err)
		}
	} else {
		retMsg,err = db.LoadGroupChatMsgAll(msgRequired.From)
		if err != nil {
			log.Println(err)
		}
	}
	for i := 0; i <= len(retMsg); i ++ {
		DoSend(MainWindow, &IpcMsg.IpcMsg{
			Type:        IpcMsg.TypeMessage,
			ContextByte: nil,
			Msg:         retMsg[i],
		})
	}
	log.Println("finish change target")
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
	//log.Println(msg)
	switch msg.Type {
	case IpcMsg.TypeLogin:
		DoSend(LoginWindow, Login(msg))
		break
	case IpcMsg.TypeMessage:
		SendMessage(msg)
		break
	case IpcMsg.TypeMessageRequired:
		SetRecvTarget(msg)
		break
	}
	return nil
}