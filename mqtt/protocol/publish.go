package protocol

import (
	"encoding/binary"
	"fmt"
)

// PUBLISH报文可以从客户端发向服务器，也可以从服务器发向客户端,用于向指定的topic发布消息
type PublishPacket struct {
	header

	topic   []byte
	payload []byte
}

func NewPublishPacket() *PublishPacket {
	pp := &PublishPacket{}
	pp.SetType(PUBLISH)

	return pp
}

func (pp PublishPacket) String() string {
	return fmt.Sprintf("%s, Topic=%q, Packet ID=%d, QoS=%d, Retained=%t, Dup=%t, Payload=%v",
		pp.header, pp.topic, pp.packetID, pp.QoS(), pp.Retain(), pp.Dup(), pp.payload)
}

// Dup返回一个PUBLISH报文是否是重复投递
// 如果报文控制标志中的DUP flag被设置为0，表示该报文是第一次发送。如果设置为1，表示该报文是再一次投递的
func (pp *PublishPacket) Dup() bool {
	return ((pp.Flags() >> 3) & 0x1) == 1
}

func (pp *PublishPacket) SetDup(v bool) {
	if v {
		pp.typeFlag |= 0x8 // 00001000
	} else {
		pp.typeFlag &= 247 // 11110111
	}
}

// RETAIN标志位只能用在PUBLISH报文中。如果该flag被设置为1且报文是从客户端发向服务器的，服务器必须要存储
// 该消息和QoS等级，然后等待订阅者来订阅此条消息
func (pp *PublishPacket) Retain() bool {
	return (pp.Flags() & 0x1) == 1
}

func (pp *PublishPacket) SetRetain(v bool) {
	if v {
		pp.typeFlag |= 0x1 // 00000001
	} else {
		pp.typeFlag &= 254 // 11111110
	}
}

// QoS是消息质量保证等级，合法的值包括：QosAtMostOnce, QosAtLeastOnce and QosExactlyOnce 。
func (pp *PublishPacket) QoS() byte {
	return (pp.Flags() >> 1) & 0x3
}

func (pp *PublishPacket) SetQoS(v byte) error {
	if v != 0x0 && v != 0x1 && v != 0x2 {
		return fmt.Errorf("publish/SetQoS: Invalid QoS %d.", v)
	}

	pp.typeFlag = (pp.typeFlag & 249) | (v << 1) // 249 = 11111001

	return nil
}

// Topic是一个频道，PUBLISH消息会被发布到指定的频道中,然后客户端可以对频道进行订阅
func (pp *PublishPacket) Topic() []byte {
	return pp.topic
}

func (pp *PublishPacket) SetTopic(v []byte) error {
	if !ValidTopic(v) {
		return fmt.Errorf("publish/SetTopic: Invalid topic name (%s). Must not be empty or contain wildcard characters", string(v))
	}
	pp.topic = v

	return nil
}

// payload是publish的消息报体部分,对发布者来说，就是应用消息
// 数据内容和格式是发布者自定义的，payload的长度这样计算：固定报头中的剩余长度 - 可变报头的长度。
// 包含零长度的PUBLISH是合法的
func (pp *PublishPacket) Payload() []byte {
	return pp.payload
}

func (pp *PublishPacket) SetPayload(v []byte) {
	pp.payload = v
}

func (pp *PublishPacket) Len() int {
	return pp.header.msglen() + pp.msglen()
}

func (pp *PublishPacket) Decode(src []byte) (int, error) {
	total := 0

	// 解码固定报头
	hn, err := pp.header.decode(src[total:])
	if err != nil {
		return total, err
	}
	total += hn

	// 解码topic
	n := 0
	pp.topic, n, err = readLPBytes(src[total:])
	if err != nil {
		return total, err
	}
	total += n

	if !ValidTopic(pp.topic) {
		return total, fmt.Errorf("publish/Decode: Invalid topic name (%s). Must not be empty or contain wildcard characters", string(pp.topic))
	}

	// 只有QoS 1或2时，才有packetID
	if pp.QoS() != 0 {
		pp.packetID = binary.BigEndian.Uint16(src[total : total+2])
		total += 2
	}

	// 解码payload
	// payload长度 = 剩余长度 －可变报头长度
	l := int(pp.remLen) - (total - hn)
	pp.payload = src[total : total+l]
	total += l

	return total, nil
}

func (pp *PublishPacket) Encode() (int, []byte, error) {

	if len(pp.topic) == 0 {
		return 0, nil, fmt.Errorf("publish/Encode: Topic name is empty.")
	}

	if len(pp.payload) == 0 {
		return 0, nil, fmt.Errorf("publish/Encode: Payload is empty.")
	}

	ml := pp.msglen()

	if err := pp.SetRemainingLength(int32(ml)); err != nil {
		return 0, nil, err
	}

	// hl := pp.header.msglen()

	// if len(dst) < hl+ml {
	// 	return 0, nil, fmt.Errorf("publish/Encode: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
	// }
	dst := make([]byte, pp.Len())

	total := 0

	n, err := pp.header.encode(dst[total:])
	total += n
	if err != nil {
		return 0, nil, err
	}

	n, err = writeLPBytes(dst[total:], pp.topic)
	total += n
	if err != nil {
		return 0, nil, err
	}

	//QoS不为0时，必须要传PacketID
	if pp.QoS() != 0 {
		if pp.PacketID() == 0 {
			return 0, nil, fmt.Errorf("publish/Encode: invalid packetid %d when qos == 0", pp.PacketID())
		}

		binary.BigEndian.PutUint16(dst[total:total+2], pp.packetID)

		total += 2

	}

	copy(dst[total:], pp.payload)
	total += len(pp.payload)

	return total, dst, nil
}

func (pp *PublishPacket) msglen() int {
	total := 2 + len(pp.topic) + len(pp.payload)
	if pp.QoS() != 0 {
		total += 2
	}

	return total
}
