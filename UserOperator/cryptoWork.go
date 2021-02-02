package UserOperator

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/dxyinme/LukaClient/db"
	"github.com/dxyinme/LukaComm/chatMsg"
	utilCrypto "github.com/dxyinme/LukaComm/util/crypto"
	"io/ioutil"
	"log"
)

func GetNowUserPrivateKey() (cryptor *utilCrypto.Cryptor) {
	if NowLoginUser == nil {
		log.Println("GetNowUserPrivateKey: No Login Status User")
		return nil
	}
	uid := NowLoginUser.Name
	pemFileByte, err := ioutil.ReadFile(priKeyLoad + uid + ".pem")
	if err != nil {
		log.Println("GetNowUserPrivateKey: read pemFile error")
		return nil
	}
	block, _ := pem.Decode(pemFileByte)
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Printf("GetNowUserPrivateKey: %v", err)
		return nil
	}
	cryptor = &utilCrypto.Cryptor{}
	cryptor.SetPriKey(priKey)
	return
}

func decodeAESCipherText(msg *chatMsg.Msg) {
	senderInfo, err := db.GetUserInfoByUid(msg.From)
	if err != nil {
		log.Println(err)
	}
	aesCrypto := utilCrypto.NewAESCrypto(senderInfo.AESKey)
	msg.Content, err = aesCrypto.DecodeCBC(msg.Content)
	if err != nil {
		log.Println(err)
	}
}

func encodeAESPlainText(msg *chatMsg.Msg) {
	recverInfo, err := db.GetUserInfoByUid(msg.Target)
	if err != nil {
		log.Println(err)
	}
	aesCrypto := utilCrypto.NewAESCrypto(recverInfo.AESKey)
	msg.Content, err = aesCrypto.EncodeCBC(msg.Content)
	if err != nil {
		log.Println(err)
	}
}