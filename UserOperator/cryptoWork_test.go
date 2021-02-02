package UserOperator

import (
	"github.com/dxyinme/LukaClient/IpcMsg"
	"github.com/dxyinme/LukaClient/db"
	CynicUClient "github.com/dxyinme/LukaComm/CynicU/Client"
	"github.com/dxyinme/LukaComm/chatMsg"
	"github.com/dxyinme/LukaComm/util/crypto"
	"log"
	"os"
	"testing"
	"time"
)

// mention!!
// please test this under AuthMain are listening
func TestGetNowUserPrivateKey(t *testing.T) {
	client = &CynicUClient.Client{}
	client.SetTimeout(10 * time.Second)
	_ = os.Mkdir(priKeyLoad, os.ModePerm)
	err := TryGenRSAKey("test")
	if err != nil {
		t.Fatal(err)
	}
	pubKey, err := client.GetAuthPubKey("localhost:20020", "test")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(len(pubKey))
	cipherText, err := crypto.EncodePub([]byte("xixiqqxixi"), pubKey)
	if err != nil {
		t.Fatal(err)
	}
	NowLoginUser = &IpcMsg.Login{
		Name:     "test",
		Password: "",
	}
	cryptor := GetNowUserPrivateKey()
	plainText, err := cryptor.Decode(cipherText)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(string(plainText))
	if string(plainText) != "xixiqqxixi" {
		t.Fatal("no equal")
	}
	_ = os.RemoveAll(priKeyLoad)
}

func TestAES(t *testing.T) {
	err := db.NewConn("test.db")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecCmd(db.CREATE_USERINFO_TABLE)
	if err != nil {
		t.Fatal(err)
	}
	err = db.SaveUserInfo(&IpcMsg.UserInfo{
		Uid:    "test",
		Name:   "test",
		AESKey: []byte("1234567887654321"),
	})
	if err != nil {
		t.Fatal(err)
	}
	msg := &chatMsg.Msg{
		From:           "test",
		Target:         "test",
		Content:        []byte("hello"),
		MsgType:        chatMsg.MsgType_Single,
		MsgContentType: chatMsg.MsgContentType_Text,
		SecretLevel:    1,
		MsgId:          "",
	}
	encodeAESPlainText(msg)
	log.Println(string(msg.Content))
	decodeAESCipherText(msg)
	log.Println(string(msg.Content))
	if string(msg.Content) != "hello" {
		t.Fatal("no equal")
	}
	_ = os.Remove("test.db")
}