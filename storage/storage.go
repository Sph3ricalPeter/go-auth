package storage

type Storage interface {
	CreateUser(username, password string) error
	GetUser(username string) (string, error)
	VerifyUser(username, password string) error
	RegisterRefreshToken(refreshToken string)
	DeleteRefreshToken(refreshToken string)
	IsRefreshTokenValid(refreshToken string) bool
}
