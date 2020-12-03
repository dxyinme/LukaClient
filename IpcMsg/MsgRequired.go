package IpcMsg

import "github.com/dxyinme/LukaComm/chatMsg"

type MsgRequired struct {
	From string
	MsgType chatMsg.MsgType
}
