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
