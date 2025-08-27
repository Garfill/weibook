package main

import (
  "time"
  "weibook/internal/repo"
  "weibook/internal/repo/dao"
  "weibook/internal/service"
  "weibook/internal/www/user"

  "github.com/gin-contrib/cors"
  "github.com/gin-contrib/sessions"
  "github.com/gin-contrib/sessions/redis"
  "github.com/gin-gonic/gin"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
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
  db, err := gorm.Open(mysql.Open("root:12345678@tcp(localhost:3306)/weibook?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
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
  // V1 版本，数据设置在cookie 不安全
  //store := cookie.NewStore([]byte("secret"))

  // 信息安全：
  //1. 身份认证key和数据加密key
  // 2. 数据操作权限
  // 注意这里的key要使用随机生成的，除了保存到redis还会保存到cookie，不过cookie里没有数据只有sid
  store, _ := redis.NewStore(
    10,
    "tcp", "localhost:6379", "", "",
    []byte("OEnEc62tqMFBOYRQEWQKmFWBvcpViJHV"), []byte("5H5v7Qqhct6EQBZ0DfsibYwi1J2l52xh"))
  redis.SetKeyPrefix(store, "wei_session")
  server.Use(sessions.Sessions("wei_session", store))

  // 自定义中间件
  //server.Use(middleware.NewLoginMiddleBuilder().Build())
}

func initUser(db *gorm.DB) *user.UserHandler {
  dao := dao.NewUserDAO(db)
  repo := repo.NewUserRepo(dao)
  svc := service.NewUserService(repo)
  userHandler := user.NewUserHandler(svc)

  return userHandler
}
