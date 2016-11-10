package protocol

// PINGREQ报文从客户端发向服务器,有三个目标
// 1.告诉服务器客户端依旧存活，特别是在没有发送其它报文时
// 2.要求服务器回复一个PINGRESP，以保证服务器是存活的
// 3.保证网络是存活的
type PingreqPacket struct {
	header
}

func NewPingreqPacket() *PingreqPacket {
	pp := &PingreqPacket{}
	pp.SetType(PINGREQ)

	return pp
}

func (pp *PingreqPacket) Decode(src []byte) (int, error) {
	return pp.header.decode(src)
}

// func (pp *PingreqPacket) Encode(dst []byte) (int, error) {
// 	return pp.header.encode(dst)
// }

func (pp *PingreqPacket) Encode() (int, []byte, error) {
	dst := make([]byte, pp.Len())
	n, err := pp.header.encode(dst)
	return n, dst, err
}
