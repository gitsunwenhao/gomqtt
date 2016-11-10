package protocol

// PINGRESP包是从服务器发向客户端
type PingrespPacket struct {
	header
}

func NewPingrespPacket() *PingrespPacket {
	pp := &PingrespPacket{}
	pp.SetType(PINGRESP)

	return pp
}

func (pp *PingrespPacket) Decode(src []byte) (int, error) {
	return pp.header.decode(src)
}

// func (pp *PingrespPacket) Encode(dst []byte) (int, error) {
// 	return pp.header.encode(dst)
// }

func (pp *PingrespPacket) Encode() (int, []byte, error) {
	dst := make([]byte, pp.Len())
	n, err := pp.header.encode(dst)
	return n, dst, err
}
