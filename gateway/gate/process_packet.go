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
		err = puback(ci, p)

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
