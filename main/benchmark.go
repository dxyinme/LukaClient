package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	CynicUClient "github.com/dxyinme/LukaComm/CynicU/Client"
	"github.com/dxyinme/LukaComm/chatMsg"
)

var (
	ipPort   = flag.String("ipport", "127.0.0.1:10137", "server IP:PORT")
	testFile = flag.String("testFile", "test/normal.test", "test case file")
)

// Operator content for operator
// send : s,name,target,content
// pull : p,name
type Operator struct {
	op      string
	name    string
	content string
	target  string
}

func main() {
	flag.Parse()
	startTime := time.Now()
	client := &CynicUClient.Client{}
	err := client.Initial(*ipPort, time.Second*3)
	var (
		sendCnt int = 0
		pullCnt int = 0
		wg      sync.WaitGroup
	)
	file, err := os.Open(*testFile)
	if err != nil {
		log.Println(err)
	}
	reader := bufio.NewReader(file)
	for {
		line, errRd := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		opNow := strings.Split(line, ",")
		if len(opNow) != 2 && len(opNow) != 4 {
			break
		}

		if errRd == nil {
			if opNow[0] == "s" {
				sendCnt++
				wg.Add(1)
				go func() {
					nowMsg := &chatMsg.Msg{
						From:           opNow[1],
						Target:         opNow[2],
						Content:        []byte(opNow[3]),
						MsgType:        chatMsg.MsgType_Single,
						MsgContentType: chatMsg.MsgContentType_Text,
					}
					log.Printf("from [%s], target [%s]", nowMsg.From, nowMsg.Target)
					errSend := client.SendTo(nowMsg)
					if errSend != nil {
						log.Println(errSend)
					}
					wg.Done()
				}()
			} else {
				pullCnt++
				log.Printf("pull [%s]", opNow[1])
				packNow, err := client.PullAll(&chatMsg.PullReq{
					From: opNow[1],
				})
				if err != nil {
					log.Println(err)
				}
				for _, v := range packNow.MsgList {
					log.Printf("from %s, target %s , content %s", v.From, v.Target, string(v.Content))
				}
			}
		} else {
			break
		}
	}
	wg.Wait()
	log.Printf("[cost time] : %d ms , [operator] send : %d , pull : %d", (time.Now().Sub(startTime).Milliseconds()), sendCnt, pullCnt)
}
