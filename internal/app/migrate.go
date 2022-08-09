package app

import (
	"errors"
	"fmt"
	"strings"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"

	"github.com/dimk00z/go-shortener-praktikum/pkg/logger"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func doMigrations(l *logger.Logger, databaseURL string) {
	sslMode := "?sslmode=disable"
	if !strings.Contains(databaseURL, sslMode) {
		databaseURL += sslMode
	}
	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", databaseURL)
		if err == nil {
			break
		}
		l.Error(err)
		l.Debug(fmt.Printf("Migrate: postgres is trying to connect, attempts left: %d", attempts))
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		l.Fatal(fmt.Printf("Migrate: postgres connect error: %s", err))
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		l.Fatal(fmt.Printf("Migrate: up error: %s", err))
	}

	if errors.Is(err, migrate.ErrNoChange) {
		l.Debug("Migrate: no change")
		return
	}

	l.Debug("Migrate: up success")
}
