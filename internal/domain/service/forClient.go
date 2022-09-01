package service

import (
	"analytic-service/internal/adapters/http/dto"
	"context"
)

func (t *TasksService) GetTotalTime(ctx context.Context) (*dto.TotalTimeResponse, error) {
	res, err := t.db.GetTotalTime(ctx)
	if err != nil {
		t.l.Errorf("error during the execution of the request to get the total time: %s", err.Error())
		return nil, err
	}
	return &dto.TotalTimeResponse{Time: res}, nil
}

func (t *TasksService) GetNumAgreedTasks(ctx context.Context) (*dto.NumAgreedTasksResponse, error) {
	num, err := t.db.GetNumAgreedTasks(ctx)
	if err != nil {
		t.l.Errorf("error during the execution of the request to get the total agreed tasks: %s", err.Error())
		return nil, err
	}
	return &dto.NumAgreedTasksResponse{Num: num}, nil
}

func (t *TasksService) GetNumRejectedTasks(ctx context.Context) (*dto.NumRejectedTaskResponse, error) {
	num, err := t.db.GetNumRejectedTasks(ctx)
	if err != nil {
		t.l.Errorf("error during the execution of the request to get the total rejected tasks: %s", err.Error())
		return nil, err
	}
	return &dto.NumRejectedTaskResponse{Num: num}, nil
}
