package internal

import (
	"context"
)

type RedisRepository[M any] interface {
	Create(ctx context.Context, key string, exp *M, seconds int) (err error)
	Get(ctx context.Context, key string) (res *M, err error)
	Delete(ctx context.Context, key string) (err error)
	Sadd(ctx context.Context, key string, value string) (err error)
	Sadds(ctx context.Context, key string, values []string) (err error)
	Srem(ctx context.Context, key string, value string) (err error)
	SIsMember(ctx context.Context, key string, value string) (isMember bool, err error)
	// SMembers(ctx context.Context, key string) (values []string, err error)
}
