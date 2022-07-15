package service

import "context"

func (t *TasksService) CreateTask(ctx context.Context, taskID string) error {
	return nil
}

func (t *TasksService) UpdateTasksState(ctx context.Context, taskID string, state string) error {
	return nil
}

func (t *TasksService) SetTimeStart(ctx context.Context, id string, login string, startTime string) error {
	return nil
}

func (t *TasksService) SetTimeEnd(ctx context.Context, id string, login string, endTime string) error {
	return nil
}
