package protocol

// ConnackCode是CONNACK报文的返回码，是对应CONNECT的报文结果
type ConnackCode byte

const (
	// 连接成功
	ConnectionAccepted ConnackCode = iota

	// 协议版本错误，3.1.1对应的是0x4版本
	ErrInvalidProtocolVersion

	// 客户端的标示符不合法
	ErrIdentifierRejected

	// 网络连接已建立，但是mqtt服务不可用
	ErrServerUnavailable

	// 用户名或者密码的格式有误
	ErrBadUsernameOrPassword

	// 连接鉴权失败
	ErrNotAuthorized
)

func (cc ConnackCode) Value() byte {
	return byte(cc)
}

func (cc ConnackCode) Desc() string {
	switch cc {
	case 0:
		return "Connection accepted"
	case 1:
		return "The Server does not support the level of the MQTT protocol requested by the Client"
	case 2:
		return "The Client identifier is correct UTF-8 but not allowed by the server"
	case 3:
		return "The Network Connection has been made but the MQTT service is unavailable"
	case 4:
		return "The data in the user name or password is malformed"
	case 5:
		return "The Client is not authorized to connect"
	}

	return ""
}

// Valid checks to see if the ConnackCode is valid. Currently valid codes are <= 5
func (cc ConnackCode) Valid() bool {
	return cc <= 5
}

// Error returns the corresonding error string for the ConnackCode
func (cc ConnackCode) Error() string {
	switch cc {
	case 0:
		return "Connection accepted"
	case 1:
		return "Connection Refused, unacceptable protocol version"
	case 2:
		return "Connection Refused, identifier rejected"
	case 3:
		return "Connection Refused, Server unavailable"
	case 4:
		return "Connection Refused, bad user name or password"
	case 5:
		return "Connection Refused, not authorized"
	}

	return "Unknown error"
}
