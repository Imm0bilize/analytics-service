package repository

import (
	"context"
)

func (r *Repository) getNumRecordsByState(ctx context.Context, state string) (int, error) {
	var count int
	err := r.Conn.QueryRowContext(ctx, "select count(*) from tasks_app.tasks_state where tasks_app.tasks_state.state=$1", state).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) GetTotalTime(ctx context.Context) (string, error) {
	var totalTime string
	err := r.Conn.QueryRowContext(
		ctx,
		"select sum(tasks_app.tasks_state.time_agreement) from tasks_app.tasks_state where tasks_app.tasks_state.state!=$1 and tasks_app.tasks_state.state!=$2",
		taskStateCreated,
		taskStateProcessing,
	).Scan(&totalTime)
	if err != nil {
		return "", ErrCreateQuery
	}
	return totalTime, nil
}

func (r *Repository) GetNumAgreedTasks(ctx context.Context) (int, error) {
	count, err := r.getNumRecordsByState(ctx, taskStateAccepted)
	if err != nil {
		return 0, ErrCreateQuery
	}
	return count, nil
}

func (r *Repository) GetNumRejectedTasks(ctx context.Context) (int, error) {
	count, err := r.getNumRecordsByState(ctx, taskStateRejected)
	if err != nil {
		return 0, ErrCreateQuery
	}
	return count, nil
}
