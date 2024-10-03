package main

import (
	"flag"
	"fmt"
	"github.com/Sph3ricalPeter/go-auth/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"strconv"
)

func main() {
	m := newMigrate()

	flag.BoolFunc("up", "migrate up", func(string) error {
		if err := m.Up(); err != nil {
			panic(fmt.Errorf("error while migrating up: %w", err))
		}
		return nil
	})
	flag.BoolFunc("down", "migrate down", func(string) error {
		if err := m.Down(); err != nil {
			panic(fmt.Errorf("error while migrating down: %w", err))
		}
		return nil
	})
	flag.BoolFunc("force", "migrate force", func(v string) error {
		if err := migrateForceVersion(m, v); err != nil {
			panic(fmt.Errorf("error while migrating force: %w", err))
		}
		return nil
	})

	flag.Parse()
}

func newMigrate() *migrate.Migrate {
	m, err := migrate.New(
		"file://db/migrations",
		config.DbConnStr(),
	)
	if err != nil {
		panic(err)
	}
	return m
}

func migrateForceVersion(m *migrate.Migrate, v string) error {
	version, err := strconv.Atoi(v)
	if err != nil {
		return err
	}
	return m.Force(version)
}
