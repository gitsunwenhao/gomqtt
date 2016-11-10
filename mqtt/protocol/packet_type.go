package protocol

import (
	"fmt"
)

// 报文类型，MQTT报头的第一个字节，该字节最高的4位代表了控制报文，低4位为保留的标志位
type PacketType byte

const (
	//保留值，不合法的报文类型
	RESERVED PacketType = iota

	//客户端到服务器的连接请求
	CONNECT

	//服务器对客户端发起的连接请求给予回复
	CONNACK

	// PUBLISH: Client to Server, or Server to Client. Publish message.
	//发布消息，可以是客户端到服务端，也可以是服务端至客户端
	PUBLISH

	// PUBACK: Client to Server, or Server to Client. Publish acknowledgment for
	// QoS 1 messages.
	//发布消息的回执，QoS 1的消息才有
	PUBACK

	// PUBACK: Client to Server, or Server to Client. Publish received for QoS 2 messages.
	// Assured delivery part 1.
	//QoS 2消息特有，PUBLISH后，对端给予的第一次回复
	PUBREC

	// PUBREL: Client to Server, or Server to Client. Publish release for QoS 2 messages.
	// Assured delivery part 1.
	//QoS 2特有，收到对端的PUBREC后，发送PUBREL
	PUBREL

	//QOS2特有，Publish Complete
	PUBCOMP

	//客户端发起订阅服务器的请求
	SUBSCRIBE

	//服务器对SUBSCRIBE给予回执
	SUBACK

	//客户端发起取消订阅的请求
	UNSUBSCRIBE

	//服务器对UNSUBSCRIBE给予回执
	UNSUBACK

	// 客户端向服务器发起PING请求
	PINGREQ

	// 服务器对PING的回应
	PINGRESP

	// 客户端向服务器发起断开连接的请求
	DISCONNECT

	// 保留位
	RESERVED2
)

//For print
func (pt PacketType) String() string {
	return pt.Name()
}

//报文名称
func (pt PacketType) Name() string {
	switch pt {
	case RESERVED:
		return "RESERVED"
	case CONNECT:
		return "CONNECT"
	case CONNACK:
		return "CONNACK"
	case PUBLISH:
		return "PUBLISH"
	case PUBACK:
		return "PUBACK"
	case PUBREC:
		return "PUBREC"
	case PUBREL:
		return "PUBREL"
	case PUBCOMP:
		return "PUBCOMP"
	case SUBSCRIBE:
		return "SUBSCRIBE"
	case SUBACK:
		return "SUBACK"
	case UNSUBSCRIBE:
		return "UNSUBSCRIBE"
	case UNSUBACK:
		return "UNSUBACK"
	case PINGREQ:
		return "PINGREQ"
	case PINGRESP:
		return "PINGRESP"
	case DISCONNECT:
		return "DISCONNECT"
	case RESERVED2:
		return "RESERVED2"
	}

	return "UNKNOWN"
}

//报文描述
func (pt PacketType) Desc() string {
	switch pt {
	case RESERVED:
		return "Reserved"
	case CONNECT:
		return "Client request to connect to Server"
	case CONNACK:
		return "Connect acknowledgement"
	case PUBLISH:
		return "Publish message"
	case PUBACK:
		return "Publish acknowledgement"
	case PUBREC:
		return "Publish received (assured delivery part 1)"
	case PUBREL:
		return "Publish release (assured delivery part 2)"
	case PUBCOMP:
		return "Publish complete (assured delivery part 3)"
	case SUBSCRIBE:
		return "Client subscribe request"
	case SUBACK:
		return "Subscribe acknowledgement"
	case UNSUBSCRIBE:
		return "Unsubscribe request"
	case UNSUBACK:
		return "Unsubscribe acknowledgement"
	case PINGREQ:
		return "PING request"
	case PINGRESP:
		return "PING response"
	case DISCONNECT:
		return "Client is disconnecting"
	case RESERVED2:
		return "Reserved"
	default:
		return "UNKNOWN"
	}
}

//对应报文类型的默认标识位
func (pt PacketType) DefaultFlags() byte {
	switch pt {
	case PUBREL:
		return 2
	case SUBSCRIBE:
		return 2
	case UNSUBSCRIBE:
		return 2
	default:
		return 0
	}
}

//对应报文类型的默认标识位  兼容重发报文
func (pt PacketType) DefaultFlags10() byte {
	return 10
}

// New creates a new message based on the message type. It is a shortcut to call
// one of the New*Message functions. If an error is returned then the message type
// is invalid.
//根据报文类型创建新的报文，返回Packet接口
func (pt PacketType) New() (Packet, error) {
	switch pt {
	case CONNECT:
		return NewConnectPacket(), nil
	case CONNACK:
		return NewConnackPacket(), nil
	case PUBLISH:
		return NewPublishPacket(), nil
	case PUBACK:
		return NewPubackPacket(), nil
	case PUBREC:
		return NewPubrecPacket(), nil
	case PUBREL:
		return NewPubrelPacket(), nil
	case PUBCOMP:
		return NewPubcompPacket(), nil
	case SUBSCRIBE:
		return NewSubscribePacket(), nil
	case SUBACK:
		return NewSubackPacket(), nil
	case UNSUBSCRIBE:
		return NewUnsubscribePacket(), nil
	case UNSUBACK:
		return NewUnsubackPacket(), nil
	case PINGREQ:
		return NewPingreqPacket(), nil
	case PINGRESP:
		return NewPingrespPacket(), nil
	case DISCONNECT:
		return NewDisconnectPacket(), nil
	default:
		return nil, fmt.Errorf("msgtype/NewMessage: Invalid packet type %d", pt)
	}

}

//验证报文类型的合法性
func (pt PacketType) Valid() bool {
	return pt > RESERVED && pt < RESERVED2
}
