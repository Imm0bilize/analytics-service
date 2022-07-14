package postrgre

import (
	"analytic-service/internal/adapters/database"
	"analytic-service/internal/config"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

const (
	taskStateAgreed     = "agreed"
	taskStateRejected   = "rejected"
	taskStateCreated    = "created"
	taskStateProcessing = "processing"
)

type db struct {
	instance *sql.DB
}

func New(cfg *config.Config) (*db, error) {
	dbCfg, err := pgx.ParseConfig(
		fmt.Sprintf("postrgresql://%s:%s@%s:%s/app?sslmode=disable", cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port),
	)

	if err != nil {
		return nil, database.ErrParseConfigFile
	}
	dbCfg.PreferSimpleProtocol = true

	db_ := stdlib.OpenDB(*dbCfg)

	if err := db_.Ping(); err != nil {
		return nil, database.ErrConnectionToDb
	} else {
		return &db{instance: db_}, nil
	}
}

func (d *db) Shutdown() error {
	return d.instance.Close()
}

func (d *db) getNumRecordsByState(ctx context.Context, state string) (int, error) {
	var count int
	err := d.instance.QueryRowContext(ctx, "select count(*) from tasks_app.tasks_state where state=$1", state).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *db) GetTotalTime(ctx context.Context) (string, error) {
	var totalTime string
	err := d.instance.QueryRowContext(
		ctx,
		"select sum(tasks_app.tasks_state.time_agreement) from tasks_app.tasks_state where tasks_app.tasks_state.state!=$1 and tasks_app.tasks_state.state!=$2",
		taskStateCreated,
		taskStateProcessing,
	).Scan(&totalTime)
	if err != nil {
		return "", database.ErrCreateQuery
	}
	return totalTime, nil
}

func (d *db) GetNumAgreedTasks(ctx context.Context) (int, error) {
	count, err := d.getNumRecordsByState(ctx, taskStateAgreed)
	if err != nil {
		return 0, database.ErrCreateQuery
	}
	return count, nil
}

func (d *db) GetNumRejectedTasks(ctx context.Context) (int, error) {
	count, err := d.getNumRecordsByState(ctx, taskStateRejected)
	if err != nil {
		return 0, database.ErrCreateQuery
	}
	return count, nil
}
