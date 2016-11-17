package gate

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"fmt"

	"os"

	"github.com/corego/tools"
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
		Addrs   []string
		Streams string
		Rooms   string
	}

	Mqtt struct {
		QosMax       byte
		MaxKeepalive uint16
	}

	StreamAddrs map[string]string
	RoomAddrs   map[string]string
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

	// 初始化Logger
	InitLogger(Conf.Common.LogPath, Conf.Common.LogLevel, Conf.Common.IsDebug)

	// stream hot update
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   Conf.Etcd.Addrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		Logger.Fatal("can't connect to etcd", zap.Error(err))
	}

	watchEtcd(cli)
	uploadEtcd(cli)

	fmt.Println(Conf)
}

// update the stream addrs
func watchEtcd(cli *clientv3.Client) {

	// update the stream addrs
	go func() {
		Conf.StreamAddrs = make(map[string]string)
		rch := cli.Watch(context.TODO(), Conf.Etcd.Streams, clientv3.WithPrefix())

		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == 0 { // PUT
					Conf.StreamAddrs[string(ev.Kv.Key)] = string(ev.Kv.Value)
				} else {
					delete(Conf.StreamAddrs, string(ev.Kv.Key))
				}
			}

			// Logger.Debug("get new stream addrs", zap.Object("addrs", Conf.StreamAddrs))
		}
	}()

	// update the room addrs
	go func() {
		Conf.RoomAddrs = make(map[string]string)
		rch := cli.Watch(context.TODO(), Conf.Etcd.Rooms, clientv3.WithPrefix())

		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == 0 { // PUT
					Conf.RoomAddrs[string(ev.Kv.Key)] = string(ev.Kv.Value)
				} else {
					delete(Conf.RoomAddrs, string(ev.Kv.Key))
				}
			}

			// Logger.Debug("get new room addrs", zap.Object("addrs", Conf.RoomAddrs))
		}
	}()
}

func uploadEtcd(cli *clientv3.Client) {
	key := Conf.Etcd.Rooms + "/" + getHost()
	ip := tools.LocalIP()
	Logger.Debug("local ip", zap.String("ip", ip))

	go func() {
		for {
			// upload self ip
			Grant, err := cli.Grant(context.TODO(), 120)
			if err != nil {
				Logger.Warn("etcd grant error", zap.Error(err))
			}

			_, err = cli.Put(context.TODO(), key, ip, clientv3.WithLease(Grant.ID))
			if err != nil {
				Logger.Warn("etcd put error", zap.Error(err))
			}

			time.Sleep(20 * time.Second)
		}
	}()

}

func getHost() string {
	host, err := os.Hostname()
	if err != nil {
		Logger.Fatal("get hostname error", zap.Error(err))
	}

	// in debug enviroment,we need to start several nodes in one machine,so pid is needed
	if Conf.Common.IsDebug {
		return fmt.Sprintf("%s-%d", host, os.Getpid())
	}

	return host
}
