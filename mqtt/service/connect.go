package service

import (
	"net"

	"github.com/aiyun/gomqtt/mqtt/protocol"
)

func connectPacket(conn net.Conn) (*protocol.ConnectPacket, error) {
	buf, err := Read(conn)
	if err != nil {
		return nil, err
	}

	cp := protocol.NewConnectPacket()

	_, err = cp.Decode(buf)
	return cp, err
}
