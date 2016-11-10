package protocol

import (
	"encoding/binary"
	"fmt"
	"regexp"
)

//客户端ID检测
var clientIdRegexp *regexp.Regexp

func init() {
	clientIdRegexp = regexp.MustCompile("^[0-9a-zA-Z]*$")
}

// 当服务器和客户端建立连接后，客户端发送的第一个包必须是CONNECT包
// CONNECT包只能发送一次，如果服务器收到同一个客户端多次CONNECT包，就要关闭该连接
// connectFlags连接标识位，一个字节共8bit
// 7: username flag
// 6: password flag
// 5: will retain
// 4-3: will QoS
// 2: will flag
// 1: clean session
// 0: reserved
type ConnectPacket struct {
	header

	connectFlags byte

	version byte

	keepAlive uint16

	protoName,
	clientId,
	willTopic,
	willMessage,
	username,
	password []byte
}

// NewConnectPacket创建CONNECT包.
func NewConnectPacket() *ConnectPacket {
	cp := &ConnectPacket{}
	cp.SetType(CONNECT)

	return cp
}

//用于打印CONNECT时
func (cp ConnectPacket) String() string {
	return fmt.Sprintf("%s, Connect Flags=%08b, Version=%d, KeepAlive=%d, Client ID=%q, Will Topic=%q, Will Message=%q, Username=%q, Password=%q",
		cp.header,
		cp.connectFlags,
		cp.Version(),
		cp.KeepAlive(),
		cp.ClientId(),
		cp.WillTopic(),
		cp.WillMessage(),
		cp.Username(),
		cp.Password(),
	)
}

// 返回客户端连接服务器时选择的版本号，mqtt3.1.1的协议版本号是4
func (cp *ConnectPacket) Version() byte {
	return cp.version
}

// 设置版本号
func (cp *ConnectPacket) SetVersion(ver byte) error {
	if _, ok := SupportedVersions[ver]; !ok {
		return fmt.Errorf("SetVersion: Invalid version number %d", ver)
	}

	cp.version = ver

	return nil
}

// CleanSession返回代表session处理的bit位
// 客户端和服务器都可以存储seesion状态,以便在不同的连接下继续进行网络传输，
// 这个bit位用于控制session的生命周期
func (cp *ConnectPacket) CleanSession() bool {
	return ((cp.connectFlags >> 1) & 0x1) == 1
}

func (cp *ConnectPacket) SetCleanSession(v bool) {
	if v {
		cp.connectFlags |= 0x2 // 00000010
	} else {
		cp.connectFlags &= 253 // 11111101
	}

}

// WillFlag返回遗愿bit位的设置情况,如果该位被设置为1，那么一旦连接建立，遗愿消息必须存储在服务器且和
// 该连接紧密相关
func (cp *ConnectPacket) WillFlag() bool {
	return ((cp.connectFlags >> 2) & 0x1) == 1
}

func (cp *ConnectPacket) SetWillFlag(v bool) {
	if v {
		cp.connectFlags |= 0x4 // 00000100
	} else {
		cp.connectFlags &= 251 // 11111011
	}
}

// WillQos代表遗愿消息发布时的Qos等级，用2bit表示
func (cp *ConnectPacket) WillQos() byte {
	return (cp.connectFlags >> 3) & 0x3
}

func (cp *ConnectPacket) SetWillQos(qos byte) error {
	if qos != QosAtMostOnce && qos != QosAtLeastOnce && qos != QosExactlyOnce {
		return fmt.Errorf("SetWillQos: Invalid QoS level %d", qos)
	}

	cp.connectFlags = (cp.connectFlags & 231) | (qos << 3) // 231 = 11100111

	return nil
}

// WillRetain位控制遗愿消息发布时是否先做存储
func (cp *ConnectPacket) WillRetain() bool {
	return ((cp.connectFlags >> 5) & 0x1) == 1
}

func (cp *ConnectPacket) SetWillRetain(v bool) {
	if v {
		cp.connectFlags |= 32 // 00100000
	} else {
		cp.connectFlags &= 223 // 11011111
	}

}

// UsernameFlag位指示了，在payload中是否存在用户名
func (cp *ConnectPacket) UsernameFlag() bool {
	return ((cp.connectFlags >> 7) & 0x1) == 1
}

func (cp *ConnectPacket) SetUsernameFlag(v bool) {
	if v {
		cp.connectFlags |= 128 // 10000000
	} else {
		cp.connectFlags &= 127 // 01111111
	}

}

// payload中是否存在密码值
func (cp *ConnectPacket) PasswordFlag() bool {
	return ((cp.connectFlags >> 6) & 0x1) == 1
}

func (cp *ConnectPacket) SetPasswordFlag(v bool) {
	if v {
		cp.connectFlags |= 64 // 01000000
	} else {
		cp.connectFlags &= 191 // 10111111
	}

}

