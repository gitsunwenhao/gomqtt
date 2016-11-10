package protocol

import (
	"encoding/binary"
	"fmt"
)

// PUBACK包是在QOS=1时对PUBLISH包的回应
type PubackPacket struct {
	header
}

func NewPubackPacket() *PubackPacket {
	pp := &PubackPacket{}
	pp.SetType(PUBACK)

	return pp
}

func (pp PubackPacket) String() string {
	return fmt.Sprintf("%s, Packet ID=%d", pp.header, pp.packetID)
}

func (pp *PubackPacket) Len() int {
	return pp.header.msglen() + pp.msglen()
}

func (pp *PubackPacket) Decode(src []byte) (int, error) {
	total := 0

	//Decode出固定报头
	n, err := pp.header.decode(src[total:])
	total += n
	if err != nil {
		return total, err
	}

	//2字节的pakcetId
	pp.packetID = binary.BigEndian.Uint16(src[total : total+2])
	total += 2

	return total, nil
}

func (pp *PubackPacket) Encode() (int, []byte, error) {

	// hl := pp.header.msglen()
	ml := pp.msglen()
	// if len(dst) < hl+ml {
	// 	return 0, fmt.Errorf("puback/Encode: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
	// }

	if err := pp.SetRemainingLength(int32(ml)); err != nil {
		return 0, nil, err
	}

	dst := make([]byte, pp.Len())
	total := 0

	n, err := pp.header.encode(dst[total:])
	total += n
	if err != nil {
		return 0, nil, err
	}

	binary.BigEndian.PutUint16(dst[total:total+2], pp.packetID)
	total += 2

	return total, dst, nil
}

func (pp *PubackPacket) msglen() int {
	// 这里的可变报文，仅仅包含PacketId
	return 2
}
