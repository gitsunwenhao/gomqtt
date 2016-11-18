package gate

import (
	proto "github.com/aiyun/gomqtt/mqtt/protocol"
	"github.com/aiyun/gomqtt/mqtt/service"
)

func subscribe(ci *connInfo, p *proto.SubscribePacket) error {
	var rets []byte

	for i, t := range p.Topics() {
		qos, err := subToStream(t, p.Qos()[i])
		if err != nil {

		}

		rets = append(rets, qos)
	}

	// give back the suback
	pb := proto.NewSubackPacket()
	pb.SetPacketID(p.PacketID())

	// return the final qos level
	pb.AddReturnCodes(rets)
	service.WritePacket(ci.c, pb)

	return nil
}

func unsubscribe(ci *connInfo, p *proto.UnsubscribePacket) error {
	pb := proto.NewUnsubackPacket()
	pb.SetPacketID(p.PacketID())

	service.WritePacket(ci.c, pb)
	return nil
}

func subToStream(t []byte, qos byte) (byte, error) {
	return 1, nil
}
