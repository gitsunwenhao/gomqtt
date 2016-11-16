package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/uber-go/zap"
)

// UpdateAddr struct
type UpdateAddr struct {
	Client    *clientv3.Client
	Grant     *clientv3.LeaseGrantResponse
	ReportKey string
}

// Init init UpdateAddr
func (upa *UpdateAddr) Init() {
	reportKey, err := GetRegisterKey()
	if err != nil {
		Logger.Panic("Init", zap.Error(err))
	}
	keylen := len(Conf.EtcdC.Reportdir)
	if keylen > 0 && Conf.EtcdC.Reportdir[keylen-1] != '/' {
		upa.ReportKey = Conf.EtcdC.Reportdir + "/" + reportKey
	} else {
		upa.ReportKey = Conf.EtcdC.Reportdir + reportKey
	}

	Logger.Info("Init", zap.String("@Key", upa.ReportKey), zap.String("addr", Conf.GrpcC.Addr))

	cfg := clientv3.Config{
		Endpoints:   Conf.EtcdC.Addrs,
		DialTimeout: time.Duration(Conf.EtcdC.Dltimeout) * time.Second,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		Logger.Panic("Etcd", zap.Error(err))
	}

	upa.Client = client

	return
}

func (upa *UpdateAddr) Start() {
	go upa.WathchWork()
	go upa.RegisterWork()
}

func (upa *UpdateAddr) Close() error {
	return upa.Client.Close()
}

func (upa *UpdateAddr) WathchWork() {
	for {
		rch := upa.Client.Watch(context.Background(), Conf.EtcdC.Reportdir, clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				comput := false
				// PUT
				if ev.Type == 0 {
					exist := gStream.cache.Sas.Add(string(ev.Kv.Key), string(ev.Kv.Value))
					if !exist {
						comput = true
						// hast add grpc addr
						gStream.hash.Add(string(ev.Kv.Value))
						// new stream, 重新计算客户端落在哪台一台stream
						Logger.Info("Watch Insert", zap.String("key", string(ev.Kv.Key)), zap.String("value", string(ev.Kv.Value)))
					}
				} else {
					// Delete
					if keyip, ok := gStream.cache.Sas.Get(string(ev.Kv.Key)); ok {
						gStream.hash.Remove(keyip)
						gStream.cache.Sas.Del(string(ev.Kv.Key))
						Logger.Info("Watch Delete", zap.String("key", string(ev.Kv.Key)), zap.String("value", string(ev.Kv.Value)))
						comput = true
					}
				}
				// 重新计算在线用户是否需要保存在本台stream上
				if comput {
					Logger.Debug("需要重新计算用户是否落在本机上")
				}
			}
			Logger.Debug("get new stream addrs", zap.Object("addrs", gStream.cache.Sas))
		}
	}
}

// RegisterWork register stream addr
func (upa *UpdateAddr) RegisterWork() {
	for {
		time.Sleep(time.Duration(Conf.EtcdC.ReportTime) * time.Second)
		upa.Put(upa.ReportKey, Conf.GrpcC.Addr, Conf.EtcdC.TTL)
	}
}

func (upa *UpdateAddr) Put(key, value string, ttl int64) error {
	Grant, err := upa.Client.Grant(context.TODO(), ttl)
	if err != nil {
		Logger.Panic("Etcd", zap.Error(err), zap.Int64("@ReportTime", ttl))
	}
	_, err = upa.Client.Put(context.TODO(), key, value, clientv3.WithLease(Grant.ID))
	if err != nil {
		Logger.Error("Put", zap.String("@key", key), zap.String("@value", value), zap.Error(err))
		return err
	}
	return nil
}

// NewUpdateAddr new UpdateAddr
func NewUpdateAddr() *UpdateAddr {
	upa := &UpdateAddr{}
	return upa
}

func GetRegisterKey() (string, error) {
	host, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%d", host, os.Getpid()), nil
}
