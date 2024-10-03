package config

import "fmt"

const (
	Host      string = "localhost"
	Port      string = "8080"
	AccSecret string = "acc-secret"
	RefSecret string = "ref-secret"
)

const (
	DbUser string = "dev"
	DbPass string = "123123"
	DbHost string = "localhost"
	DbPort string = "5432"
	DbName string = "go-auth"
)

func DbConnStr() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DbUser, DbPass, DbHost, DbPort, DbName)
}
