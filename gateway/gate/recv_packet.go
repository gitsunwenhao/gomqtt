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
		ci.stopped = true
	}()

	for {
		if ci.stopped {
			goto STOP
		}

		pt, needContinue := read(ci)
		if !needContinue {
			goto STOP
		}

		err := processPacket(ci, pt)
		if err != nil {

		}

		ci.inCount++
	}

STOP:
}

func read(ci *connInfo) (proto.Packet, bool) {
	needContinue := true

	ci.c.SetReadDeadline(time.Now().Add(5 * time.Second))

	pt, buf, n, err := service.ReadPacket(ci.c)
	if err != nil {
		_, ok := err.(net.Error)
		if !ok {
			Logger.Warn("Read packet error", zap.Error(err), zap.Object("buf", buf), zap.Int("bytes", n), zap.Int("cid", ci.id))
			needContinue = false
		}
	}

	return pt, needContinue
}
