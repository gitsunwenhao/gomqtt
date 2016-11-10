package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// UNSUBSCRIBE是客户端发送的取消订阅的报文
type UnsubscribePacket struct {
	header

	topics [][]byte
}

func NewUnsubscribePacket() *UnsubscribePacket {
	up := &UnsubscribePacket{}
	up.SetType(UNSUBSCRIBE)

	return up
}

func (up UnsubscribePacket) String() string {
	msgstr := fmt.Sprintf("%s", up.header)

	for i, t := range up.topics {
		msgstr = fmt.Sprintf("%s, Topic%d=%s", msgstr, i, string(t))
	}

	return msgstr
}

func (up *UnsubscribePacket) Topics() [][]byte {
	return up.topics
}

func (up *UnsubscribePacket) AddTopic(topic []byte) {
	if up.TopicExists(topic) {
		return
	}

	up.topics = append(up.topics, topic)
}

func (up *UnsubscribePacket) RemoveTopic(topic []byte) {
	for i, t := range up.topics {
		if bytes.Equal(t, topic) {
			up.topics = append(up.topics[:i], up.topics[i+1:]...)
			break
		}
	}
}

func (up *UnsubscribePacket) TopicExists(topic []byte) bool {
	for _, t := range up.topics {
		if bytes.Equal(t, topic) {
			return true
		}
	}

	return false
}

func (up *UnsubscribePacket) Len() int {
	return up.header.msglen() + up.msglen()
}

func (up *UnsubscribePacket) Decode(src []byte) (int, error) {
	total := 0

	hn, err := up.header.decode(src[total:])
	total += hn
	if err != nil {
		return total, err
	}

	up.packetID = binary.BigEndian.Uint16(src[total : total+2])
	total += 2

	rl := int(up.remLen) - (total - hn)
	for rl > 0 {
		t, n, err := readLPBytes(src[total:])
		total += n
		if err != nil {
			return total, err
		}

		up.topics = append(up.topics, t)
		rl = rl - n - 1
	}

	if len(up.topics) == 0 {
		return 0, fmt.Errorf("unsubscribe/Decode: Empty topic list")
	}

	return total, nil
}

func (up *UnsubscribePacket) Encode() (int, []byte, error) {
	//hl := up.header.msglen()
	ml := up.msglen()

	// if len(dst) < hl+ml {
	// 	return 0, fmt.Errorf("unsubscribe/Encode: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
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

	// UNSUBSCRIBE必须要有PackeID
	if up.PacketID() == 0 {
		return 0, nil, fmt.Errorf("subscribe/Encode: invalid packetid %d", up.PacketID())
	}

	binary.BigEndian.PutUint16(dst[total:total+2], up.packetID)
	total += 2

	for _, t := range up.topics {
		n, err := writeLPBytes(dst[total:], t)
		total += n
		if err != nil {
			return 0, nil, err
		}
	}

	return total, dst, nil
}

func (up *UnsubscribePacket) msglen() int {
	// packet ID
	total := 2

	for _, t := range up.topics {
		total += 2 + len(t)
	}

	return total
}
