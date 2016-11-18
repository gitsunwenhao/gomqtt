package service

type Cache struct {
	As  *Accounts    //用户列表
	Sas *StreamAddrs //stream地址缓存列表
}

func NewCache() *Cache {
	cache := &Cache{
		As:  NewAccounts(),
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
