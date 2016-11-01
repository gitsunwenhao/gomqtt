package sessions

import (
	"fmt"
	"sync"
)

func init() {
	Register("mem", NewMemProvider())
}

type memProvider struct {
	st map[string]*Session
	mu sync.RWMutex
}

func NewMemProvider() *memProvider {
	return &memProvider{
		st: make(map[string]*Session),
	}
}

func (mp *memProvider) New(id string) (*Session, error) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.st[id] = &Session{id: id}
	return mp.st[id], nil
}

func (mp *memProvider) Get(id string) (*Session, error) {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	sess, ok := mp.st[id]
	if !ok {
		return nil, fmt.Errorf("store/Get: No session found for key %s", id)
	}

	return sess, nil
}

func (mp *memProvider) Del(id string) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	delete(mp.st, id)
}

func (mp *memProvider) Count() int {
	return len(mp.st)
}

func (mp *memProvider) Close() error {
	mp.st = make(map[string]*Session)
	return nil
}
