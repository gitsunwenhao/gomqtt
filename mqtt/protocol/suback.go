package protocol

import (
	"encoding/binary"
	"fmt"
)

// SUBACK是服务器对SUBSCRIBE的响应
type SubackPacket struct {
	header

	returnCodes []byte
}

func NewSubackPacket() *SubackPacket {
	sp := &SubackPacket{}
	sp.SetType(SUBACK)

	return sp
}

func (sp SubackPacket) String() string {
	return fmt.Sprintf("%s, Packet ID=%d, Return Codes=%v", sp.header, sp.PacketID(), sp.returnCodes)
}

// 返回码列表，说明了订阅时允许的QoS等级
func (sp *SubackPacket) ReturnCodes() []byte {
	return sp.returnCodes
}

func (sp *SubackPacket) AddReturnCodes(ret []byte) error {
	for _, c := range ret {
		if c != QosAtMostOnce && c != QosAtLeastOnce && c != QosExactlyOnce && c != QosFailure {
			return fmt.Errorf("suback/AddReturnCode: Invalid return code %d. Must be 0, 1, 2, 0x80.", c)
		}

		sp.returnCodes = append(sp.returnCodes, c)
	}

	return nil
}

func (sp *SubackPacket) AddReturnCode(ret byte) error {
	return sp.AddReturnCodes([]byte{ret})
}

func (sp *SubackPacket) Len() int {
	return sp.header.msglen() + sp.msglen()
}

func (sp *SubackPacket) Decode(src []byte) (int, error) {
	total := 0

	// 解码出固定报头
	hn, err := sp.header.decode(src[total:])
	total += hn
	if err != nil {
		return total, err
	}

	// 获取PacketId
	sp.packetID = binary.BigEndian.Uint16(src[total : total+2])
	total += 2

	//获取订阅的返回码
	l := int(sp.remLen) - (total - hn)
	sp.returnCodes = src[total : total+l]
	total += len(sp.returnCodes)

	for i, code := range sp.returnCodes {
		if code != 0x00 && code != 0x01 && code != 0x02 && code != 0x80 {
			return total, fmt.Errorf("suback/Decode: Invalid return code %d for topic %d", code, i)
		}
	}

	return total, nil
}

func (sp *SubackPacket) Encode() (int, []byte, error) {

	for i, code := range sp.returnCodes {
		if code != 0x00 && code != 0x01 && code != 0x02 && code != 0x80 {
			return 0, nil, fmt.Errorf("suback/Encode: Invalid return code %d for topic %d", code, i)
		}
	}

	//hl := sp.header.msglen()
	ml := sp.msglen()

	// if len(dst) < hl+ml {
	// 	return 0, fmt.Errorf("suback/Encode: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
	// }

	if err := sp.SetRemainingLength(int32(ml)); err != nil {
		return 0, nil, err
	}
	dst := make([]byte, sp.Len())
	total := 0

	//编码固定报头
	n, err := sp.header.encode(dst[total:])
	total += n
	if err != nil {
		return 0, nil, err
	}

	//编码PackeID
	binary.BigEndian.PutUint16(dst[total:total+2], sp.packetID)
	total += 2

	//编码返回码
	copy(dst[total:], sp.returnCodes)
	total += len(sp.returnCodes)

	return total, dst, nil
}

func (sp *SubackPacket) msglen() int {
	return 2 + len(sp.returnCodes)
}
