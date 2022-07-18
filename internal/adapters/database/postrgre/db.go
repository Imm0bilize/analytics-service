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
	taskStateAccepted   = "accepted"
	taskStateRejected   = "rejected"
	taskStateCreated    = "created"
	taskStateProcessing = "processing"
)

type db struct {
	instance *sql.DB
}

func New(cfg *config.Config) (*db, error) {
	dbCfg, err := pgx.ParseConfig(
		fmt.Sprintf("postgresql://%s:%s@%s:%s/postgres?sslmode=disable", cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port),
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

// Query

func (d *db) getNumRecordsByState(ctx context.Context, state string) (int, error) {
	var count int
	err := d.instance.QueryRowContext(ctx, "select count(*) from tasks_app.tasks_state where tasks_app.tasks_state.state=$1", state).Scan(&count)
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
	count, err := d.getNumRecordsByState(ctx, taskStateAccepted)
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

func closingTxIfError(err error, tx *sql.Tx) error {
	if err != nil {
		txErr := tx.Rollback()
		if txErr != nil {
			return fmt.Errorf("rollback execution error {%s} after request execution {%w}", txErr.Error(), err)
		}
		return fmt.Errorf("query execution error: %w", err)
	}
	return nil
}

// Commands

func (d *db) CreateTask(ctx context.Context, taskID string) error {
	tx, err := d.instance.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = d.instance.ExecContext(ctx, "insert into tasks_app.tasks_state (id) values ($1)", taskID)
	if err := closingTxIfError(err, tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error when committing a transaction {%s}", err.Error())
	}
	return nil
}

func (d *db) SetTimeStart(ctx context.Context, taskID, login, startTime, newState string) error {
	tx, err := d.instance.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = d.instance.ExecContext(ctx, "insert into tasks_app.user_accept_time (task_id, email, time_start) values ($1, $2, $3)", taskID, login, startTime)
	if err := closingTxIfError(err, tx); err != nil {
		return err
	}

	if newState != "" {
		_, err = d.instance.ExecContext(ctx, "update tasks_app.tasks_state set state=$1 where id=$2", newState, taskID)
		if err := closingTxIfError(err, tx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error when committing a transaction {%s}", err.Error())
	}
	return nil
}

func (d *db) SetTimeEnd(ctx context.Context, taskID, login, endTime, newState string) error {
	tx, err := d.instance.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = d.instance.ExecContext(ctx, "update tasks_app.user_accept_time set time_end=$1 where task_id=$2 and email=$3", endTime, taskID, login)
	if err := closingTxIfError(err, tx); err != nil {
		return err
	}

	if newState == taskStateAccepted || newState == taskStateRejected {
		_, err = d.instance.ExecContext(ctx, "update tasks_app.tasks_state set time_agreement=("+
			"select sum(tasks_app.user_accept_time.time_end - tasks_app.user_accept_time.time_start) from tasks_app.user_accept_time where task_id=$1)"+
			"where id=$2",
			taskID, taskID,
		)

		if err := closingTxIfError(err, tx); err != nil {
			return err
		}

		_, err = d.instance.ExecContext(ctx, "update tasks_app.tasks_state set state=$1 where id=$2", newState, taskID)
		if err := closingTxIfError(err, tx); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error when committing a transaction {%s}", err.Error())
	}
	return nil
}
