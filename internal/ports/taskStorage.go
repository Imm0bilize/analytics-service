package ports

import (
	"context"
)

type TaskStorage interface {
	GetTotalTime(ctx context.Context) (string, error)
	GetNumAgreedTasks(ctx context.Context) (int, error)
	GetNumRejectedTasks(ctx context.Context) (int, error)

	// TODO: add createTask etc
}
