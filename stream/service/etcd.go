package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/uber-go/zap"
)

// SerStatus server report key addr
type SerStatus struct {
	Key  string //上报key
	Addr string //上报地址
}

// Etcd struct
type Etcd struct {
	Client    *clientv3.Client
	ReportKey string
	Grant     *clientv3.LeaseGrantResponse
}

// Init init etcd
func (e *Etcd) Init() {
	reportKey, err := GetRegisterKey()
	if err != nil {
		Logger.Panic("Init", zap.Error(err))
	}
	keylen := len(Conf.EtcdC.Reportdir)
	if keylen > 0 && Conf.EtcdC.Reportdir[keylen-1] != '/' {
		e.ReportKey = Conf.EtcdC.Reportdir + "/" + reportKey
	} else {
		e.ReportKey = Conf.EtcdC.Reportdir + reportKey
	}

	Logger.Info("Init", zap.String("@Key", e.ReportKey), zap.String("addr", Conf.GrpcC.Addr))

	cfg := clientv3.Config{
		Endpoints:   Conf.EtcdC.Addrs,
		DialTimeout: time.Duration(Conf.EtcdC.Dltimeout) * time.Second,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		Logger.Panic("Etcd", zap.Error(err))
	}

	e.Client = client

	return
}

func (e *Etcd) Start() {
	go e.WathchWork()
	go e.RegisterWork()
}

func (e *Etcd) Close() error {
	return e.Client.Close()
}

func (e *Etcd) WathchWork() {
	for {
		rch := e.Client.Watch(context.Background(), Conf.EtcdC.Reportdir, clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				// PUT
				if ev.Type == 0 {
					log.Println("Put", string(ev.Kv.Value))
					exist := gStream.streamAddrs.Add(string(ev.Kv.Key), string(ev.Kv.Value))
					if !exist {
						// new stream, 重新计算客户端落在哪台一台stream
						Logger.Info("WathchWork", zap.String("key", string(ev.Kv.Key)), zap.String("Value", string(ev.Kv.Value)))
					}
				} else {
					// Delete
					log.Println("Delete", string(ev.Kv.Key))
					gStream.streamAddrs.Del(string(ev.Kv.Key))
				}
			}
			Logger.Debug("get new stream addrs", zap.Object("addrs", gStream.streamAddrs))
		}
	}
}

// RegisterWork register stream addr
func (e *Etcd) RegisterWork() {
	for {
		time.Sleep(time.Duration(Conf.EtcdC.ReportTime) * time.Second)
		e.Put(e.ReportKey, Conf.GrpcC.Addr, Conf.EtcdC.TTL)
	}
}

func (e *Etcd) Put(key, value string, ttl int64) error {
	Grant, err := e.Client.Grant(context.TODO(), ttl)
	if err != nil {
		Logger.Panic("Etcd", zap.Error(err), zap.Int64("@ReportTime", ttl))
	}
	_, err = e.Client.Put(context.TODO(), key, value, clientv3.WithLease(Grant.ID))
	if err != nil {
		Logger.Error("Put", zap.String("@key", key), zap.String("@value", value), zap.Error(err))
		return err
	}
	return nil
}

// NewEtcd new etcd
func NewEtcd() *Etcd {
	etcd := &Etcd{}
	return etcd
}

// GetRegisterKey get etcd Register key
func GetRegisterKey() (string, error) {
	host, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%d", host, os.Getpid()), nil
}
