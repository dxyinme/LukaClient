package db

import (
	"github.com/dxyinme/LukaComm/chatMsg"
	"log"
	"testing"
)

func TestDB(t *testing.T) {
	err := NewConn("test.db")
	if err != nil {
		t.Error(err)
	}

	_, err = ExecCmd(CREATE_USERINFO_TABLE)
	if err != nil {
		t.Error(err)
	}
	_, err = ExecCmd(CREATE_MSG_TABLE)
	if err != nil {
		t.Error(err)
	}

	_ = SaveChatMsg(&chatMsg.Msg{
		From:           "tb1",
		Target:         "tb2",
		Content:        []byte("我是你爹 im fire4"),
		MsgType:        chatMsg.MsgType_Single,
		MsgContentType: chatMsg.MsgContentType_Text,
		SendTime:       "1919-08-10",
		GroupName:      "",
		Spread:         false,
		MsgId:          "tb100000000000101100",
	})

	_ = SaveChatMsg(&chatMsg.Msg{
		From:           "tb1",
		Target:         "tb2",
		Content:        []byte("我是你爹 im fire1"),
		MsgType:        chatMsg.MsgType_Single,
		MsgContentType: chatMsg.MsgContentType_Text,
		SendTime:       "1919-08-10",
		GroupName:      "",
		Spread:         false,
		MsgId:          "tb100000000000101000",
	})

	_ = SaveChatMsg(&chatMsg.Msg{
		From:           "tb1",
		Target:         "tb2",
		Content:        []byte("我是你爹 im fire2"),
		MsgType:        chatMsg.MsgType_Single,
		MsgContentType: chatMsg.MsgContentType_Text,
		SendTime:       "1919-08-10",
		GroupName:      "",
		Spread:         false,
		MsgId:          "tb100000000001010001",
	})

	_ = SaveChatMsg(&chatMsg.Msg{
		From:           "tb1",
		Target:         "tb2",
		Content:        []byte("我是你爹 im fire3"),
		MsgType:        chatMsg.MsgType_Single,
		MsgContentType: chatMsg.MsgContentType_Text,
		SendTime:       "1919-08-10",
		GroupName:      "",
		Spread:         false,
		MsgId:          "tb100000000000101010",
	})

	var ret []*chatMsg.Msg

	ret, err = LoadSingleChatMsgAll("tb1", "tb2")
	if err != nil {
		t.Error(err)
	}
	for _,v := range ret {
		log.Println(v)
	}

	err = Close()

	if err != nil {
		t.Error(err)
	}
}
