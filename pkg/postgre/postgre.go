package postgre

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
	"time"
)

type DB struct {
	Conn *sql.DB
}

func attemptPingDB(db *sql.DB, logger logrus.FieldLogger, nAttempts int) error {
	if nAttempts <= 0 {
		return ErrNAttempts
	}

	startDelayTime := time.Second * 2

	for nAttempts > 0 {
		err := db.Ping()
		if err == nil {
			return nil
		}
		logger.Warningf("failed to connect to the database, retry via %v", nAttempts)
		time.Sleep(startDelayTime)

		nAttempts--
		startDelayTime *= 2
	}
	return ErrConnectionToDb
}

func makeMigrate(user, password, host, port string) error {
	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("postgresql://%s:%s@%s:%s/postgres?sslmode=disable", user, password, host, port),
	)
	if err != nil {
		return errors.New("failed to create an object for migration")
	}

	if err := m.Up(); err != nil {
		return errors.New("error during migration")
	}
	return nil
}

func New(logger logrus.FieldLogger, user, password, host, port string, nAttempts int, isNeedMigrate bool) (*DB, error) {
	cfg, err := pgx.ParseConfig(
		fmt.Sprintf("postgresql://%s:%s@%s:%s/postgres?sslmode=disable", user, password, host, port),
	)
	if err != nil {
		return nil, ErrParseConfigFile
	}
	cfg.PreferSimpleProtocol = true

	db := stdlib.OpenDB(*cfg)

	if err := attemptPingDB(db, logger, nAttempts); err != nil {
		return nil, err
	}

	if isNeedMigrate {
		if err := makeMigrate(user, password, host, port); err != nil {
			return nil, err
		}
	}
	return &DB{Conn: db}, nil

}

func (d *DB) Shutdown() error {
	return d.Conn.Close()
}
