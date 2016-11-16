package service

import "sync"
import "github.com/uber-go/zap"

type StreamAddrs struct {
	sync.RWMutex
	Addrs map[string]string
}

func NewStreamAddrs() *StreamAddrs {
	sa := &StreamAddrs{
		Addrs: make(map[string]string),
	}
	return sa
}

// Init 将本机stream 自己的grpc地址保存进map
func (sa *StreamAddrs) Init(key, addr string) {
	sa.Lock()
	sa.Addrs[key] = addr
	sa.Unlock()
	Logger.Info("Stre  amAddrs", zap.String("key", key), zap.String("addr", addr))
}

// Add insert key-value, if exist return true,else return false
func (sa *StreamAddrs) Add(key, addr string) bool {
	sa.RLock()
	if _, ok := sa.Addrs[key]; ok {
		sa.RUnlock()
		return true
	}
	sa.RUnlock()

	sa.Lock()
	sa.Addrs[key] = addr
	sa.Unlock()
	return false
}

func (sa *StreamAddrs) Get(key string) (string, bool) {
	sa.RLock()
	addr, ok := sa.Addrs[key]
	sa.RUnlock()

	return addr, ok
}

func (sa *StreamAddrs) Del(key string) {
	sa.Lock()
	delete(sa.Addrs, key)
	sa.Unlock()
}
