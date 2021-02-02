package db

import (
	"github.com/dxyinme/LukaClient/IpcMsg"
	"log"
	"testing"
)

func TestUserInfo(t *testing.T) {
	err := NewConn("test.db")
	if err != nil {
		t.Error(err)
	}

	_, err = ExecCmd(CREATE_USERINFO_TABLE)
	if err != nil {
		t.Error(err)
	}
	AESKey := []byte("0348fflleeew")
	if err != nil {
		t.Fatal(err)
	}
	test := &IpcMsg.UserInfo{
		Uid: "test_uid",
		Name: "test_name",
		AESKey: AESKey,
	}
	err = SaveUserInfo(test)
	if err != nil {
		t.Fatal(err)
	}
	res,err := GetUserInfoByUid("test_uid")
	if err != nil {
		t.Fatal(err)
	}
	if string(res.AESKey) != string(test.AESKey) || res.Name != test.Name {
		t.Fatal(err)
	}
	log.Printf("resName=%s, testName=%s\n",res.Name, test.Name)
	test.Uid = "qqqi"
	err = SaveUserInfo(test)
	if err != nil {
		t.Fatal(err)
	}
	test.AESKey = []byte("hehehehehehehehe")
	err = UpdateUserInfoByUid(test)
	if err != nil {
		t.Fatal(err)
	}
}