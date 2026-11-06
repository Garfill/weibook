package middleware

import (
  "github.com/gin-gonic/gin"
  "github.com/redis/go-redis/v9"
  "github.com/ulule/limiter/v3"
  limiterGin "github.com/ulule/limiter/v3/drivers/middleware/gin"
  limiterRedis "github.com/ulule/limiter/v3/drivers/store/redis"
)

type RateLimiter struct {
  limiter *limiter.Limiter
}

func NewRateLimiter(rdb redis.Cmdable) *RateLimiter {
  // 限流
  rate, err := limiter.NewRateFromFormatted("1000-S")
  if err != nil {
    panic("Failed to initialize rate limiter")
  }
  client := rdb.(*redis.Client)
  store, err := limiterRedis.NewStoreWithOptions(client, limiter.StoreOptions{
    Prefix: "limiter_prefix_",
  })
  if err != nil {
    panic("Failed to initialize limit store")
  }
  l := limiter.New(store, rate)
  return &RateLimiter{
    limiter: l,
  }
}

func (r *RateLimiter) Build() gin.HandlerFunc {
  return limiterGin.NewMiddleware(r.limiter)
}
