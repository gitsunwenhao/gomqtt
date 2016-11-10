package protocol

import (
	"encoding/binary"
	"fmt"
)

// 固定报头
// 第一个字节是控制报文，后面是剩余长度，该长度是变长的,1-4个字节
type header struct {
	// 剩余长度 = 可变报头 ＋ payload
	remLen int32

	// 控制报文：前4位是报文类型，后四位是flag的保留位
	typeFlag byte

	// 部分消息是需要packet ID的,2个字节表示uint16,大端表示

	packetID uint16
}

// Header的打印形式
func (h header) String() string {
	return fmt.Sprintf("Type=%q, Flags=%08b, Remaining Length=%d",
		h.Type().Name(), h.Flags(), h.remLen)
}

// 返回报文类型名的字符串表示
func (h *header) Name() string {
	return h.Type().Name()
}

// 返回报文类型的描述信息
func (h *header) Desc() string {
	return h.Type().Desc()
}

// 返回报文类型
func (h *header) Type() PacketType {
	return PacketType(h.typeFlag >> 4)
}

// 设置报文的控制类型，同时设置默认的flag位
// 注意这里不更新updated标志,因此encode后的包体大小不会发生变化
func (h *header) SetType(t PacketType) error {
	if !t.Valid() {
		return fmt.Errorf("header/SetType1: Invalid control packet type %d", t)
	}

	// 这里设置了报文的类型和Flag，Flag是取的默认的
	h.typeFlag = byte(t)<<4 | (t.DefaultFlags() & 0xf)

	return nil
}

// 返回flags
func (h *header) Flags() byte {
	return h.typeFlag & 0x0f
}

// 剩余长度，报文中除去固定报头的剩余部分
func (h *header) RemainingLength() int32 {
	return h.remLen
}

// 设置剩余长度，不得超过268435455字节，这个也是mqtt报文最大的长度
func (h *header) SetRemainingLength(l int32) error {
	if l > maxRemainingLength || l < 0 {
		return fmt.Errorf("header/SetLength1: Remaining length (%d) out of bound (max %d, min 0)",
			l, maxRemainingLength)
	}

	h.remLen = l

	return nil
}

func (h *header) Len() int {
	return h.msglen()
}

// 返回PacketID
func (h *header) PacketID() uint16 {
	return h.packetID
}

func (h *header) SetPacketID(id uint16) {
	h.packetID = id
}

//编码固定报头
func (h *header) encode(dst []byte) (int, error) {
	ml := h.msglen()

	if len(dst) < ml {
		return 0, fmt.Errorf("header/Encode1: Insufficient buffer size. Expecting %d, got %d.", ml, len(dst))
	}

	total := 0

	if h.remLen > maxRemainingLength || h.remLen < 0 {
		return total, fmt.Errorf("header/Encode2: Remaining length (%d) out of bound (max %d, min 0)", h.remLen, maxRemainingLength)
	}

	if !h.Type().Valid() {
		return total, fmt.Errorf("header/Encode3: Invalid message type %d", h.Type())
	}

	//第一个字节为控制报文
	dst[total] = h.typeFlag
	total += 1

	//写入剩余长度，采用了变长编码
	n := binary.PutUvarint(dst[total:], uint64(h.remLen))
	total += n

	return total, nil
}

//解码固定报头
func (h *header) decode(src []byte) (int, error) {
	total := 0

	mtype := h.Type()

	//读取并写入控制报文
	h.typeFlag = src[total]

	if !mtype.Valid() {
		return total, fmt.Errorf("header/Decode1: Invalid message type %d.", mtype)
	}

	//只有PUBLISH才有Flag位，其它的Flag都是默认的
	if mtype != PUBLISH && h.Flags() != mtype.DefaultFlags() && h.Flags() != mtype.DefaultFlags10() {
		return total, fmt.Errorf("header/Decode3: Invalid message (%d) flags. Expecting %d, got %d",
			mtype, mtype.DefaultFlags(), h.Flags())
	}

	//验证PUBLISH报文的qos
	if mtype == PUBLISH && !ValidQos((h.Flags()>>1)&0x3) {
		return total, fmt.Errorf("header/Decode4: Invalid QoS (%d) for PUBLISH message.",
			(h.Flags()>>1)&0x3)
	}
	total++

	//读取剩余长度(使用了变长编码的读取函数)
	rl, m := binary.Uvarint(src[total:])
	total += m
	h.remLen = int32(rl)

	if h.remLen > maxRemainingLength || rl < 0 {
		return total, fmt.Errorf("header/Decode5: Remaining length (%d) out of bound (max %d, min 0)", h.remLen, maxRemainingLength)
	}

	if int(h.remLen) > len(src[total:]) {
		return total, fmt.Errorf("header/Decode6: Remaining length (%d) is greater than remaining buffer (%d)", h.remLen, len(src[total:]))
	}

	return total, nil
}

func (h *header) msglen() int {
	total := 1

	if h.remLen <= 127 {
		total += 1
	} else if h.remLen <= 16383 {
		total += 2
	} else if h.remLen <= 2097151 {
		total += 3
	} else {
		total += 4
	}

	return total
}
