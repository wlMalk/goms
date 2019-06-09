package cache

import (
	"context"
)

type Cache interface {
	Set(ctx context.Context, key, value interface{}) (err error)
	Get(ctx context.Context, key interface{}) (value interface{}, err error)
}
