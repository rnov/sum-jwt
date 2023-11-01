package storage

import (
	"sync"
)

// UserAccess - is a virtual memory storage for user Tokens that are used to authenticate them
type UserAccess struct {
	*sync.RWMutex
	Storage map[string]bool
}

func NewUserAccess() *UserAccess {
	return &UserAccess{
		RWMutex: new(sync.RWMutex),
		Storage: make(map[string]bool),
	}
}

// AddUserToken - add user token to the storage.
func (u *UserAccess) AddUserToken(token string) bool {
	u.Lock()
	defer u.Unlock()
	u.Storage[token] = true
	return true
}

// IsActiveToken - checks whether a token is still active/valid in our system.
func (u *UserAccess) IsActiveToken(token string) bool {
	u.RLock()
	defer u.RUnlock()
	return u.Storage[token]
}
