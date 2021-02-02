package IpcMsg

type Secret struct {
	From string
	Target string
	AESKey []byte
}
