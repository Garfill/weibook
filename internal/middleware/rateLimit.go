package middleware

import (
  "fmt"

  "github.com/gin-gonic/gin"
  "github.com/redis/go-redis/v9"
  "github.com/ulule/limiter/v3"
  limiterGin "github.com/ulule/limiter/v3/drivers/middleware/gin"
  limiterRedis "github.com/ulule/limiter/v3/drivers/store/redis"
)

type RateLimiter struct {
  limiter *limiter.Limiter
  prefix  string
}

// 判断限流情况
func (r *RateLimiter) Limited(ctx *gin.Context) (bool, error) {
  key := fmt.Sprintf("%s:%s", r.prefix, ctx.ClientIP())
  limiterContext, err := r.limiter.Get(ctx, key)
  if err != nil {
    // 查询限流情况报错
    return true, err
  }
  return limiterContext.Reached, nil
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
    prefix:  "limiter_prefix_",
  }
}

func (r *RateLimiter) Build() gin.HandlerFunc {
  return limiterGin.NewMiddleware(r.limiter)
}
