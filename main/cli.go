package main

import (
	"flag"
	"fmt"
	"github.com/dxyinme/LukaComm/CynicU/SendMsg"
	"github.com/dxyinme/LukaComm/util/Const"
	"log"
	"time"

	CynicUClient "github.com/dxyinme/LukaComm/CynicU/Client"
	"github.com/dxyinme/LukaComm/chatMsg"
	"github.com/golang/glog"
)

var (
	op      = flag.String("o", "", "input your operation")
	host    = flag.String("host", "127.0.0.1:10199", "host")
	name    = flag.String("n", "", "sender name")
	content = flag.String("c", "", "message content")
	target  = flag.String("t", "", "message target")
	group   = flag.String("group", "", "group name")
)

func main() {
	var (
		resp *chatMsg.MsgPack
		err  error
	)
	flag.Parse()
	client := &CynicUClient.Client{}
	err = client.Initial(*host, time.Second*3)
	if err != nil {
		glog.Error(err)
	}
	switch *op {
	case "lg":
		err = client.GroupOp(Const.LeaveGroup, &chatMsg.GroupReq{
			Uid:       *name,
			GroupName: *group,
			IsCopy:    false,
		})
		if err != nil {
			log.Println(err)
		}
		break
	case "jg":
		err = client.GroupOp(Const.JoinGroup, &chatMsg.GroupReq{
			Uid:       *name,
			GroupName: *group,
			IsCopy:    false,
		})
		if err != nil {
			log.Println(err)
		}
		break
	case "cg":
		err = client.GroupOp(Const.CreateGroup, &chatMsg.GroupReq{
			Uid:       *name,
			GroupName: *group,
			IsCopy:    false,
		})
		if err != nil {
			log.Println(err)
		}
		break
	case "dg":
		err = client.GroupOp(Const.DeleteGroup, &chatMsg.GroupReq{
			Uid:       *name,
			GroupName: *group,
			IsCopy:    false,
		})
		if err != nil {
			log.Println(err)
		}
		break
	case "sg":
		err = client.SendTo(&chatMsg.Msg{
			From:           *name,
			Target:         "",
			Content:        []byte(*content),
			MsgType:        chatMsg.MsgType_Group,
			MsgContentType: chatMsg.MsgContentType_Text,
			GroupName:      *group,
			Spread:         true,
		})
		if err != nil {
			log.Println(err)
		}
		break
	case "s":
		err = client.SendTo(&chatMsg.Msg{
			From:           *name,
			Target:         *target,
			Content:        []byte(*content),
			MsgType:        chatMsg.MsgType_Single,
			MsgContentType: chatMsg.MsgContentType_Text,
		})
		break
	case "su":
		clientUDP := SendMsg.NewClient(*host)
		err = clientUDP.SendTo(&chatMsg.Msg{
			From:           *name,
			Target:         *target,
			Content:        []byte(*content),
			MsgType:        chatMsg.MsgType_Single,
			MsgContentType: chatMsg.MsgContentType_Text,
		})
		if err != nil {
			log.Println(err)
		}
		break
	case "p":
		resp, err = client.Pull(&chatMsg.PullReq{From: *name})
		if err != nil {
			break
		}
		for _, msg := range resp.MsgList {
			fmt.Println(msg)
		}
		break
	case "pa":
		resp, err = client.PullAll(&chatMsg.PullReq{From: *name})
		if err != nil {
			break
		}
		for _, msg := range resp.MsgList {
			fmt.Println(msg)
		}
		break
	}
}
