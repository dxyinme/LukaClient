package db

import (
	"database/sql"
	"fmt"
	"github.com/dxyinme/LukaClient/IpcMsg"
)

func parseResIpcMsgUserInfo(res *sql.Rows) ([]*IpcMsg.UserInfo,error) {
	var (
		userInfoList = make([]*IpcMsg.UserInfo, 0)
		err error
	)
	for res.Next() {
		var (
			UID string
			Name string
			AESKey []byte
		)
		err = res.Scan(&UID, &Name, &AESKey)
		if err != nil {
			return userInfoList, err
		}
		userInfoList = append(userInfoList, &IpcMsg.UserInfo{
			Uid:    UID,
			Name:   Name,
			AESKey: AESKey,
		})
	}
	return userInfoList, err
}

func SaveUserInfo(userInfo *IpcMsg.UserInfo) (err error) {
	var (
		stmt *sql.Stmt
	)
	_, err = GetUserInfoByUid(userInfo.Uid)
	if err == nil {
		return fmt.Errorf("uid [%s], is existed", userInfo.Uid)
	}
	stmt, err = databaseConnector.Prepare(INSERT_USERINFO)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userInfo.Uid, userInfo.Name, userInfo.AESKey)
	if err != nil {
		return err
	}
	return nil
}

func GetUserInfoByUid(uid string) (userInfo *IpcMsg.UserInfo, err error) {
	var (
		stmt *sql.Stmt
		res *sql.Rows
	)
	stmt, err = databaseConnector.Prepare(SELECT_USERINFO_BY_UID)
	if err != nil {
		return nil, err
	}
	res, err = stmt.Query(uid)
	userInfoList, err := parseResIpcMsgUserInfo(res)
	if err != nil {
		return nil, err
	}
	if len(userInfoList) == 1 {
		return userInfoList[0], nil
	} else if len(userInfoList) == 0 {
		return nil, fmt.Errorf("uid [%s], is not saved in here", uid)
	} else {
		return nil, fmt.Errorf("uid [%s], exist time is [%d]", uid, len(userInfoList))
	}
}

func UpdateUserInfoByUid(userInfo *IpcMsg.UserInfo) (err error) {
	var (
		stmt *sql.Stmt
	)
	stmt, err = databaseConnector.Prepare(UPDATE_USERINFO_BY_UID)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userInfo.Uid, userInfo.Name, userInfo.AESKey, userInfo.Uid)
	if err != nil {
		return err
	}
	return nil
}