// 客户端传递控制报文间的最大时间间隔(相邻两次 )
func (cp *ConnectPacket) KeepAlive() uint16 {
	return cp.keepAlive
}

func (cp *ConnectPacket) SetKeepAlive(t uint16) {
	cp.keepAlive = t
}

// 每个客户端的ID都要是唯一的，这个ID应该用作sessionID
func (cp *ConnectPacket) ClientId() []byte {
	return cp.clientId
}

// SetClientId sets an ID that identifies the Client to the Server.
func (cp *ConnectPacket) SetClientId(id []byte) error {
	if len(id) > 0 && !cp.validClientId(id) {
		return ErrIdentifierRejected
	}

	cp.clientId = id

	return nil
}

// 遗愿消息要发布到的topic，如果该标志位为1，那么payload中必须要有will topic
func (cp *ConnectPacket) WillTopic() []byte {
	return cp.willTopic
}

func (cp *ConnectPacket) SetWillTopic(t []byte) {
	cp.willTopic = t

	if len(t) > 0 {
		cp.SetWillFlag(true)
	} else if len(cp.willMessage) == 0 {
		cp.SetWillFlag(false)
	}

}

// 遗愿消息体
func (cp *ConnectPacket) WillMessage() []byte {
	return cp.willMessage
}

func (cp *ConnectPacket) SetWillMessage(msg []byte) {
	cp.willMessage = msg

	if len(msg) > 0 {
		cp.SetWillFlag(true)
	} else if len(cp.willTopic) == 0 {
		cp.SetWillFlag(false)
	}

}

// 用户名
func (cp *ConnectPacket) Username() []byte {
	return cp.username
}

func (cp *ConnectPacket) SetUsername(u []byte) {
	cp.username = u

	if len(u) > 0 {
		cp.SetUsernameFlag(true)
	} else {
		cp.SetUsernameFlag(false)
	}

}

// 密码
func (cp *ConnectPacket) Password() []byte {
	return cp.password
}

func (cp *ConnectPacket) SetPassword(pw []byte) {
	cp.password = pw

	if len(pw) > 0 {
		cp.SetPasswordFlag(true)
	} else {
		cp.SetPasswordFlag(false)
	}

}

func (cp *ConnectPacket) Len() int {
	return cp.header.msglen() + cp.msglen()
}

// 对于CONNECT包，下面两个方法返回的error可能会是ConnackReturnCode,所以要检查error值。如果返回的是
// 一般性的错误，那么就要认为这个包是不合法的
// 调用者可以通过ValidConnackError(err)来检查返回的错误是否是Connack错误，如果是这个错误，那么
// 调用者应该向客户端返回相应的CONNACK消息
func (cp *ConnectPacket) Decode(src []byte) (int, error) {
	total := 0

	n, err := cp.header.decode(src[total:])
	if err != nil {
		return total + n, err
	}
	total += n

	if n, err = cp.decodeMessage(src[total:]); err != nil {
		return total + n, err
	}
	total += n

	return total, nil
}

func (cp *ConnectPacket) Encode() (int, []byte, error) {
	if cp.Type() != CONNECT {
		return 0, nil, fmt.Errorf("connect/Encode: Invalid message type. Expecting %d, got %d", CONNECT, cp.Type())
	}

	_, ok := SupportedVersions[cp.version]
	if !ok {
		return 0, nil, ErrInvalidProtocolVersion
	}

	// hl := cp.header.msglen()
	ml := cp.msglen()

	// if len(dst) < hl+ml {
	// 	return 0, nil, fmt.Errorf("connect/Encode: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
	// }
	dst := make([]byte, cp.Len())
	if err := cp.SetRemainingLength(int32(ml)); err != nil {
		return 0, nil, err
	}

	total := 0

	n, err := cp.header.encode(dst[total:])
	total += n
	if err != nil {
		return total, nil, err
	}

	n, err = cp.encodeMessage(dst[total:])
	total += n
	if err != nil {
		return total, nil, err
	}

	return total, dst, nil
}

func (cp *ConnectPacket) encodeMessage(dst []byte) (int, error) {
	total := 0

	n, err := writeLPBytes(dst[total:], []byte(SupportedVersions[cp.version]))
	total += n
	if err != nil {
		return total, err
	}

	dst[total] = cp.version
	total += 1

	dst[total] = cp.connectFlags
	total += 1

	binary.BigEndian.PutUint16(dst[total:], cp.keepAlive)
	total += 2

	n, err = writeLPBytes(dst[total:], cp.clientId)
	total += n
	if err != nil {
		return total, err
	}

	if cp.WillFlag() {
		n, err = writeLPBytes(dst[total:], cp.willTopic)
		total += n
		if err != nil {
			return total, err
		}

		n, err = writeLPBytes(dst[total:], cp.willMessage)
		total += n
		if err != nil {
			return total, err
		}
	}

	if cp.UsernameFlag() && len(cp.username) > 0 {
		n, err = writeLPBytes(dst[total:], cp.username)
		total += n
		if err != nil {
			return total, err
		}
	}

	if cp.PasswordFlag() && len(cp.password) > 0 {
		n, err = writeLPBytes(dst[total:], cp.password)
		total += n
		if err != nil {
			return total, err
		}
	}

	return total, nil
}

