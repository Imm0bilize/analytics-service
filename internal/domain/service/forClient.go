package service

import (
	"analytic-service/internal/domain/models/dto"
	"context"
)

func (t *TasksService) GetTotalTime(ctx context.Context) (*dto.TotalTimeDTO, error) {
	res, err := t.db.GetTotalTime(ctx)
	if err != nil {
		t.l.ErrorF("error during the execution of the request to get the total time: %s", err.Error())
		return nil, err
	}
	return &dto.TotalTimeDTO{Time: res}, nil
}

func (t *TasksService) GetNumAgreedTasks(ctx context.Context) (*dto.NumAgreedTasksDTO, error) {
	num, err := t.db.GetNumAgreedTasks(ctx)
	if err != nil {
		t.l.ErrorF("error during the execution of the request to get the total agreed tasks: %s", err.Error())
		return nil, err
	}
	return &dto.NumAgreedTasksDTO{Num: num}, nil
}

func (t *TasksService) GetNumRejectedTasks(ctx context.Context) (*dto.NumRejectedTaskDTO, error) {
	num, err := t.db.GetNumRejectedTasks(ctx)
	if err != nil {
		t.l.ErrorF("error during the execution of the request to get the total rejected tasks: %s", err.Error())
		return nil, err
	}
	return &dto.NumRejectedTaskDTO{Num: num}, nil
}
