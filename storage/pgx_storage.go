package storage

import (
	"context"
	"fmt"
	"github.com/Sph3ricalPeter/go-auth/config"
	"github.com/jackc/pgx/v5"
)

type PgxStorage struct {
	conn *pgx.Conn
}

func NewPgxStorage() *PgxStorage {
	conn, err := pgx.Connect(context.Background(), config.DbConnStr())
	if err != nil {
		panic(err)
	}
	return &PgxStorage{
		conn: conn,
	}
}

func (s *PgxStorage) CreateUser(username, password string) error {
	_, err := s.conn.Exec(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2)", username, password)
	return err
}

func (s *PgxStorage) GetUser(username string) (*User, error) {
	user := &User{}
	err := s.conn.QueryRow(context.Background(), "SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.Id, &user.Username, &user.Password)
	return user, err
}

func (s *PgxStorage) VerifyUser(username, password string) error {
	var p string
	err := s.conn.QueryRow(context.Background(), "SELECT password FROM users WHERE username = $1", username).Scan(&p)
	if err != nil {
		return err
	}
	if p != password {
		return fmt.Errorf("invalid username or password")
	}
	return nil
}

func (s *PgxStorage) RegisterRefreshToken(refreshToken string) {
	_, _ = s.conn.Exec(context.Background(), "INSERT INTO refresh_tokens (token) VALUES ($1)", refreshToken)
}

func (s *PgxStorage) DeleteRefreshToken(refreshToken string) {
	_, _ = s.conn.Exec(context.Background(), "DELETE FROM refresh_tokens WHERE token = $1", refreshToken)
}

func (s *PgxStorage) IsRefreshTokenValid(refreshToken string) bool {
	var exists bool
	_ = s.conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM refresh_tokens WHERE token = $1)", refreshToken).Scan(&exists)
	return exists
}
