package ports

import (
	"context"
)

type TaskStorage interface {
	// For clients
	GetTotalTime(ctx context.Context) (string, error)
	GetNumAgreedTasks(ctx context.Context) (int, error)
	GetNumRejectedTasks(ctx context.Context) (int, error)

	// For management service
	CreateTask(ctx context.Context, taskID string) error
	SetTimeStart(ctx context.Context, taskID, login, startTime, newState string) error
	SetTimeEnd(ctx context.Context, taskID, login, endTime, newState string) error
}
