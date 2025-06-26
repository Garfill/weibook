package main

import (
  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "time"
  "weibook/internal/www/user"
)

func main() {
  server := gin.Default()

  // cors
  server.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"POST"},
    AllowHeaders:     []string{"Origin"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    AllowOriginFunc: func(origin string) bool {
      return origin == "http://localhost:3000"
    },
    MaxAge: 12 * time.Hour,
  }))

  // 注册handler
  userHandler := user.NewUserHandler()
  userHandler.RegisterRoutes(server)

  server.Run(":8080")
}
