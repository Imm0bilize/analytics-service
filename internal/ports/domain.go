package ports

import (
	"analytic-service/internal/adapters/http/dto"
	"context"
)

type QueryDomain interface {
	GetTotalTime(ctx context.Context) (*dto.TotalTimeResponse, error)
	GetNumAgreedTasks(ctx context.Context) (*dto.NumAgreedTasksResponse, error)
	GetNumRejectedTasks(ctx context.Context) (*dto.NumRejectedTaskResponse, error)
}

type CommandDomain interface {
	CreateTask(ctx context.Context, taskID string) error
	SetTimeStart(ctx context.Context, id, login, startTime, newState string) error
	SetTimeEnd(ctx context.Context, id, login, endTime, newState string) error
}
