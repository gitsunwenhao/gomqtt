package gate

import (
	"errors"

	proto "github.com/aiyun/gomqtt/mqtt/protocol"
	"github.com/uber-go/zap"
)

func processPacket(ci *connInfo, pt proto.Packet) error {
	switch t := pt.(type) {
	case *proto.DisconnectPacket:
		return errors.New("recv disconnect packet")

	case *proto.PublishPacket:

	case *proto.PubackPacket:

	case *proto.SubscribePacket:

	case *proto.UnsubscribePacket:

	case *proto.PingreqPacket:

	default:
		Logger.Warn("recv invalid packet type", zap.Object("invalid_type", t), zap.Int("cid", ci.id))
	}

	return nil
}
