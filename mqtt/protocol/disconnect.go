package protocol

// DISCONNECT是从客户端发向服务器的最后一个控制报文
type DisconnectPacket struct {
	header
}

func NewDisconnectPacket() *DisconnectPacket {
	dp := &DisconnectPacket{}
	dp.SetType(DISCONNECT)

	return dp
}

func (dp *DisconnectPacket) Decode(src []byte) (int, error) {
	return dp.header.decode(src)
}

func (dp *DisconnectPacket) Encode(dst []byte) (int, error) {
	return dp.header.encode(dst)
}
