package gate

import (
	"errors"
	"net"
	"time"

	"fmt"

	proto "github.com/aiyun/gomqtt/mqtt/protocol"
	"github.com/aiyun/gomqtt/mqtt/service"
	"github.com/corego/tools"
	"github.com/uber-go/zap"
)

func serve(c net.Conn) {
	// init a new connInfo
	ci := &connInfo{}

	//generate a uuid for this conn
	ci.id = 1
	ci.c = c
	Logger.Debug("a new connection has established", zap.Int("cid", ci.id), zap.String("ip", c.RemoteAddr().String()))

	defer func() {
		c.Close()
		delCI(ci.id)
	}()

	//----------------Connection init---------------------------------------------
	err := initConnection(ci)
	if err != nil {
		return
	}

	// save ci
	saveCI(ci)

	ci.stopped = make(chan struct{})
	go recvPacket(ci)

	// loop reading data
	for {
		select {
		case <-ci.stopped:
			Logger.Info("user's main thread is going to stop")
			goto STOP
		}
	}

STOP:
}

func initConnection(ci *connInfo) error {

	// the first packet is connect type,we need to restrain the read deadline
	ci.c.SetReadDeadline(time.Now().Add(10 * time.Second))

	reply := proto.NewConnackPacket()

	pt, buf, n, err := service.ReadPacket(ci.c)
	if err != nil {
		Logger.Warn("Read packet error", zap.Error(err), zap.String("buf", fmt.Sprintf("%v", buf)), zap.Int("bytes", n), zap.Int("cid", ci.id))

		if code, ok := err.(proto.ConnackCode); ok {
			reply.SetReturnCode(code)
			service.WritePacket(ci.c, reply)
		}
		return err
	}

	cp, ok := pt.(*proto.ConnectPacket)
	if !ok {
		Logger.Warn("this first packet is not connect type", zap.String("packet_type", fmt.Sprintf("%T", cp)), zap.Int("cid", ci.id))

		reply.SetReturnCode(proto.ErrIdentifierRejected)
		service.WritePacket(ci.c, reply)
		return errors.New("invalid packet")
	}

	ci.cp = cp

	Logger.Debug("user connected!", zap.String("user", tools.Bytes2String(ci.cp.Username())), zap.String("password", tools.Bytes2String(ci.cp.Password())), zap.Int("cid", ci.id),
		zap.Float64("keepalive", float64(cp.KeepAlive())))

	// validate the user
	ok = userValidate(ci.cp.Username(), ci.cp.Password())
	if !ok {
		Logger.Debug("user invalid", zap.Int("cid", ci.id))

		reply.SetReturnCode(proto.ErrIdentifierRejected)
		service.WritePacket(ci.c, reply)
		return errors.New("invalid user")
	}

	reply.SetReturnCode(proto.ConnectionAccepted)
	if err := service.WritePacket(ci.c, reply); err != nil {
		Logger.Info("write packet error", zap.Error(err), zap.Int("cid", ci.id))
		return err
	}

	// if keepalive == 0 ,we should specify a default keepalive
	if ci.cp.KeepAlive() == 0 {
		ci.cp.SetKeepAlive(Conf.Mqtt.MaxKeepalive)
	}

	return nil
}
