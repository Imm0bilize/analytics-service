package postgre

import (
	"database/sql"
	"fmt"
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
		panic("the number of attempts to connect to the database must be a positive number")
	}

	startDelayTime := time.Second * 2

	for nAttempts > 0 {
		err := db.Ping()
		if err == nil {
			return nil
		}
		logger.Warningf("failed to connect to the database, retry via %v", nAttempts)
		time.Sleep(startDelayTime)
		startDelayTime *= 2
	}
	return ErrConnectionToDb
}

func New(logger logrus.FieldLogger, user, password, host, port string, nAttempts int) (*DB, error) {
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
	return &DB{Conn: db}, nil

}

func (d *DB) Shutdown() error {
	return d.Conn.Close()
}
