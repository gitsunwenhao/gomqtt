package gate

import (
	"fmt"
	"net"
	"time"

	proto "github.com/aiyun/gomqtt/mqtt/protocol"
	"github.com/aiyun/gomqtt/mqtt/service"
	"github.com/uber-go/zap"
)

func recvPacket(ci *connInfo) {
	defer func() {
		ci.stopped <- true
	}()

	for {
		//check
		now := time.Now()

		// according to mqtt protocol, interval should be greater than 1.5 * keepalive
		if now.Sub(ci.lastPacketTime).Seconds() > 1.5*float64(ci.cp.KeepAlive()) {
			Logger.Info("not activity,stopped")
			break
		}

		pt, needContinue := read(ci)
		if ci.isStopped || !needContinue {
			break
		}

		if pt == nil {
			continue
		}

		err := processPacket(ci, pt)
		if err != nil {
			break
		}

		// update the last packet time
		ci.lastPacketTime = time.Now()
		ci.inCount++
	}
}

func read(ci *connInfo) (proto.Packet, bool) {
	needContinue := true

	ci.c.SetReadDeadline(time.Now().Add(5 * time.Second))

	pt, buf, n, err := service.ReadPacket(ci.c)
	if err != nil {
		_, ok := err.(net.Error)
		if !ok {
			Logger.Warn("Read packet error", zap.Error(err), zap.String("buf", fmt.Sprintf("%v", buf)), zap.Int("bytes", n), zap.Int("cid", ci.id))
			needContinue = false
		}
	}

	return pt, needContinue
}
