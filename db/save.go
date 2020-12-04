package db

import (
	"database/sql"
	"github.com/dxyinme/LukaComm/chatMsg"
	"time"
)


func parseResChatMsg(res *sql.Rows) (ret []*chatMsg.Msg, err error){
	ret = make([]*chatMsg.Msg,0)
	for res.Next() {
		var (
			msgId          	string
			msgType        	int
			msgContentType 	int
			content        	string
			sendTime       	string
			msgFrom        	string
			msgTarget      	string
			groupName      	string
			recvTime		int64
		)
		err = res.Scan(&msgId, &msgType, &msgContentType,
			&content, &sendTime, &msgFrom, &msgTarget, &groupName, &recvTime)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &chatMsg.Msg{
			From:           msgFrom,
			Target:         msgTarget,
			Content:        []byte(content),
			MsgType:        chatMsg.MsgType(msgType),
			MsgContentType: chatMsg.MsgContentType(msgContentType),
			SendTime:       sendTime,
			GroupName:      groupName,
			Spread:         false,
			MsgId:          msgId,
		})
	}
	return
}

func SaveChatMsg(x *chatMsg.Msg) (err error) {
	var (
		stmt *sql.Stmt
	)
	stmt,err = databaseConnector.Prepare(INSERT_MSG)
	if err != nil {
		return
	}
	recvTime := time.Now().Unix()
	_, err = stmt.Exec(x.MsgId,x.MsgType,x.MsgContentType,
		x.Content,x.SendTime,x.From,x.Target,x.GroupName,recvTime)
	return
}

func LoadGroupChatMsgAll(from string) (ret []*chatMsg.Msg, err error) {
	var (
		sqLine string
		stmt *sql.Stmt
		res *sql.Rows
	)
	sqLine = SELECT_MSG_BY_GROUP_ALL
	stmt, err = databaseConnector.Prepare(sqLine)
	if err != nil {
		return nil, err
	}
	res, err = stmt.Query(from)
	if err != nil {
		return nil, err
	}
	ret, err = parseResChatMsg(res)
	return
}

func LoadSingleChatMsgAll(from ,target string) (ret []*chatMsg.Msg, err error) {
	var (
		sqLine string
		stmt *sql.Stmt
		res *sql.Rows
	)
	sqLine = SELECT_MSG_SINGLE_ALL
	stmt, err = databaseConnector.Prepare(sqLine)
	if err != nil {
		return nil, err
	}
	res, err = stmt.Query(from, target, target, from)
	if err != nil {
		return nil, err
	}
	ret, err = parseResChatMsg(res)
	return
}