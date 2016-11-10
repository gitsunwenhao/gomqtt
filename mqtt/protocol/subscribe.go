package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// SUBSCRIBE从客户端发向服务器器。每个订阅可以订阅一个或者多个topic。服务器通过PUBLISH将消息发布到topic中
type SubscribePacket struct {
	header

	topics [][]byte
	qos    []byte
}

func NewSubscribePacket() *SubscribePacket {
	sp := &SubscribePacket{}
	sp.SetType(SUBSCRIBE)

	return sp
}

func (sp SubscribePacket) String() string {
	msgstr := fmt.Sprintf("%s, Packet ID=%d", sp.header, sp.PacketID())

	for i, t := range sp.topics {
		msgstr = fmt.Sprintf("%s, Topic[%d]=%q/%d", msgstr, i, string(t), sp.qos[i])
	}

	return msgstr
}

// Topic是sub/pub的频道
func (sp *SubscribePacket) Topics() [][]byte {
	return sp.topics
}

func (sp *SubscribePacket) AddTopic(topic []byte, qos byte) error {
	if !ValidQos(qos) {
		return fmt.Errorf("Invalid QoS %d", qos)
	}

	for i, t := range sp.topics {
		//若topic已存在，更新qos
		if bytes.Equal(t, topic) {
			sp.qos[i] = qos
			return nil
		}
	}

	sp.topics = append(sp.topics, topic)
	sp.qos = append(sp.qos, qos)
	return nil
}

func (sp *SubscribePacket) RemoveTopic(topic []byte) {
	for i, t := range sp.topics {
		if bytes.Equal(t, topic) {
			sp.topics = append(sp.topics[:i], sp.topics[i+1:]...)
			sp.qos = append(sp.qos[:i], sp.qos[i+1:]...)
			break
		}
	}
}

func (sp *SubscribePacket) TopicExists(topic []byte) bool {
	for _, t := range sp.topics {
		if bytes.Equal(t, topic) {
			return true
		}
	}

	return false
}

func (sp *SubscribePacket) TopicQos(topic []byte) byte {
	for i, t := range sp.topics {
		if bytes.Equal(t, topic) {
			return sp.qos[i]
		}
	}

	return QosFailure
}

func (sp *SubscribePacket) Qos() []byte {
	return sp.qos
}

func (sp *SubscribePacket) Len() int {
	return sp.header.msglen() + sp.msglen()
}

func (sp *SubscribePacket) Decode(src []byte) (int, error) {
	total := 0

	hn, err := sp.header.decode(src[total:])
	total += hn
	if err != nil {
		return total, err
	}

	sp.packetID = binary.BigEndian.Uint16(src[total : total+2])
	total += 2

	rl := int(sp.remLen) - (total - hn)
	for rl > 0 {
		t, n, err := readLPBytes(src[total:])
		total += n
		if err != nil {
			return total, err
		}

		sp.topics = append(sp.topics, t)

		sp.qos = append(sp.qos, src[total])
		total++

		rl = rl - n - 1
	}

	if len(sp.topics) == 0 {
		return 0, fmt.Errorf("subscribe/Decode: Empty topic list")
	}

	return total, nil
}

func (sp *SubscribePacket) Encode() (int, []byte, error) {
	// hl := sp.header.msglen()
	ml := sp.msglen()

	// if len(dst) < hl+ml {
	// 	return 0, fmt.Errorf("subscribe/Encode: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
	// }

	if err := sp.SetRemainingLength(int32(ml)); err != nil {
		return 0, nil, err
	}
	dst := make([]byte, sp.Len())

	total := 0

	n, err := sp.header.encode(dst[total:])
	total += n
	if err != nil {
		return 0, nil, err
	}

	// SUBSCRIBE报文必须要有PackeId
	if sp.PacketID() == 0 {
		return 0, nil, fmt.Errorf("subscribe/Encode: invalid packetid %d", sp.PacketID())
	}

	binary.BigEndian.PutUint16(dst[total:total+2], sp.packetID)
	total += 2

	for i, t := range sp.topics {
		n, err := writeLPBytes(dst[total:], t)
		total += n
		if err != nil {
			return 0, nil, err
		}

		dst[total] = sp.qos[i]
		total++
	}

	return total, dst, nil
}

func (sp *SubscribePacket) msglen() int {
	// packet ID
	total := 2

	for _, t := range sp.topics {
		total += 2 + len(t) + 1
	}

	return total
}
