package IpcMsg

import (
	"github.com/dxyinme/LukaComm/chatMsg"
	"github.com/dxyinme/LukaComm/util"
)

const (
	TypeNothing			= -1

	TypeErr 			= 0

	TypeLogin 			= 1

	TypeMessage 		= 2

	TypeLoginFinished 	= 3

	//TypeVideo			= 4
	// chatWindow is on, message is required
	TypeMessageRequired	= 6

	TypeNewWindow 		= 7

	TypeGroupOperator 	= 8

	TypeWindowGroupWindow 	= "group"
)

type IpcMsg struct {
	Type 		int		`json:"Type"`
	ContextByte []byte	`json:"ContextByte"`


	Msg 		interface{}
}

func (m *IpcMsg) Marshal(Type int, v interface{}) {
	contextByte, err := util.IJson.Marshal(v)
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
		err = util.IJson.Unmarshal(m.ContextByte, &tmp)
		m.Msg = tmp
		if err != nil {
			return err
		}
		break
	case TypeLogin:
		var tmp Login
		err = util.IJson.Unmarshal(m.ContextByte, &tmp)
		m.Msg = tmp
		if err != nil {
			return err
		}
		break
	case TypeMessageRequired:
		var tmp MsgRequired
		err = util.IJson.Unmarshal(m.ContextByte, &tmp)
		m.Msg = tmp
		if err != nil {
			return err
		}
		break
	case TypeNewWindow:
		var tmp NewWindow
		err = util.IJson.Unmarshal(m.ContextByte, &tmp)
		m.Msg = tmp
		if err != nil {
			return err
		}
		break
	case TypeGroupOperator:
		var tmp GroupOperator
		err = util.IJson.Unmarshal(m.ContextByte, &tmp)
		m.Msg = tmp
		if err != nil {
			return err
		}
		break
	}
	return nil
}