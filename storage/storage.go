package storage

import (
	"fmt"
)

type Storage interface {
	CreateUser(username, password string) error
	GetUser(username string) (*User, error)
	VerifyUser(username, password string) error
	RegisterRefreshToken(refreshToken string)
	DeleteRefreshToken(refreshToken string)
	IsRefreshTokenValid(refreshToken string) bool
}

type User struct {
	Id       int
	Username string
	Password string
}

type DummyStorage struct {
	users         map[string]*User
	refreshTokens map[string]bool
}

func NewDummyStorage() *DummyStorage {
	return &DummyStorage{
		users:         make(map[string]*User),
		refreshTokens: make(map[string]bool),
	}
}

func (s *DummyStorage) CreateUser(username, password string) error {
	if _, ok := s.users[username]; ok {
		return fmt.Errorf("user %s already exists", username)
	}
	s.users[username] = &User{
		Username: username,
		Password: password,
	}
	return nil
}

func (s *DummyStorage) GetUser(username string) (*User, error) {
	if _, ok := s.users[username]; !ok {
		return nil, fmt.Errorf("user %s does not exist", username)
	}
	return s.users[username], nil
}

func (s *DummyStorage) VerifyUser(username, password string) error {
	if p, ok := s.users[username]; !ok || p.Password != password {
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
