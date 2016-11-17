package service

import "sync"

type OlineUsers struct {
	sync.RWMutex
	Users map[string]*User
}

func NewOlineUsers() *OlineUsers {
	ous := &OlineUsers{
		Users: make(map[string]*User),
	}
	return ous
}

type OfflineUsers struct {
	sync.RWMutex
	Users map[string]*User
}

func NewOfflineUsers() *OfflineUsers {
	ofs := &OfflineUsers{
		Users: make(map[string]*User),
	}
	return ofs
}

type User struct {
}

func NewUser() *User {
	user := &User{}
	return user
}
