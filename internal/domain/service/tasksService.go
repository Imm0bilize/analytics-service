package service

import (
	"analytic-service/internal/ports"
	"analytic-service/pkg/logging"
)

type TasksService struct {
	db ports.TaskStorage
	l  logging.ILogger
}

func New(db ports.TaskStorage, logger logging.ILogger) *TasksService {
	return &TasksService{db: db, l: logger}
}
