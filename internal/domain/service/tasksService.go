package service

import (
	"analytic-service/internal/ports"
	"github.com/sirupsen/logrus"
)

type TasksService struct {
	db ports.TaskStorage
	l  logrus.FieldLogger
}

func New(db ports.TaskStorage, logger logrus.FieldLogger) *TasksService {
	return &TasksService{db: db, l: logger}
}
