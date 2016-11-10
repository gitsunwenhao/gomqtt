package protocol

import (
	"encoding/binary"
	"fmt"
)

// UNSUBACK是服务器对UNSUBSCRIBE的响应
type UnsubackPacket struct {
	header
}

// NewUnsubackMessage creates a new UNSUBACK message.
func NewUnsubackPacket() *UnsubackPacket {
	up := &UnsubackPacket{}
	up.SetType(UNSUBACK)

	return up
}

func (up UnsubackPacket) String() string {
	return fmt.Sprintf("%s, Packet ID=%d", up.header, up.packetID)
}

func (up *UnsubackPacket) Len() int {
	return up.header.msglen() + up.msglen()
}

func (up *UnsubackPacket) Decode(src []byte) (int, error) {
	total := 0

	//Decode出固定报头
	n, err := up.header.decode(src[total:])
	total += n
	if err != nil {
		return total, err
	}

	//2字节的pakcetId
	up.packetID = binary.BigEndian.Uint16(src[total : total+2])
	total += 2

	return total, nil
}

func (up *UnsubackPacket) Encode() (int, []byte, error) {
	// hl := up.header.msglen()
	ml := up.msglen()

	// if len(dst) < hl+ml {
	// 	return 0, fmt.Errorf("puback/Encode: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
	// }

	if err := up.SetRemainingLength(int32(ml)); err != nil {
		return 0, nil, err
	}
	dst := make([]byte, up.Len())
	total := 0

	n, err := up.header.encode(dst[total:])
	total += n
	if err != nil {
		return 0, nil, err
	}

	binary.BigEndian.PutUint16(dst[total:total+2], up.packetID)
	total += 2

	return total, dst, nil
}

func (up *UnsubackPacket) msglen() int {
	// 这里的可变报文，仅仅包含PacketId
	return 2
}
