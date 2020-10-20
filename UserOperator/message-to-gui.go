package UserOperator

import (
	"encoding/json"
	"github.com/asticode/go-astilectron"
	"github.com/dxyinme/LukaClient/IpcMsg"
	"github.com/dxyinme/LukaComm/chatMsg"
	"log"
)

func DoSend(msg *IpcMsg.IpcMsg) {
	err := ChatWindow.SendMessage(msg)
	if err != nil {
		log.Println(err)
	}
}

func Login(msg IpcMsg.IpcMsg) *IpcMsg.IpcMsg {
	return &msg
}

func RecvMessage(msg IpcMsg.IpcMsg) *IpcMsg.IpcMsg {
	log.Println(string(msg.Msg.(chatMsg.Msg).Content))
	return nil
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
		RecvMessage(msg)
		break
	}
	return nil
}