package main

import (
  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
  "time"
  "weibook/internal/repo"
  "weibook/internal/repo/dao"
  "weibook/internal/service"
  "weibook/internal/www/user"
)

func main() {
  // 初始化mysql连接
  db, err := gorm.Open(mysql.Open("root:12345678@tcp(localhost:13306)/weibook?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
  if err != nil {
    // panic 会令整个 goroutine 结束
    panic("failed to connect database")
  }
  err = dao.InitTable(db)
  if err != nil {
    panic("failed to init table")
  }

  // user 初始化
  dao := dao.NewUserDAO(db)
  repo := repo.NewUserRepo(dao)
  svc := service.NewUserService(repo)
  userHandler := user.NewUserHandler(svc)

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
  userHandler.RegisterRoutes(server)

  server.Run(":8080")
}
