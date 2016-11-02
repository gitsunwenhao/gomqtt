package gate

import (
	"net"
	"time"

	proto "github.com/aiyun/gomqtt/mqtt/protocol"
	"github.com/aiyun/gomqtt/mqtt/service"
	"github.com/uber-go/zap"
)

func recvPacket(ci *connInfo) {
	defer func() {
		ci.stopCh <- true
	}()

	for {
		if ci.stopped {
			goto STOP
		}

		ci.c.SetReadDeadline(time.Now().Add(5 * time.Second))

		pt, buf, n, err := service.ReadPacket(ci.c)
		if err != nil {
			e, ok := err.(net.Error)
			if !ok {
				Logger.Warn("Read packet error", zap.Error(err), zap.Object("buf", buf), zap.Int("bytes", n), zap.Int("cid", ci.id))
				goto STOP
			}

			if e.Timeout() {
				continue
			}

			if e.Temporary() {
				time.Sleep(1 * time.Second)
				continue
			}
		}

		switch t := pt.(type) {
		case *proto.DisconnectPacket:
			goto STOP

		case *proto.PublishPacket:

		case *proto.PubackPacket:

		case *proto.SubscribePacket:

		case *proto.UnsubscribePacket:

		case *proto.PingreqPacket:

		default:
			Logger.Warn("recv invalid packet type", zap.Object("invalid_type", t), zap.Int("cid", ci.id))
		}
	}

STOP:
}
