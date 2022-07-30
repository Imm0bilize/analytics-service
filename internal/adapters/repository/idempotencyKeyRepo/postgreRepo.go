package idempotencyKeyRepo

import (
	pg "analytic-service/pkg/postgre"
	"context"
	"fmt"
)

type IdempotencyKeyRepo struct {
	*pg.DB
}

func NewIdempotencyKeyRepo(db *pg.DB) *IdempotencyKeyRepo {
	return &IdempotencyKeyRepo{db}
}

func (i IdempotencyKeyRepo) CheckIdempotencyKeyInStore(ctx context.Context, key string) (bool, error) {
	var count int
	err := i.Conn.QueryRowContext(ctx, "select count(*) from tasks_app.idempotency_keys where tasks_app.idempotency_keys.key=$1", key).Scan(&count)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (i IdempotencyKeyRepo) Commit(ctx context.Context, key string) error {
	tx, err := i.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err = i.Conn.ExecContext(ctx, "insert into tasks_app.idempotency_keys (key) values ($1)", key); err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return fmt.Errorf("rollback execution error {%s} after request execution {%w}", txErr.Error(), err)
		} else {
			return fmt.Errorf("query execution error: %w", err)
		}

	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error when committing a transaction {%s}", err.Error())
	} else {
		return nil
	}
}
