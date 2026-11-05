package ioc

import (
  "time"
  "weibook/internal/web/middleware"
  web "weibook/internal/web/user"

  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "github.com/redis/go-redis/v9"
  "github.com/ulule/limiter/v3"
  limiterGin "github.com/ulule/limiter/v3/drivers/middleware/gin"
  limiterRedis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func InitGin(middles []gin.HandlerFunc, handler *web.UserHandler) *gin.Engine {
  server := gin.Default()
  server.Use(middles...)
  handler.RegisterRoutes(server)
  return server
}

func InitMiddlewares(redis redis.Cmdable) []gin.HandlerFunc {
  return []gin.HandlerFunc{
    middleware.NewLoginJWTMiddleBuilder().Build(),
    CORSHandler(),
    initRateLimit(redis),
  }
}

func CORSHandler() gin.HandlerFunc {
  return cors.New(cors.Config{
    AllowOrigins: []string{"http://localhost:3000"},
    //AllowMethods:     []string{"POST", "GET"},
    // x-jwt-token 是前端请求携带的token字段
    AllowHeaders:     []string{"Content-Type", "x-jwt-token"},
    ExposeHeaders:    []string{"x-jwt-token"},
    AllowCredentials: true,
    //AllowOriginFunc: func(origin string) bool {
    //  return origin == "http://localhost:3000"
    //},
    MaxAge: 12 * time.Hour,
  })
}

func initRateLimit(redisCmd redis.Cmdable) gin.HandlerFunc {
  // 限流
  rate, err := limiter.NewRateFromFormatted("100-S")
  if err != nil {
    panic("Failed to initialize rate limiter")
  }
  client := redisCmd.(*redis.Client)
  store, err := limiterRedis.NewStoreWithOptions(client, limiter.StoreOptions{
    Prefix: "limiter_prefix_",
  })
  if err != nil {
    panic("Failed to initialize limit store")
  }
  mw := limiterGin.NewMiddleware(limiter.New(store, rate))
  return mw
}
