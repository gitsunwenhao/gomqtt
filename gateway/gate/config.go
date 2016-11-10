package gate

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"fmt"

	"github.com/coreos/etcd/clientv3"
	"github.com/naoina/toml"
	"github.com/uber-go/zap"
)

type Config struct {
	Common struct {
		Version  string
		IsDebug  bool
		LogLevel string
		LogPath  string
	}

	Provider struct {
		Invoked   []string
		TcpAddr   string
		EnableTls bool
		TlsCert   string
		TlsKey    string
	}

	Etcd struct {
		Addrs []string
		Key   string
	}

	StreamAddrs map[string]string
}

var Conf = &Config{}

func loadConfig(staticConf bool) {
	var contents []byte
	var err error

	if staticConf {
		//静态配置
		contents, err = ioutil.ReadFile("configs/gateway.toml")
	} else {
		contents, err = ioutil.ReadFile("/etc/gomqtt/gateway.toml")
	}

	if err != nil {
		log.Fatal("load config error", zap.Error(err))
	}

	tbl, err := toml.Parse(contents)
	if err != nil {
		log.Fatal("parse config error", zap.Error(err))
	}

	toml.UnmarshalTable(tbl, Conf)

	go updateStreamAddr()

	fmt.Println(Conf)

	// 初始化Logger
	InitLogger(Conf.Common.LogPath, Conf.Common.LogLevel, Conf.Common.IsDebug)
}

// update the stream addrs
func updateStreamAddr() {
	// stream hot update
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   Conf.Etcd.Addrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		Logger.Fatal("can't connect to etcd", zap.Error(err))
	}
	defer cli.Close()

	Conf.StreamAddrs = make(map[string]string)
	rch := cli.Watch(context.TODO(), Conf.Etcd.Key, clientv3.WithPrefix())

	for wresp := range rch {
		for _, ev := range wresp.Events {
			if ev.Type == 0 { // PUT
				Conf.StreamAddrs[string(ev.Kv.Key)] = string(ev.Kv.Value)
			} else {
				delete(Conf.StreamAddrs, string(ev.Kv.Key))
			}
		}

		Logger.Debug("get new stream addrs", zap.Object("addrs", Conf.StreamAddrs))
	}
}
