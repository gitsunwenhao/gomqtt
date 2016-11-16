package service

type Cache struct {
	Ols *OlineUsers   //在线用户列表
	Ofs *OfflineUsers //离线用户列表
	Sas *StreamAddrs  //stream地址缓存列表
}

func NewCache() *Cache {
	cache := &Cache{
		Ols: NewOlineUsers(),
		Ofs: NewOfflineUsers(),
		Sas: NewStreamAddrs(),
	}
	return cache
}

func (cache *Cache) Init() {

}

func (cache *Cache) Start() {

}

func (cache *Cache) Close() error {
	return nil
}
