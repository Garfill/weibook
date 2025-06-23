package main

import (
  "github.com/gin-gonic/gin"
  "weibook/internal/www/user"
)

func main() {
  server := gin.Default()

  userHandler := user.NewUserHandler()
  userHandler.RegisterRoutes(server)

  server.Run(":8080")
}
