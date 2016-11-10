package protocol

import (
	"encoding/binary"
	"fmt"
)

// PUBCOMP是对PUBREL的回应，它是QoS 2交换中的第四步也是最后一步
type PubcompPacket struct {
	header
}

func NewPubcompPacket() *PubcompPacket {
	pp := &PubcompPacket{}
	pp.SetType(PUBCOMP)

	return pp
}

func (pp PubcompPacket) String() string {
	return fmt.Sprintf("%s, Packet ID=%d", pp.header, pp.packetID)
}

func (pp *PubcompPacket) Len() int {
	return pp.header.msglen() + pp.msglen()
}

func (pp *PubcompPacket) Decode(src []byte) (int, error) {
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

func (pp *PubcompPacket) Encode() (int, []byte, error) {

	//hl := pp.header.msglen()
	ml := pp.msglen()

	// if len(dst) < hl+ml {
	// 	return 0, fmt.Errorf("puback/Encode2: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
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

func (pp *PubcompPacket) msglen() int {
	// 这里的可变报文，仅仅包含PacketId
	return 2
}
