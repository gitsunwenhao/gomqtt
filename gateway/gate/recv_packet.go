package gate

import (
	"fmt"
	"time"

	"github.com/aiyun/gomqtt/mqtt/service"
	"github.com/uber-go/zap"
)

func recvPacket(ci *connInfo) {
	defer func() {
		close(ci.stopped)
	}()

	wait := time.Duration(ci.cp.KeepAlive())*time.Second - 10*time.Second

	for {
		// We need to considering about the network delay,so here allows 10 seconds delay.
		ci.c.SetReadDeadline(time.Now().Add(wait))

		pt, buf, n, err := service.ReadPacket(ci.c)
		if err != nil {
			Logger.Warn("Read packet error", zap.Error(err), zap.String("buf", fmt.Sprintf("%v", buf)), zap.Int("bytes", n), zap.Int("cid", ci.id))
			break
		}

		err = processPacket(ci, pt)
		if err != nil {
			break
		}

		ci.inCount++
	}
}
