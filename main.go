package main

import (
	"fmt"
	CynicUClient "github.com/dxyinme/LukaComm/CynicU/Client"
	"github.com/dxyinme/LukaComm/chatMsg"
	"github.com/golang/glog"
	"time"
)

var (
	name string
)

func main() {
	var (
		err error
	)
	client := &CynicUClient.Client{}
	fmt.Println(
		"op '--s %Target %Content' for SendTo\n " +
			"op '--p' for pull\n " +
			"op '--e' for exit\n ")
	err = client.Initial("127.0.0.1:10137", time.Second * 2)
	if err != nil {
		glog.Errorf("client initial failed : %v", err)
	}
	defer client.Close()
	fmt.Print("\nInput Name :")
	_,err = fmt.Scanf("%s", &name)
	if err != nil {
		glog.Errorf("client get name failed : %v", err)
	}
	for {
		var (
			op string
			target string
			content string
			pack *chatMsg.MsgPack
		)
		_, err = fmt.Scanf("%s", &op)
		if op == "--s" {
			_, err = fmt.Scanf("%s %s", &target , &content)
			err = client.SendTo(&chatMsg.Msg{
				From: name,
				Target: target,
				Content: []byte(content),
				MsgContentType: chatMsg.MsgContentType_Text,
				MsgType: chatMsg.MsgType_Single,
			})
			if err != nil {

			}
		} else if op == "--p" {
			pack, err = client.Pull(&chatMsg.Ack{From: name})
			if err != nil {
				glog.Errorf("Pull failed : %v", err)
				continue
			}
			for i := 0; i < len(pack.MsgList); i ++ {
				fmt.Println(pack.MsgList[i])
			}
		} else if op == "--e" {
			break
		} else {
			continue
		}
	}
}