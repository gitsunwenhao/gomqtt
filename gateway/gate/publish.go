package gate

import (
	proto "github.com/aiyun/gomqtt/mqtt/protocol"
	"github.com/aiyun/gomqtt/mqtt/service"
)

func publish(ci *connInfo, p *proto.PublishPacket) error {
	// need give back the ack
	if p.QoS() == 1 {
		pb := proto.NewPubackPacket()
		pb.SetPacketID(p.PacketID())
		service.WritePacket(ci.c, pb)
	}
	return nil
}

func puback(ci *connInfo, p *proto.PubackPacket) error {
	return nil
}
