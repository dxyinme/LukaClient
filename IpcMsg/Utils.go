package IpcMsg

import (
	"log"
	"strconv"
	"strings"
)

func ArrayBufferToByteArray(arrayBufferStringPtr *string) []byte {
	var (
		result []byte
		strResult []string
	)
	strResult = strings.Split(*arrayBufferStringPtr,",")
	for _,v := range strResult {
		iter, err := strconv.Atoi(v)
		if err != nil {
			log.Println(err)
		}
		result = append(result, byte(iter))
	}
	return result
}