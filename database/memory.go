package database

import (
	"slices"
	"sync"
)

type MemoryDBType struct {
	users map[string][]string
	mu    sync.RWMutex
}

var MemoryDB = MemoryDBType{
	users: make(map[string][]string),
}

func InMemoryDB() *MemoryDBType {
	return &MemoryDB
}

// Saves refresh token in the memory DB for access token generation
func (m *MemoryDBType) SaveRefreshToken(username string, refreshToken string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.users[username] = append(m.users[username], refreshToken)
	return nil
}

// Returns refresh tokens for the given user
func (m *MemoryDBType) GetRefreshTokens(username string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.users[username], nil
}

func (m *MemoryDBType) DeleteRefreshToken(username string, refreshToken string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, rToken := range m.users[username] {
		if rToken == refreshToken {
			m.users[username] = slices.Delete(m.users[username], i, i+1)
			return nil
		}
	}
	return nil
}

func (m *MemoryDBType) DeleteRefreshTokens(username string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.users, username)
	return nil
}
