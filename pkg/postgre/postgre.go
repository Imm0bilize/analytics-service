package postgre

import (
	"database/sql"
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
		logger.Errorf("failed to connect to the database, retry via %v", startDelayTime)
		time.Sleep(startDelayTime)

		nAttempts--
		startDelayTime *= 2
	}
	return ErrConnectionToDb
}

func makeMigrate(user, password, host, port, dbName string) error {
	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName),
	)
	if err != nil {
		return fmt.Errorf("failed to create an object for migration: %s", err.Error())
	}

	if err := m.Up(); err != nil {
		return fmt.Errorf("error during migration: %s", err.Error())
	}
	return nil
}

func New(logger logrus.FieldLogger, user, password, host, port, dbName string, nAttempts int, isNeedMigrate bool) (*DB, error) {
	cfg, err := pgx.ParseConfig(
		fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName),
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
		if err := makeMigrate(user, password, host, port, dbName); err != nil {
			return nil, err
		}
	}
	return &DB{Conn: db}, nil

}

func (d *DB) Shutdown() error {
	return d.Conn.Close()
}
