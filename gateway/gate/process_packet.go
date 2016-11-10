package gate

import (
	"errors"

	"fmt"

	proto "github.com/aiyun/gomqtt/mqtt/protocol"
	"github.com/aiyun/gomqtt/mqtt/service"
	"github.com/uber-go/zap"
)

func processPacket(ci *connInfo, pt proto.Packet) error {
	var err error
	switch p := pt.(type) {
	case *proto.DisconnectPacket: // recv Disconnect
		Logger.Info("Disconnect")
		err = errors.New("recv disconnect packet")

	case *proto.PublishPacket: // recv publish
		err = publish(ci, p)

	case *proto.PubackPacket:

	case *proto.SubscribePacket:
		err = subscribe(ci, p)

	case *proto.UnsubscribePacket:
		err = unsubscribe(ci, p)

	case *proto.PingreqPacket:
		Logger.Info("recv ping req")
		pingReq(ci)
	default:
		Logger.Warn("recv invalid packet type", zap.String("invalid_type", fmt.Sprintf("%T", pt)), zap.Int("cid", ci.id))
	}

	return err
}

func pingReq(ci *connInfo) {
	pb := proto.NewPingrespPacket()
	service.WritePacket(ci.c, pb)
}

func publish(ci *connInfo, p *proto.PublishPacket) error {
	// need give back the ack
	if p.QoS() == 1 {
		pb := proto.NewPubackPacket()
		pb.SetPacketID(p.PacketID())
		service.WritePacket(ci.c, pb)
	}
	return nil
}

func subscribe(ci *connInfo, p *proto.SubscribePacket) error {
	pb := proto.NewSubackPacket()
	pb.SetPacketID(p.PacketID())

	// return the final qos level
	for i := 0; i < len(p.Qos()); i++ {
		pb.AddReturnCodes([]byte{proto.QosAtLeastOnce})
	}

	service.WritePacket(ci.c, pb)

	return nil
}

func unsubscribe(ci *connInfo, p *proto.UnsubscribePacket) error {
	pb := proto.NewUnsubackPacket()
	pb.SetPacketID(p.PacketID())

	service.WritePacket(ci.c, pb)
	return nil
}
