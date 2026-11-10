package rateLimit

import "context"

type Limiter interface {
  // 是否有触发限流
  Limited(ctx context.Context, key string) (bool, error)
}
