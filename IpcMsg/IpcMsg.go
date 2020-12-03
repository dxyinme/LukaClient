package IpcMsg

import (
	"encoding/json"
	"github.com/dxyinme/LukaComm/chatMsg"
)

const (
	TypeNothing			= -1

	TypeErr 			= 0

	TypeLogin 			= 1

	TypeMessage 		= 2

	TypeLoginFinished 	= 3

	//TypeVideo			= 4
	// chatWindow is on, message is required
	TypeMessageRequired	= 5
)

type IpcMsg struct {
	Type 		int		`json:"Type"`
	ContextByte []byte	`json:"ContextByte"`


	Msg 		interface{}
}

func (m *IpcMsg) Marshal(Type int, v interface{}) {
	contextByte, err := json.Marshal(v)
	if err != nil {
		return
	}
	m.Type = Type
	m.ContextByte = contextByte
	m.Msg = v
}

func (m *IpcMsg) Unmarshal(vPtr interface{}) error {

	return nil
}

func (m *IpcMsg) Marshalify() error {
	var err error
	switch m.Type {
	case TypeMessage:
		var tmp chatMsg.Msg
		err = json.Unmarshal(m.ContextByte, &tmp)
		m.Msg = tmp
		if err != nil {
			return err
		}
		break
	case TypeLogin:
		var tmp Login
		err = json.Unmarshal(m.ContextByte, &tmp)
		m.Msg = tmp
		if err != nil {
			return err
		}
		break
	case TypeMessageRequired:
		var tmp MsgRequired
		err = json.Unmarshal(m.ContextByte, &tmp)
		if err != nil {
			return err
		}
		break
	}
	return nil
}