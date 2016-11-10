// Author - Sunface
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	maxLPString        uint16 = 65535
	maxRemainingLength int32  = 268435455 // 最大剩余长度
)

const (
	// QoS  0 : 最多发送一次，接收方也不需要发送回执
	QosAtMostOnce byte = iota

	// QoS 1 : 至少发送一次，这个等级会确保消息至少被接收方收到一次
	// QoS 1 的PUBLISH变长包头中有唯一的标示符，而且接收方会通过发送PUBACK包来告知服务器消息已经收到
	QosAtLeastOnce

	// QoS 2 : 刚好一次
	// 这是等级最高的服务，该服务下发送方仅发送一次包，接收方也只接收一次，不能丢包也不能重复发包
	// 该等级会给通信带来较大的负担，需要4次通信
	QosExactlyOnce

	// 若客户端订阅某个特定主题时，服务器发生了错误，那么会返回QosFailure
	QosFailure = 0x80
)

// 支持的版本号， 目前最常用的是MQTT
var SupportedVersions map[byte]string = map[byte]string{
	0x3: "MQIsdp",
	0x4: "MQTT",
}

// 所有的报文类型通用的接口
type Packet interface {
	// 返回Message的报文的字符串类型，例如"PUBLISH"、"SUBSCRIBE",该值是静态不可变的
	Name() string

	// 报文描述信息，例如CONNECT报文会返回"Client request to connect to Server."
	// 该信息也是静态不可变的
	Desc() string

	// 返回报文的MessageType字节
	Type() PacketType

	// 报文的唯一ID
	PacketID() uint16

	// 对message进行编码， 写入字节数组并返回
	Encode() (int, []byte, error)

	// 对字节数组进行解码，生成message
	Decode([]byte) (int, error)

	Len() int
}

// 验证PUBLISH时Topic的合法性，不能包含通配符，例如+和#
func ValidTopic(topic []byte) bool {
	return len(topic) > 0 && bytes.IndexByte(topic, '#') == -1 && bytes.IndexByte(topic, '+') == -1
}

// 验证QoS的合法性
func ValidQos(qos byte) bool {
	return qos == QosAtLeastOnce || qos == QosAtMostOnce || qos == QosExactlyOnce
}

// 在连接时，验证版本号
func ValidVersion(v byte) bool {
	_, ok := SupportedVersions[v]
	return ok
}

// 判断给定的error是否是Connack的error
func ValidConnackError(err error) bool {
	return err == ErrInvalidProtocolVersion ||
		err == ErrIdentifierRejected ||
		err == ErrServerUnavailable ||
		err == ErrBadUsernameOrPassword ||
		err == ErrNotAuthorized
}

// 读取变长字段
// 首先是2字节的字段长度表示，然后是对应长度的字段主体部分
func readLPBytes(buf []byte) ([]byte, int, error) {
	if len(buf) < 2 {
		return nil, 0, fmt.Errorf("readLPBytes: Insufficient buffer size. Expecting %d, got %d.", 2, len(buf))
	}

	n, total := 0, 0

	n = int(binary.BigEndian.Uint16(buf))
	total += 2

	if len(buf) < n {
		return nil, total, fmt.Errorf("readLPBytes: Insufficient buffer size. Expecting %d, got %d.", n, len(buf))
	}

	total += n

	return buf[2:total], total, nil
}

// 写入变长字段
// 首先是2字节的字段长度表示，然后是对应长度的字段主体部分
func writeLPBytes(buf []byte, b []byte) (int, error) {
	total, n := 0, len(b)

	if n > int(maxLPString) {
		return 0, fmt.Errorf("writeLPBytes: Length (%d) greater than %d bytes.", n, maxLPString)
	}

	if len(buf) < 2+n {
		return 0, fmt.Errorf("writeLPBytes: Insufficient buffer size. Expecting %d, got %d.", 2+n, len(buf))
	}

	binary.BigEndian.PutUint16(buf, uint16(n))
	total += 2

	copy(buf[total:], b)
	total += n

	return total, nil
}
