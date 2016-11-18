package service

import "github.com/labstack/echo"

type Stream struct {
	upa   *UpdateAddr
	rpc   *Rpc
	cache *Cache
	hash  *Hash
}

var gStream *Stream

func New() *Stream {
	stream := &Stream{}
	return stream
}

func (s *Stream) Init() {

	// init etcd
	upa := NewUpdateAddr()
	upa.Init()

	// init cache
	cache := NewCache()
	cache.Init()

	// init rpc service
	rpc := NewRpc()
	rpc.Init()

	// hash 初始化
	hash := NewHash()
	//	init other

	s.rpc = rpc
	s.upa = upa
	s.cache = cache
	s.hash = hash

	gStream = s
}

func (s *Stream) Start(isStatic bool) {

	loadConfig(isStatic)

	// stream 初始化所有功能服务
	s.Init()

	// rpc start
	s.rpc.Start()

	// upa start
	s.upa.Start()

	go httpStart()
}

func (s *Stream) Close() error {
	s.upa.Close()
	s.rpc.Close()
	return nil
}

func httpStart() {
	// 启动Http服务
	e := echo.New()

	// 启动配置热更新
	e.GET("/reload", reload)

	// e.Run(standard.New(":8907"))

	err := e.Start(":8907")
	if err != nil {
		e.Logger.Fatal(err.Error())
	}
}

// GetAddrKey  获取上报stream地址的key
// func GetAddrKey(rootDir string) (string, error) {

// 	keylen := len(rootDir)
// 	if keylen > 0 && rootDir[keylen-1] != '/' {
// 		rootDir = rootDir + "/"
// 	}

// 	host, err := os.Hostname()
// 	if err != nil {
// 		return "", err
// 	}
// 	return fmt.Sprintf("%s%s-%d", rootDir, host, os.Getpid()), nil
// }
