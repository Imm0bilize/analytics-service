package service

import "context"

func (t *TasksService) CreateTask(ctx context.Context, taskID string) error {
	return t.db.CreateTask(ctx, taskID)
}

func (t *TasksService) SetTimeStart(ctx context.Context, taskID, login, startTime, newState string) error {
	return t.db.SetTimeStart(ctx, taskID, login, startTime, newState)
}

func (t *TasksService) SetTimeEnd(ctx context.Context, taskID, login, endTime, newState string) error {
	return t.db.SetTimeEnd(ctx, taskID, login, endTime, newState)
}
