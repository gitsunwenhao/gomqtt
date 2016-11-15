package service

import (
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

type Stream struct {
	etcd        *Etcd
	streamAddrs *StreamAddrs
	grpc        *Grpc
}

var gStream *Stream

func New() *Stream {
	return &Stream{}
}

func (s *Stream) Init() {

	// init etcd
	etcd := NewEtcd()
	etcd.Init()

	// init stream addrs
	streamAddrs := NewStreamAddrs()
	streamAddrs.Init(etcd.ReportKey, Conf.GrpcC.Addr)

	// init grpc service
	grpc := NewGrpc()
	grpc.Init()

	//	init other

	s.streamAddrs = streamAddrs
	s.grpc = grpc
	s.etcd = etcd

	gStream = s
}

func (s *Stream) Start(isStatic bool) {
	log.Println(isStatic)
	loadConfig(isStatic)

	// stream 初始化所有功能服务
	s.Init()

	// grpc start
	s.grpc.Start()

	// etcd start
	s.etcd.Start()

	go httpStart()
}

func (s *Stream) Close() error {
	s.etcd.Close()
	return nil
}

func httpStart() {
	// 启动Http服务
	e := echo.New()

	// 启动配置热更新
	e.GET("/reload", reload)

	e.Run(standard.New(":8907"))
}
