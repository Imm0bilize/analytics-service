package service

import "context"

func (t *TasksService) CreateTask(ctx context.Context, taskID string) error {
	if err := t.db.CreateTask(ctx, taskID); err != nil {
		t.l.Errorf("error while writing new task in db: err {%s}, taskID {%s}", err.Error(), taskID)
		return err
	}
	return nil
}

func (t *TasksService) SetTimeStart(ctx context.Context, taskID, login, startTime, newState string) error {
	if err := t.db.SetTimeStart(ctx, taskID, login, startTime, newState); err != nil {
		t.l.Errorf("error while set start time in db: err {%s}, taskID {%s}", err.Error(), taskID)
		return err
	}
	return nil
}

func (t *TasksService) SetTimeEnd(ctx context.Context, taskID, login, endTime, newState string) error {
	if err := t.db.SetTimeEnd(ctx, taskID, login, endTime, newState); err != nil {
		t.l.Errorf("error while set end time in db: err {%s}, taskID {%s}", err.Error(), taskID)
		return err
	}
	return nil
}
