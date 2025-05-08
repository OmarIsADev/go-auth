package database

import (
	"slices"
	"strings"
	"sync"

	"github.com/omarisadev/go-auth/models"
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

	if m.users[username] == nil {
		user, _ := GetUserByUsername(username)
		m.users[username] = strings.Split(user.RefreshTokens, "+")
	}

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

func (m *MemoryDBType) Save() error {

	for user, tokens := range m.users {
		DB.Model(&models.User{}).Where("username = ?", user).Update("refresh_tokens", strings.Join(tokens, "+"))
	}

	return nil
}
