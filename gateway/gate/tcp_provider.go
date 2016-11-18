package gate

import (
	"io"
	"net"
	"strings"
	"time"

	"fmt"

	"crypto/tls"

	"github.com/uber-go/zap"
)

/* Websocket Provider */

type TcpProvider struct {
}

func (tp *TcpProvider) Start() {
	var ln net.Listener
	var err error

	if !Conf.Provider.EnableTls { //start tcp
		ln, err = net.Listen("tcp", Conf.Provider.TcpAddr)
		if err != nil {
			Logger.Fatal("tcp Listen", zap.Error(err))
		}

		Logger.Debug("tcp provider startted", zap.String("addr", Conf.Provider.TcpAddr))

	} else { // start tls
		cert, err := tls.LoadX509KeyPair(Conf.Provider.TlsCert, Conf.Provider.TlsKey)
		if err != nil {
			Logger.Fatal("tls load cert", zap.Error(err))
			return
		}

		config := &tls.Config{Certificates: []tls.Certificate{cert}}
		ln, err = tls.Listen("tcp", Conf.Provider.TcpAddr, config)
		if err != nil {
			Logger.Fatal("tls listen", zap.Error(err))
			return
		}

		Logger.Debug("tls provider startted", zap.String("addr", Conf.Provider.TcpAddr))
	}

	// start accepting
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
