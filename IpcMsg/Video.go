package IpcMsg

type Video struct {
	Avid string `json:"avid"`
	// Media bytes content, save as string byte1,byte2,...,byteN
	Media string `json:"media"`
}