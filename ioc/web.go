package ioc

import (
  "time"
  "weibook/internal/web/middleware"
  web "weibook/internal/web/user"

  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
)

func InitGin(middles []gin.HandlerFunc, handler *web.UserHandler) *gin.Engine {
  server := gin.Default()
  server.Use(middles...)
  handler.RegisterRoutes(server)
  return server
}

func InitMiddlewares() []gin.HandlerFunc {
  return []gin.HandlerFunc{
    middleware.NewLoginJWTMiddleBuilder().Build(), // jwt
    CORSHandler(),
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
