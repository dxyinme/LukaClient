package db

import (
	"database/sql"
	"github.com/dxyinme/LukaComm/chatMsg"
)


func parseResChatMsg(res *sql.Rows) (ret []*chatMsg.Msg, err error){
	ret = make([]*chatMsg.Msg,0)
	for res.Next() {
		var (
			msgId          string
			msgType        int
			msgContentType int
			content        string
			sendTime       string
			msgFrom        string
			msgTarget      string
			groupName      string
		)
		err = res.Scan(&msgId, &msgType, &msgContentType,
			&content, &sendTime, &msgFrom, &msgTarget, &groupName)
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
	_, err = stmt.Exec(x.MsgId,x.MsgType,x.MsgContentType,
		x.Content,x.SendTime,x.From,x.Target,x.GroupName)
	return
}

func LoadChatMsg(from string, isGroup bool, limit int) (ret []*chatMsg.Msg, err error) {
	var (
		sqLine string
		stmt *sql.Stmt
		res *sql.Rows
		)
	if isGroup {
		sqLine = SELECT_MSG_BY_GROUP
	} else {
		sqLine = SELECT_MSG_BY_FROM
	}
	stmt, err = databaseConnector.Prepare(sqLine)
	if err != nil {
		return nil, err
	}
	res, err = stmt.Query(from, limit)
	if err != nil {
		return nil, err
	}
	ret, err = parseResChatMsg(res)
	return
}

func LoadChatMsgAll(from string, isGroup bool) (ret []*chatMsg.Msg, err error) {
	var (
		sqLine string
		stmt *sql.Stmt
		res *sql.Rows
	)
	if isGroup {
		sqLine = SELECT_MSG_BY_GROUP_ALL
	} else {
		sqLine = SELECT_MSG_BY_FROM_ALL
	}
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