package repository

import (
	pg "analytic-service/pkg/postgre"
)

const (
	taskStateAccepted   = "accepted"
	taskStateRejected   = "rejected"
	taskStateCreated    = "created"
	taskStateProcessing = "processing"
)

type Repository struct {
	*pg.DB
}

func NewPgRepo(db *pg.DB) *Repository {
	return &Repository{db}
}
