package ioc

import (
  "time"
  "weibook/internal/middleware"
  "weibook/internal/web/user"
  "weibook/internal/web/wechat"

  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "github.com/redis/go-redis/v9"
)

func InitGin(middles []gin.HandlerFunc, userHdl *user.UserHandler, wechatHdl *wechat.OAuth2WechatHandler) *gin.Engine {
  server := gin.Default()
  server.Use(middles...)
  userHdl.RegisterRoutes(server)
  wechatHdl.RegisterRoutes(server)
  return server
}

func InitMiddlewares(rdb redis.Cmdable) []gin.HandlerFunc {
  return []gin.HandlerFunc{
    CORSHandler(),
    middleware.NewRateLimiter(rdb).Build(),
    middleware.NewLoginJWTMiddleBuilder().Build(),
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
