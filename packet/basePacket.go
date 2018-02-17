package packet

type BaseMessage interface {
	Init([]byte) error
	Make([]byte) []byte
	Verify()error
	Body()[]byte
	NeedReply()bool
}

type BasePacket interface {
	Init([]byte) error
	First()([]byte, error)
	Next()([]byte, error)
}
