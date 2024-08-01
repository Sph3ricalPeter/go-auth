package storage

import (
	"fmt"
)

type DummyStorage struct {
	users         map[string]string
	refreshTokens map[string]bool
}

func NewDummyStorage() *DummyStorage {
	return &DummyStorage{
		users:         make(map[string]string),
		refreshTokens: make(map[string]bool),
	}
}

func (s *DummyStorage) CreateUser(username, password string) error {
	if _, ok := s.users[username]; ok {
		return fmt.Errorf("user %s already exists", username)
	}
	s.users[username] = password
	return nil
}

func (s *DummyStorage) GetUser(username string) (string, error) {
	if _, ok := s.users[username]; !ok {
		return "", fmt.Errorf("user %s does not exist", username)
	}
	return s.users[username], nil
}

func (s *DummyStorage) VerifyUser(username, password string) error {
	if p, ok := s.users[username]; !ok || p != password {
		return fmt.Errorf("invalid username or password")
	}
	return nil
}

func (s *DummyStorage) RegisterRefreshToken(refreshToken string) {
	s.refreshTokens[refreshToken] = true
}

func (s *DummyStorage) DeleteRefreshToken(refreshToken string) {
	delete(s.refreshTokens, refreshToken)
}

func (s *DummyStorage) IsRefreshTokenValid(refreshToken string) bool {
	valid, exists := s.refreshTokens[refreshToken]
	return exists && valid
}
