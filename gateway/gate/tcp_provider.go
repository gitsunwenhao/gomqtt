package gate

import (
	"io"
	"net"
	"strings"
	"time"

	"fmt"

	"github.com/uber-go/zap"
)

/* Websocket Provider */

type TcpProvider struct {
}

func (tp *TcpProvider) Start() {
	if !Conf.Provider.EnableTls {
		ln, err := net.Listen("tcp", Conf.Provider.TcpAddr)
		if err != nil {
			Logger.Fatal("Listen", zap.Error(err))
		}

		for {
			c, err := accept(ln)
			if err != nil {
				Logger.Warn("accept tcp connection error", zap.Error(err))
				if c != nil {
					c.Close()
				}
			}
			go serve(c)
		}
	}
	Logger.Debug("tcp provider startted")
}

func (tp *TcpProvider) Close() error {
	return nil
}

func accept(ln net.Listener) (net.Conn, error) {
	c, err := ln.Accept()
	if err != nil {
		if c != nil {
			return c, fmt.Errorf("net.Listener returned non-nil conn and non-nil error : %v", err)
		}

		if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
			time.Sleep(time.Second)
			return c, fmt.Errorf("Temporary error when accepting new connections: %s\n", netErr)
		}

		if err != io.EOF && !strings.Contains(err.Error(), "use of closed network connection") {
			return c, fmt.Errorf("Permanent error when accepting new connections: %s\n", err)
		}

		return c, io.EOF
	}

	if c == nil {
		Logger.Fatal("BUG: net.Listener returned (nil, nil)")
	}

	return c, nil
}
