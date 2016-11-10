package protocol

import "fmt"

// CONNACK包是在客户端发送CONNECT包到服务器后，服务器返回的确认包。服务器发送客户端的第一个包必须是CONNACK
// 如果客户端在一定时间内没有收到CONNACK包，应该关必网络连接。这个时间的设置取决于应用类型和底层的通信基础设施
type ConnackPacket struct {
	header

	sessionPresent bool
	returnCode     ConnackCode
}

//创建Connack包
func NewConnackPacket() *ConnackPacket {
	cp := &ConnackPacket{}
	cp.SetType(CONNACK)

	return cp
}

func (cp ConnackPacket) String() string {
	return fmt.Sprintf("%s, Session Present=%t, Return code=%q\n", cp.header,
		cp.sessionPresent, cp.returnCode)
}

// SessionPresent,是否存在之前的会话
func (cp *ConnackPacket) SessionPresent() bool {
	return cp.sessionPresent
}

func (cp *ConnackPacket) SetSessionPresent(p bool) {
	if p {
		cp.sessionPresent = true
	} else {
		cp.sessionPresent = false
	}
}

// 对CONNECT包处理后的返回码
func (cp *ConnackPacket) ReturnCode() ConnackCode {
	return cp.returnCode
}

func (cp *ConnackPacket) SetReturnCode(ret ConnackCode) {
	cp.returnCode = ret
}

func (cp *ConnackPacket) Len() int {
	return cp.header.msglen() + cp.msglen()
}

func (cp *ConnackPacket) Decode(src []byte) (int, error) {
	var total int

	//获取固定报头,这里n=2
	n, err := cp.header.decode(src)
	total += n
	if err != nil {
		return total, err
	}

	//报头前两个字节是固定报头，第一个字节为mqtt控制报文类型，第二个字节是剩余长度(可变长度)
	//后两个字节是可变报头，第一个字节是连接确认标志Connack Acknowledge Flags，
	//高位的7bit必须为0，最后一个bit代表Session Present flag
	//第二个字节是连接返回码
	t := src[total]
	if t&254 != 0 {
		return 0, fmt.Errorf("connack/Decode.1: Invalid Connack Acknowledge Flags (%08b)", t)
	}

	cp.sessionPresent = (t&0x1 == 1)
	total++

	//获取返回码
	rc := src[total]

	if rc > 5 {
		return 0, fmt.Errorf("connack/Decode.2: Invalid CONNACK return code (%d)", rc)
	}

	cp.returnCode = ConnackCode(rc)
	total++

	return total, nil
}

func (cp *ConnackPacket) Encode() (int, []byte, error) {
	// 固定报头长度
	// hl := cp.header.msglen()
	// 报体长度:可变报头长度 ＋ 报体长度,Connack是2
	ml := cp.msglen()
	// if len(dst) < hl+ml {
	// 	return 0, fmt.Errorf("connack/Encode.2: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
	// }
	// 设置剩余长度
	if err := cp.SetRemainingLength(int32(ml)); err != nil {
		return 0, nil, err
	}

	dst := make([]byte, cp.Len())
	total := 0

	// 设置固定报头
	n, err := cp.header.encode(dst[total:])
	total += n
	if err != nil {
		return 0, nil, err
	}

	// 设置连接确认标志
	if cp.sessionPresent {
		dst[total] = 1
	}
	total++

	// 设置返回码
	if cp.returnCode > 5 {
		return total, nil, fmt.Errorf("connack/Encode.3: Invalid CONNACK return code (%d)", cp.returnCode)
	}
	dst[total] = cp.returnCode.Value()
	total++

	return total, dst, nil
}

func (cp *ConnackPacket) msglen() int {
	return 2
}
