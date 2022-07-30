package ports

import "context"

type IdempotencyKeyStorage interface {
	CheckIdempotencyKeyInStore(ctx context.Context, key string) (bool, error)
	Commit(ctx context.Context, key string) error
}
