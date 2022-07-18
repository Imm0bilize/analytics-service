package ports

import (
	"analytic-service/internal/domain/models/dto"
	"context"
)

type ClientDomain interface {
	GetTotalTime(ctx context.Context) (*dto.TotalTimeDTO, error)
	GetNumAgreedTasks(ctx context.Context) (*dto.NumAgreedTasksDTO, error)
	GetNumRejectedTasks(ctx context.Context) (*dto.NumRejectedTaskDTO, error)
}

type ManagementServerDomain interface {
	CreateTask(ctx context.Context, taskID string) error
	SetTimeStart(ctx context.Context, id, login, startTime, newState string) error
	SetTimeEnd(ctx context.Context, id, login, endTime, newState string) error
}
