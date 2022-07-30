package tasksRepo

import (
	pg "analytic-service/pkg/postgre"
)

const (
	taskStateAccepted   = "accepted"
	taskStateRejected   = "rejected"
	taskStateCreated    = "created"
	taskStateProcessing = "processing"
)

type TasksRepo struct {
	*pg.DB
}

func NewTasksRepo(db *pg.DB) *TasksRepo {
	return &TasksRepo{db}
}
