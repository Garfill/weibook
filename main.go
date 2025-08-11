package main

import (
  "github.com/gin-contrib/cors"
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/cookie"
  "github.com/gin-gonic/gin"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
  "time"
  "weibook/internal/repo"
  "weibook/internal/repo/dao"
  "weibook/internal/service"
  "weibook/internal/www/middleware"
  "weibook/internal/www/user"
)

func main() {
  // 数据库初始化
  db := initDB()
  // user 初始化
  userHandler := initUser(db)
  //初始化gin
  server := initServer()

  // session
  initSession(server)

  // 注册handler
  userHandler.RegisterRoutes(server)

  server.Run(":8080")
}

func initDB() *gorm.DB {
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

  return db
}

func initServer() *gin.Engine {
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

  return server
}

func initSession(server *gin.Engine) {
  store := cookie.NewStore([]byte("secret"))
  server.Use(sessions.Sessions("wei_session", store))

  server.Use(middleware.NewLoginMiddleBuilder().Build())
}

func initUser(db *gorm.DB) *user.UserHandler {
  dao := dao.NewUserDAO(db)
  repo := repo.NewUserRepo(dao)
  svc := service.NewUserService(repo)
  userHandler := user.NewUserHandler(svc)

  return userHandler
}
