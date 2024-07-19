package main

import (
	"errors"
	"fmt"
	"sync"
)

// TODO #1: implement in-memory user store

var userStore = NewUserStore()

func NewUserStore() *UserStore {
	return &UserStore{data: make(map[string]UserInfo)}
}

type UserStore struct {
	mu   sync.Mutex
	data map[string]UserInfo
}

func (u *UserStore) Save(info UserInfo) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	_, ok := u.data[info.UserName]
	if ok {
		return errors.New("user already exists")
	}
	fmt.Println(u.data, info)
	u.data[info.UserName] = info
	return nil
}

func (u *UserStore) Get(username string) (UserInfo, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	info, ok := u.data[username]
	if ok {
		return info, nil
	}
	return UserInfo{}, errors.New("user not found")
}

type UserInfo struct {
	// TODO: implement me  (username, password, full_name, address)
	UserName string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Address  string `json:"address"`
}
