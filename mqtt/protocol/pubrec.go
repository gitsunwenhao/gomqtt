package protocol

import (
	"encoding/binary"
	"fmt"
)

// PUBREC包是对PUBLISH包的回应，是QoS 2的第二步
type PubrecPacket struct {
	header
}

func NewPubrecPacket() *PubrecPacket {
	pp := &PubrecPacket{}
	pp.SetType(PUBREC)

	return pp
}

func (pp PubrecPacket) String() string {
	return fmt.Sprintf("%s, Packet ID=%d", pp.header, pp.packetID)
}

func (pp *PubrecPacket) Len() int {
	return pp.header.msglen() + pp.msglen()
}

func (pp *PubrecPacket) Decode(src []byte) (int, error) {
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

func (pp *PubrecPacket) Encode() (int, []byte, error) {
	//	hl := pp.header.msglen()
	ml := pp.msglen()

	// if len(dst) < hl+ml {
	// 	return 0, nil, fmt.Errorf("puback/Encode: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
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

func (pp *PubrecPacket) msglen() int {
	// 这里的可变报文，仅仅包含PacketId
	return 2
}