func (cp *ConnectPacket) decodeMessage(src []byte) (int, error) {
	var err error
	n, total := 0, 0

	cp.protoName, n, err = readLPBytes(src[total:])
	total += n
	if err != nil {
		return total, err
	}

	cp.version = src[total]
	total++

	if verstr, ok := SupportedVersions[cp.version]; !ok {
		return total, ErrInvalidProtocolVersion
	} else if verstr != string(cp.protoName) {
		return total, ErrInvalidProtocolVersion
	}

	cp.connectFlags = src[total]
	total++

	if cp.connectFlags&0x1 != 0 {
		return total, fmt.Errorf("connect/decodeMessage: Connect Flags reserved bit 0 is not 0")
	}

	if cp.WillQos() > QosExactlyOnce {
		return total, fmt.Errorf("connect/decodeMessage: Invalid QoS level (%d) for %s message", cp.WillQos(), cp.Name())
	}

	if !cp.WillFlag() && (cp.WillRetain() || cp.WillQos() != QosAtMostOnce) {
		return total, fmt.Errorf("connect/decodeMessage: Protocol violation: If the Will Flag (%t) is set to 0 the Will QoS (%d) and Will Retain (%t) fields MUST be set to zero", cp.WillFlag(), cp.WillQos(), cp.WillRetain())
	}

	if cp.UsernameFlag() && !cp.PasswordFlag() {
		return total, fmt.Errorf("connect/decodeMessage: Username flag is set but Password flag is not set")
	}

	if len(src[total:]) < 2 {
		return 0, fmt.Errorf("connect/decodeMessage: Insufficient buffer size. Expecting %d, got %d.", 2, len(src[total:]))
	}

	cp.keepAlive = binary.BigEndian.Uint16(src[total:])
	total += 2

	cp.clientId, n, err = readLPBytes(src[total:])
	total += n
	if err != nil {
		return total, err
	}

	// 如果客户端不传ClientId,那么CleanSession必须设置为1
	if len(cp.clientId) == 0 && !cp.CleanSession() {
		return total, ErrIdentifierRejected
	}

	// ClientId只支持0-9,a-z,A-Z这些字符
	if len(cp.clientId) > 0 && !cp.validClientId(cp.clientId) {
		return total, ErrIdentifierRejected
	}

	if cp.WillFlag() {
		cp.willTopic, n, err = readLPBytes(src[total:])
		total += n
		if err != nil {
			return total, err
		}

		cp.willMessage, n, err = readLPBytes(src[total:])
		total += n
		if err != nil {
			return total, err
		}
	}

	// 在3.1协议中，可以允许用户位设置了，但是用户为空
	if cp.UsernameFlag() && len(src[total:]) > 0 {
		cp.username, n, err = readLPBytes(src[total:])
		total += n
		if err != nil {
			return total, err
		}
	}

	// 在3.1协议中，可以允许密码位设置了，但是密码为空
	if cp.PasswordFlag() && len(src[total:]) > 0 {
		cp.password, n, err = readLPBytes(src[total:])
		total += n
		if err != nil {
			return total, err
		}
	}

	return total, nil
}

func (cp *ConnectPacket) msglen() int {
	total := 0

	ver, ok := SupportedVersions[cp.version]
	if !ok {
		return total
	}

	// 2字节的协议长度前缀
	// n字节协议名
	// 1字节协议版本
	// 1字节连接标识
	// 2字节keepalive时间
	total += 2 + len(ver) + 1 + 1 + 2

	// Add the clientID length, 2 is the length prefix
	// 添加clientID长度，2个字节的长度前缀
	total += 2 + len(cp.clientId)

	// 添加遗愿topic和遗愿消息，以及它们的长度前缀
	if cp.WillFlag() {
		total += 2 + len(cp.willTopic) + 2 + len(cp.willMessage)
	}

	// 添加用户名长度
	if cp.UsernameFlag() && len(cp.username) > 0 {
		total += 2 + len(cp.username)
	}

	if cp.PasswordFlag() && len(cp.password) > 0 {
		total += 2 + len(cp.password)
	}

	return total
}

func (cp *ConnectPacket) validClientId(cid []byte) bool {
	if cp.Version() == 0x3 {
		return true
	}

	return clientIdRegexp.Match(cid)
}
