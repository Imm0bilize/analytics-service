package tasksRepo

import (
	"context"
	"database/sql"
	"fmt"
)

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

func (r *TasksRepo) CreateTask(ctx context.Context, taskID string) error {
	tx, err := r.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = r.Conn.ExecContext(ctx, "insert into tasks_app.tasks_state (id) values ($1)", taskID)
	if err := closingTxIfError(err, tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error when committing a transaction {%s}", err.Error())
	}
	return nil
}

func (r *TasksRepo) SetTimeStart(ctx context.Context, taskID, login, startTime, newState string) error {
	tx, err := r.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = r.Conn.ExecContext(ctx, "insert into tasks_app.user_accept_time (task_id, email, time_start) values ($1, $2, $3)", taskID, login, startTime)
	if err := closingTxIfError(err, tx); err != nil {
		return err
	}

	if newState != "" {
		_, err = r.Conn.ExecContext(ctx, "update tasks_app.tasks_state set state=$1 where id=$2", newState, taskID)
		if err := closingTxIfError(err, tx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error when committing a transaction {%s}", err.Error())
	}
	return nil
}

func (r *TasksRepo) SetTimeEnd(ctx context.Context, taskID, login, endTime, newState string) error {
	tx, err := r.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = r.Conn.ExecContext(ctx, "update tasks_app.user_accept_time set time_end=$1 where task_id=$2 and email=$3", endTime, taskID, login)
	if err := closingTxIfError(err, tx); err != nil {
		return err
	}

	if newState == taskStateAccepted || newState == taskStateRejected {
		_, err = r.Conn.ExecContext(ctx, "update tasks_app.tasks_state set time_agreement=("+
			"select sum(tasks_app.user_accept_time.time_end - tasks_app.user_accept_time.time_start) from tasks_app.user_accept_time where task_id=$1)"+
			"where id=$2",
			taskID, taskID,
		)

		if err := closingTxIfError(err, tx); err != nil {
			return err
		}

		_, err = r.Conn.ExecContext(ctx, "update tasks_app.tasks_state set state=$1 where id=$2", newState, taskID)
		if err := closingTxIfError(err, tx); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error when committing a transaction {%s}", err.Error())
	}
	return nil
}
