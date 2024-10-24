package port

import "context"

// Cache defines the interface for interacting with a cache service
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration int) error
	Get(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key string) error
}
