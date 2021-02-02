package UserOperator

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/asticode/go-astilectron"
	"github.com/dxyinme/LukaClient/IpcMsg"
	"github.com/dxyinme/LukaClient/db"
	"github.com/dxyinme/LukaComm/Assigneer"
	CynicUClient "github.com/dxyinme/LukaComm/CynicU/Client"
	"github.com/dxyinme/LukaComm/CynicU/SendMsg"
	"github.com/dxyinme/LukaComm/chatMsg"
	"github.com/dxyinme/LukaComm/util"
	utilCrypto "github.com/dxyinme/LukaComm/util/crypto"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"log"
	"os"
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

func FileExist(path string) (bool,error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func TryGenRSAKey(uid string) error {
	priKey, pubKey, err := utilCrypto.NewRsaKey()
	if err != nil {
		return err
	}
	block := &pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: priKey,
	}
	isExist, err := FileExist(priKeyLoad + uid + ".pem")
	if err != nil {
		return err
	}
	if !isExist {
		file, err := os.Create(priKeyLoad + uid + ".pem")
		defer file.Close()
		if err != nil {
			return err
		}
		err = pem.Encode(file, block)
		if err != nil {
			return err
		}
		err = client.SetAuthPubKey(*AuthHost, uid, pubKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func Login(msg IpcMsg.IpcMsg) *IpcMsg.IpcMsg {
	// get the keeper host to location
	loginMsg := msg.Msg.(IpcMsg.Login)
	conn,err := grpc.Dial(*AssignHost, grpc.WithInsecure())
	loginClient := Assigneer.NewAssigneerClient(conn)
	resp, err := loginClient.SwitchKeeper(context.Background(), &Assigneer.SwitchKeeperReq{
		Uid: loginMsg.Name,
	})
	// loginMsg.Name is uid !!
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

	// udp client
	udpClient = SendMsg.NewClient(KeeperHost)
	// grpc client
	client = &CynicUClient.Client{}
	err = client.Initial(KeeperHost, time.Second * 3)
	if err != nil {
		return &IpcMsg.IpcMsg{
			Type:        IpcMsg.TypeErr,
			Msg:         err.Error(),
		}
	}
	err = TryGenRSAKey(loginMsg.Name)
	if err != nil {
		log.Println(err)
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
	tmp.SendTime = time.Now().String()
	tmpBytes, err := proto.Marshal(&tmp)
	if tmp.SecretLevel == 1 {
		encodeAESPlainText(&tmp)
		goto SEND_GRPC
	}
	if err != nil {
		log.Println(err)
		return nil
	}
	if len(tmpBytes) <= SendMsg.PacketSize {
		log.Println("send in udp")
		err = udpClient.SendTo(&tmp)
		if err != nil {
			log.Println(err)
			return nil
		}
		goto SAVE_DB
	}
	log.Println("send in grpc")
SEND_GRPC:
	err = client.SendTo(&tmp)
	if err != nil {
		log.Println(err)
		return nil
	}
SAVE_DB:
	err = db.SaveChatMsg(&tmp)
	if err != nil {
		log.Println(err)
		return nil
	}
	return nil
}

func solveKeyAgreement(msg *chatMsg.Msg) {
	cryptor := GetNowUserPrivateKey()
	AESKey , err := cryptor.Decode(msg.Content)
	if err != nil {
		log.Println(err)
	}
	err = db.SaveUserInfo(&IpcMsg.UserInfo{
		Uid:    msg.From,
		Name:   msg.From,
		AESKey: AESKey,
	})
	if err != nil {
		log.Println(err)
	}
}

func sendKeyAgreement(msg IpcMsg.Secret) {
	pubKey, err := client.GetAuthPubKey(*AuthHost, msg.Target)
	if err != nil {
		log.Println(err)
		return
	}
	keyAfterEncode, err := utilCrypto.EncodePub(msg.AESKey, pubKey)
	if err != nil {
		log.Println(err)
		return
	}
	keyAgreement := &chatMsg.Msg{
		From:           msg.From,
		Target:         msg.Target,
		Content:        keyAfterEncode,
		MsgType:        chatMsg.MsgType_Single,
		MsgContentType: chatMsg.MsgContentType_KeyAgreement,
		SendTime:       time.Now().String(),
		MsgId:          "",
	}
	err = client.SendTo(keyAgreement)
	if err != nil {
		log.Println(err)
	}
	err = db.SaveUserInfo(&IpcMsg.UserInfo{
		Uid:    msg.Target,
		Name:   msg.Target,
		AESKey: msg.AESKey,
	})
	if err != nil {
		log.Println(err)
	}
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
					if v.MsgContentType == chatMsg.MsgContentType_KeyAgreement {
						solveKeyAgreement(v)
						continue
					}
					if v.SecretLevel == 1 {
						log.Println("is AES decoding: ", v)
						decodeAESCipherText(v)
					}
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
	for i := 0; i < len(retMsg); i ++ {
		DoSend(MainWindow, &IpcMsg.IpcMsg{
			Type:        IpcMsg.TypeMessage,
			ContextByte: nil,
			Msg:         retMsg[i],
		})
	}
	log.Println("finish change target")
}


func DoGroup(op IpcMsg.GroupOperator) {
	var (
		err error
		req = &chatMsg.GroupReq{
			Uid:       op.Uid,
			GroupName: op.GroupName,
			IsCopy:    false,
		}
	)
	log.Printf("group operator %s", op.GroupOp)
	err = client.GroupOp(op.GroupOp, req)
	if err != nil {
		log.Println(err)
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
	case IpcMsg.TypeNewWindow:
		CreateWindow(msg.Msg.(IpcMsg.NewWindow).WindowType)
		break
	case IpcMsg.TypeGroupOperator:
		DoGroup(msg.Msg.(IpcMsg.GroupOperator))
		break
	case IpcMsg.TypeSecret:
		sendKeyAgreement(msg.Msg.(IpcMsg.Secret))
		break
	}

	return nil
}