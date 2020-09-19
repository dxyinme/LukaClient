package main

import (
	"flag"
	"fmt"
	CynicUClient "github.com/dxyinme/LukaComm/CynicU/Client"
	"github.com/dxyinme/LukaComm/chatMsg"
	"github.com/golang/glog"
	"time"
)

var (
	op = flag.String("o", "","input your operation")
	host = flag.String("host", "127.0.0.1:10137", "host")
	name = flag.String("n", "", "sender name")
	content = flag.String("c", "", "message content")
	target = flag.String("t", "", "message target")
)

func main() {
	var (
		resp *chatMsg.MsgPack
		err error
	)
	flag.Parse()
	client := &CynicUClient.Client{}
	err = client.Initial(*host, time.Second * 3)
	if err != nil {
		glog.Error(err)
	}
	switch *op {
	case "s":
		err = client.SendTo(&chatMsg.Msg{
			From: *name,
			Target: *target,
			Content: []byte(*content),
			MsgType: chatMsg.MsgType_Single,
			MsgContentType: chatMsg.MsgContentType_Text,
		})
		break
	case "p":
		resp, err = client.PullAll(&chatMsg.Ack{From: *name})
		if err != nil {
			break
		}
		for _,msg := range resp.MsgList {
			fmt.Println(msg)
		}
		break
	}
}